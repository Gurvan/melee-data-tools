package fighter

import (
	"fmt"
	"io"

	"github.com/Gurvan/melee-data-tools/animation"
	"github.com/Gurvan/melee-data-tools/animation/track"
	"github.com/Gurvan/melee-data-tools/binread"
	"github.com/Gurvan/melee-data-tools/descriptor"
	. "github.com/Gurvan/melee-data-tools/lib"
	"github.com/Gurvan/melee-data-tools/logger"
)

// PartPoseCount is how many poses a body part can be authored, one per way of gripping an item.
const PartPoseCount = 4

// PartPoseTable is the table at fighter data offset 0x1C. A set poses one part of the body over
// whatever the action animation is doing: picking up an item puts the carrying hand into the pose
// its hold kind asks for, and a script trigger can pose any other set the same way.
type PartPoseTable []PartPoseSet

// PartPoseSet is one body part's poses. Characters have an entry per hand plus, usually, smaller
// ones for other parts a script or a carried item poses.
type PartPoseSet struct {
	// StartBone is the bone the pose starts at. A pose walks the character's bones from
	// there in the order the model tree lists them, one bone per joint of the pose, so the
	// i-th joint of a pose drives bone StartBone+i.
	StartBone uint16
	// Bones are the bones this set claims. They are the ones a pose is allowed to move, and
	// the ones handed back to the action animation when the pose ends; a set can start its
	// walk above the bones it claims, and two sets can walk the same bones and claim
	// different ones.
	Bones []uint8
	// Poses is indexed by pose index, and holds a nil for a pose the character has none of.
	Poses []*PartPoseJoint
}

// PartPoseJoint is one joint of a pose, a tree mirroring the shape of the bones it drives.
type PartPoseJoint struct {
	// Tracks are the curves posing this joint, in the same form as an action's animation
	// tracks, and are empty for a joint the pose leaves to the action animation.
	Tracks []animation.Track
	// FrameCount is the last frame of the joint's curves.
	FrameCount float32
	Child      *PartPoseJoint
	Next       *PartPoseJoint
}

var (
	_ binread.BinReader = (*PartPoseTable)(nil)
	_ binread.BinReader = (*PartPoseSet)(nil)
	_ binread.BinReader = (*PartPoseJoint)(nil)
)

// Joints flattens one of the set's poses into the walk the game applies it with, so entry i
// drives bone StartBone+i. It is nil for a pose the character has none of.
func (s PartPoseSet) Joints(pose int) []*PartPoseJoint {
	if pose < 0 || pose >= len(s.Poses) || s.Poses[pose] == nil {
		return nil
	}
	return flattenPose(s.Poses[pose])
}

// claims reports whether the set is allowed to move the bone the i-th joint of a pose drives.
func (s PartPoseSet) claims(i int) bool {
	bone := int(s.StartBone) + i
	for _, b := range s.Bones {
		if int(b) == bone {
			return true
		}
	}
	return false
}

// posesClaimedBones reports whether a pose only moves bones the set claims, which is what
// tells a real pose apart from whatever sits past the end of a pose table.
func (s PartPoseSet) posesClaimedBones(pose *PartPoseJoint) bool {
	for i, joint := range flattenPose(pose) {
		if len(joint.Tracks) > 0 && !s.claims(i) {
			return false
		}
	}
	return true
}

// nextRootCount returns how many entries of entrySize fit between addr and the first root
// after it, and 0 when no root follows. Roots end a block the same way a pointer target does,
// but the relocation table does not list them.
func nextRootCount(args []Args, addr Addr, entrySize uint32) int {
	end := ^uint32(0)
	for _, args := range args {
		desc, ok := args["descriptor"].(*descriptor.Descriptor)
		if !ok || desc == nil {
			continue
		}
		for _, root := range desc.Roots {
			if uint32(root.Offset) > uint32(addr) && uint32(root.Offset) < end {
				end = uint32(root.Offset)
			}
		}
	}
	if end == ^uint32(0) {
		return 0
	}
	return int((end - uint32(addr)) / entrySize)
}

func flattenPose(joint *PartPoseJoint) []*PartPoseJoint {
	if joint == nil {
		return nil
	}
	joints := []*PartPoseJoint{joint}
	joints = append(joints, flattenPose(joint.Child)...)
	joints = append(joints, flattenPose(joint.Next)...)
	return joints
}

func (t *PartPoseTable) BinRead(r *binread.Reader, args ...Args) error {
	count, ok := relocEntryCount(r, args, 4)
	if !ok {
		logger.Warning.Println("PartPoseTable needs to be parsed as a part of FighterFile")
		return nil
	}

	sets := make([]PartPoseSet, count)
	for i := range sets {
		var p OptionalPtr[PartPoseSet]
		if err := r.Decode(&p, args...); err != nil {
			return fmt.Errorf("PartPoseSet(Index:%d): %w", i, err)
		}
		if p.ValuePtr != nil {
			sets[i] = *p.ValuePtr
		}
	}

	*t = sets
	return nil
}

func (s *PartPoseSet) BinRead(r *binread.Reader, args ...Args) error {
	var boneCount uint16
	if err := r.Decode(&s.StartBone); err != nil {
		return err
	}
	if err := r.Decode(&boneCount); err != nil {
		return err
	}

	var bonesAddr, posesAddr Addr
	if err := r.Decode(&bonesAddr); err != nil {
		return err
	}
	if err := r.Decode(&posesAddr); err != nil {
		return err
	}

	before := r.CurrentPosition()

	if bonesAddr != Addr(0x20) {
		if _, err := r.Seek(bonesAddr.ToSeek(), io.SeekStart); err != nil {
			return err
		}
		s.Bones = make([]uint8, boneCount)
		if err := r.Decode(&s.Bones); err != nil {
			return err
		}
	}

	if posesAddr != Addr(0x20) {
		if _, err := r.Seek(posesAddr.ToSeek(), io.SeekStart); err != nil {
			return err
		}
		// The pose table carries no count of its own, so it runs to whatever the file
		// mentions next - a pointer target or a root - capped at the poses the game can
		// address.
		count, ok := relocEntryCount(r, args, 4)
		if !ok {
			logger.Warning.Println("PartPoseSet needs to be parsed as a part of FighterFile")
			return nil
		}
		if root := nextRootCount(args, posesAddr, 4); root > 0 && root < count {
			count = root
		}
		if count > PartPoseCount {
			count = PartPoseCount
		}
		s.Poses = make([]*PartPoseJoint, count)
		for i := range s.Poses {
			var p OptionalPtr[PartPoseJoint]
			if err := r.Decode(&p, args...); err != nil {
				return fmt.Errorf("PartPose(Index:%d): %w", i, err)
			}
			// Anything past the real end of the table is another structure that happens to
			// sit there, and reads as a pose that moves bones this set does not claim.
			if p.ValuePtr != nil && s.posesClaimedBones(p.ValuePtr) {
				s.Poses[i] = p.ValuePtr
			}
		}
	}

	_, err := r.Seek(before, io.SeekStart)
	return err
}

func (j *PartPoseJoint) BinRead(r *binread.Reader, args ...Args) error {
	var child, next OptionalPtr[PartPoseJoint]
	if err := r.Decode(&child, args...); err != nil {
		return err
	}
	if err := r.Decode(&next, args...); err != nil {
		return err
	}
	j.Child, j.Next = child.ValuePtr, next.ValuePtr

	var curvesAddr Addr
	if err := r.Decode(&curvesAddr); err != nil {
		return err
	}
	// Constraint animations and flags, neither of which the pose walk reads.
	var unused [2]uint32
	if err := r.Decode(&unused); err != nil {
		return err
	}
	if curvesAddr == Addr(0x20) {
		return nil
	}

	before := r.CurrentPosition()
	if _, err := r.Seek(curvesAddr.ToSeek(), io.SeekStart); err != nil {
		return err
	}

	var flags uint32
	if err := r.Decode(&flags); err != nil {
		return err
	}
	if err := r.Decode(&j.FrameCount); err != nil {
		return err
	}
	var trackAddr Addr
	if err := r.Decode(&trackAddr); err != nil {
		return err
	}
	for trackAddr != Addr(0x20) {
		if _, err := r.Seek(trackAddr.ToSeek(), io.SeekStart); err != nil {
			return err
		}
		t, nextAddr, err := readPoseTrack(r)
		if err != nil {
			return err
		}
		j.Tracks = append(j.Tracks, t)
		trackAddr = nextAddr
	}

	_, err := r.Seek(before, io.SeekStart)
	return err
}

// readPoseTrack reads one curve of a pose and where the next one lives. Poses store their
// curves as a linked list with a header of its own, where an action's animation stores the
// same curve data behind a packed per-bone table, so the header is read by hand and the keys
// are left to the shared decoder.
func readPoseTrack(r *binread.Reader) (animation.Track, Addr, error) {
	var t animation.Track

	var nextAddr Addr
	if err := r.Decode(&nextAddr); err != nil {
		return t, nextAddr, err
	}
	var length uint32
	if err := r.Decode(&length); err != nil {
		return t, nextAddr, err
	}
	var startFrame float32
	if err := r.Decode(&startFrame); err != nil {
		return t, nextAddr, err
	}
	var header [4]uint8
	if err := r.Decode(&header); err != nil {
		return t, nextAddr, err
	}

	t.DataLength = uint16(length)
	t.StartFrame = int16(startFrame)
	t.Type = track.Type(header[0])
	t.ValueFormat = track.DataFormat(header[1] & 0b11100000)
	t.ValueScale = 1 << (header[1] & 0b00011111)
	t.TangentFormat = track.DataFormat(header[2] & 0b11100000)
	t.TangentScale = 1 << (header[2] & 0b00011111)

	keyArgs := Args{
		"valueScale":    t.ValueScale,
		"valueFormat":   t.ValueFormat,
		"tangentScale":  t.TangentScale,
		"tangentFormat": t.TangentFormat,
		"dataLength":    t.DataLength,
	}
	err := r.Decode(&t.Keys, keyArgs)
	return t, nextAddr, err
}
