package fighter

import (
	"fmt"
	"io"

	"github.com/Gurvan/melee-data-tools/binread"
	. "github.com/Gurvan/melee-data-tools/lib"
	"github.com/Gurvan/melee-data-tools/model"
)

const maxDynamicBoneSets = 10

// DynamicBoneParameters contains the per-joint constants copied into the runtime
// dynamic-bone solver. The final entry in a chain is retained even though it has no
// child link, matching the source data.
type DynamicBoneParameters struct {
	Stiffness       float32
	Convergence     float32
	NaturalRotation model.Vec3[float32]
	_               float32
	MaxDeviation    float32
	Unknown1        model.Vec3[float32]
	Unknown2        model.Vec3[float32]
	Damping         float32
	MaxAngleChange  float32
}

// DynamicBoneSet describes one consecutive child-bone chain.
type DynamicBoneSet struct {
	RootBone            int32
	Parameters          []DynamicBoneParameters
	StiffnessMultiplier float32
	Unknown             float32
	Gravity             float32
}

func (s *DynamicBoneSet) BinRead(r *binread.Reader, args ...Args) error {
	var parametersOffset Addr
	var count uint32
	if err := r.Decode(&s.RootBone, args...); err != nil {
		return err
	}
	if err := r.Decode(&parametersOffset, args...); err != nil {
		return err
	}
	if err := r.Decode(&count, args...); err != nil {
		return err
	}
	if err := r.Decode(&s.StiffnessMultiplier, args...); err != nil {
		return err
	}
	if err := r.Decode(&s.Unknown, args...); err != nil {
		return err
	}
	if err := r.Decode(&s.Gravity, args...); err != nil {
		return err
	}

	if count == 0 || parametersOffset == Addr(0x20) {
		s.Parameters = []DynamicBoneParameters{}
		return nil
	}
	if count > 1<<16 {
		return fmt.Errorf("dynamic bone set has unreasonable parameter count %d", count)
	}

	before := r.CurrentPosition()
	if _, err := r.Seek(parametersOffset.ToSeek(), io.SeekStart); err != nil {
		return err
	}
	s.Parameters = make([]DynamicBoneParameters, count)
	if err := r.Decode(&s.Parameters, args...); err != nil {
		return err
	}
	_, err := r.Seek(before, io.SeekStart)
	return err
}

// DynamicBoneCollider is a sphere attached to a fighter joint.
type DynamicBoneCollider struct {
	Bone   int32
	Offset model.Vec3[float32]
	Radius float32
}

// Dynamics is the dynamic-bone configuration stored at offset 0x2C of a fighter root.
type Dynamics struct {
	BoneSets  []DynamicBoneSet
	Colliders []DynamicBoneCollider
}

func (d *Dynamics) BinRead(r *binread.Reader, args ...Args) error {
	var boneSetCount int32
	var boneSetsOffset Addr
	var colliderCount int32
	var collidersOffset Addr
	var animationTableOffset Addr

	if err := r.Decode(&boneSetCount, args...); err != nil {
		return err
	}
	if err := r.Decode(&boneSetsOffset, args...); err != nil {
		return err
	}
	if err := r.Decode(&colliderCount, args...); err != nil {
		return err
	}
	if err := r.Decode(&collidersOffset, args...); err != nil {
		return err
	}
	if err := r.Decode(&animationTableOffset, args...); err != nil {
		return err
	}

	if boneSetCount < 0 || boneSetCount > maxDynamicBoneSets {
		return fmt.Errorf("dynamic bone set count %d is outside [0, %d]", boneSetCount, maxDynamicBoneSets)
	}
	if colliderCount < 0 || colliderCount > 1<<16 {
		return fmt.Errorf("dynamic bone collider count %d is unreasonable", colliderCount)
	}

	d.BoneSets = []DynamicBoneSet{}
	if boneSetCount > 0 && boneSetsOffset != Addr(0x20) {
		before := r.CurrentPosition()
		if _, err := r.Seek(boneSetsOffset.ToSeek(), io.SeekStart); err != nil {
			return err
		}
		d.BoneSets = make([]DynamicBoneSet, boneSetCount)
		if err := r.Decode(&d.BoneSets, args...); err != nil {
			return err
		}
		if _, err := r.Seek(before, io.SeekStart); err != nil {
			return err
		}
	}

	d.Colliders = []DynamicBoneCollider{}
	if colliderCount > 0 && collidersOffset != Addr(0x20) {
		before := r.CurrentPosition()
		if _, err := r.Seek(collidersOffset.ToSeek(), io.SeekStart); err != nil {
			return err
		}
		d.Colliders = make([]DynamicBoneCollider, colliderCount)
		if err := r.Decode(&d.Colliders, args...); err != nil {
			return err
		}
		if _, err := r.Seek(before, io.SeekStart); err != nil {
			return err
		}
	}

	return nil
}
