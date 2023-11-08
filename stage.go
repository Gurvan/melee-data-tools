package mdt

import (
	"errors"
	"fmt"
	"io"

	"github.com/Gurvan/melee-data-tools/binread"
	"github.com/Gurvan/melee-data-tools/descriptor"
	. "github.com/Gurvan/melee-data-tools/lib"
)

type CollDirectionFlags struct {
	Disabled bool `bit:"1"`
	Left     bool `bit:"1"`
	Right    bool `bit:"1"`
	Bottom   bool `bit:"1"`
	Top      bool `bit:"1"`
}

func (f *CollDirectionFlags) BinRead(r *binread.Reader, _ ...Args) error {
	var err error

	byts := make([]byte, 2)
	err = r.Decode(&byts)
	if err != nil {
		return err
	}

	bits := binread.SplitBytes(byts)
	err = binread.BitRead(bits, f, 11)
	return err
}

type CollPropertyFlags struct {
	LedgeGrab   bool `bit:"1"`
	DropThrough bool `bit:"1"`
}

func (f *CollPropertyFlags) BinRead(r *binread.Reader, _ ...Args) error {
	var err error

	byts := make([]byte, 1)
	err = r.Decode(&byts)
	if err != nil {
		return err
	}

	bits := binread.SplitBytes(byts)
	err = binread.BitRead(bits, f, 6)
	return err
}

type CollVertex struct {
	X float32
	Y float32
}

type CollLine struct {
	VertexIndex1 int16
	VertexIndex2 int16

	NextLine     int16
	PreviousLine int16

	NextLineAltGroup     int16
	PreviousLineAltGroup int16

	DirectionFlags CollDirectionFlags
	PropertyFlags  CollPropertyFlags
	Material       byte
}

type CollLineGroup struct {
	TopLineIndex int16
	TopLineCount int16

	BottomLineIndex int16
	BottomLineCount int16

	RightLineIndex int16
	RightLineCount int16

	LeftLineIndex int16
	LeftLineCount int16

	DynamicLineIndex int16
	DynamicLineCount int16

	XMin float32
	YMin float32
	XMax float32
	YMax float32

	VertexStart int16
	VertexCount int16
}

type CollData struct {
	Vertices SizedArray[CollVertex]
	Links    SizedArray[CollLine]

	TopLinksOffset int16
	TopLinksCount  int16

	BottomLinksOffset int16
	BottomLinksCount  int16

	RightLinksOffset int16
	RightLinksCount  int16

	LeftLinksOffset int16
	LeftLinksCount  int16

	DynamicLinksOffset int16
	DynamicLinksCount  int16

	LineGroups SizedArray[CollLineGroup]
}

type PointType int16

const (
	Player1Spawn         PointType = 0
	Player2Spawn         PointType = 1
	Player3Spawn         PointType = 2
	Player4Spawn         PointType = 3
	Player1Respawn       PointType = 4
	Player2Respawn       PointType = 5
	Player3Respawn       PointType = 6
	Player4Respawn       PointType = 7
	ItemSpawn1           PointType = 127
	ItemSpawn2           PointType = 128
	ItemSpawn3           PointType = 129
	ItemSpawn4           PointType = 130
	ItemSpawn5           PointType = 131
	ItemSpawn6           PointType = 132
	ItemSpawn7           PointType = 133
	ItemSpawn8           PointType = 134
	ItemSpawn9           PointType = 135
	ItemSpawn10          PointType = 136
	ItemSpawn11          PointType = 137
	ItemSpawn12          PointType = 138
	ItemSpawn13          PointType = 139
	ItemSpawn14          PointType = 140
	ItemSpawn15          PointType = 141
	ItemSpawn16          PointType = 142
	ItemSpawn17          PointType = 143
	ItemSpawn18          PointType = 144
	ItemSpawn19          PointType = 145
	ItemSpawn20          PointType = 146
	DeltaAngleCamera     PointType = 148
	TopLeftBoundary      PointType = 149
	BottomRightBoundary  PointType = 150
	TopLeftBlastZone     PointType = 151
	BottomRightBlastZone PointType = 152
	Target1              PointType = 199
	Target2              PointType = 200
	Target3              PointType = 201
	Target4              PointType = 202
	Target5              PointType = 203
	Target6              PointType = 204
	Target7              PointType = 205
	Target8              PointType = 206
	Target9              PointType = 207
	Target10             PointType = 208
	Bumper1              PointType = 252
	Bumper2              PointType = 253
	Bumper3              PointType = 254
	Bumper4              PointType = 255
	Bumper5              PointType = 256
	Bumper6              PointType = 257
	Bumper7              PointType = 258
	Bumper8              PointType = 259
)

func (t PointType) String() string {
	switch t {
	case Player1Spawn:
		return "Player1Spawn"
	case Player2Spawn:
		return "Player2Spawn"
	case Player3Spawn:
		return "Player3Spawn"
	case Player4Spawn:
		return "Player4Spawn"
	case Player1Respawn:
		return "Player1Respawn"
	case Player2Respawn:
		return "Player2Respawn"
	case Player3Respawn:
		return "Player3Respawn"
	case Player4Respawn:
		return "Player4Respawn"
	case ItemSpawn1:
		return "ItemSpawn1"
	case ItemSpawn2:
		return "ItemSpawn2"
	case ItemSpawn3:
		return "ItemSpawn3"
	case ItemSpawn4:
		return "ItemSpawn4"
	case ItemSpawn5:
		return "ItemSpawn5"
	case ItemSpawn6:
		return "ItemSpawn6"
	case ItemSpawn7:
		return "ItemSpawn7"
	case ItemSpawn8:
		return "ItemSpawn8"
	case ItemSpawn9:
		return "ItemSpawn9"
	case ItemSpawn10:
		return "ItemSpawn10"
	case ItemSpawn11:
		return "ItemSpawn11"
	case ItemSpawn12:
		return "ItemSpawn12"
	case ItemSpawn13:
		return "ItemSpawn13"
	case ItemSpawn14:
		return "ItemSpawn14"
	case ItemSpawn15:
		return "ItemSpawn15"
	case ItemSpawn16:
		return "ItemSpawn16"
	case ItemSpawn17:
		return "ItemSpawn17"
	case ItemSpawn18:
		return "ItemSpawn18"
	case ItemSpawn19:
		return "ItemSpawn19"
	case ItemSpawn20:
		return "ItemSpawn20"
	case DeltaAngleCamera:
		return "DeltaAngleCamera"
	case TopLeftBoundary:
		return "TopLeftBoundary"
	case BottomRightBoundary:
		return "BottomRightBoundary"
	case TopLeftBlastZone:
		return "TopLeftBlastZone"
	case BottomRightBlastZone:
		return "BottomRightBlastZone"
	case Target1:
		return "Target1"
	case Target2:
		return "Target2"
	case Target3:
		return "Target3"
	case Target4:
		return "Target4"
	case Target5:
		return "Target5"
	case Target6:
		return "Target6"
	case Target7:
		return "Target7"
	case Target8:
		return "Target8"
	case Target9:
		return "Target9"
	case Target10:
		return "Target10"
	case Bumper1:
		return "Bumper1"
	case Bumper2:
		return "Bumper2"
	case Bumper3:
		return "Bumper3"
	case Bumper4:
		return "Bumper4"
	case Bumper5:
		return "Bumper5"
	case Bumper6:
		return "Bumper6"
	case Bumper7:
		return "Bumper7"
	case Bumper8:
		return "Bumper8"
	default:
		return fmt.Sprintf("PointTypeUnimplemented(%d)", t)
	}
}

type GeneralPointInfo struct {
	Index int16
	Type  PointType
}

type Joint struct {
	_            [2]uint32
	Child        OptionalPtr[Joint]
	Next         OptionalPtr[Joint]
	_            [4]byte
	_            [6]uint32
	TranslationX float32
	TranslationY float32
	_            [8]byte
}

type GeneralPoints struct {
	ReferenceJoint Ptr[Joint]
	Points         SizedArray[GeneralPointInfo]
}

func (g *GeneralPoints) FindCoordinates(t PointType) (float32, float32, error) {
	var index int16 = -1
	for _, p := range g.Points.Data {
		if p.Type == t {
			index = p.Index
			break
		}
	}
	if index < 0 {
		return 0.0, 0.0, errors.New(fmt.Sprintf("Couldn't find point type %s.", t))
	}
	joint := g.ReferenceJoint.GetValue().Child.ValuePtr
	for i := int16(0); i < index-1; i++ {
		if joint == nil {
			return 0.0, 0.0, errors.New(fmt.Sprintf("Couldn't find point type %s coordinates.", t))
		}
		joint = joint.Next.ValuePtr
	}
	return joint.TranslationX, joint.TranslationY, nil
}

type MapHead struct {
	GeneralPoints SizedArray[GeneralPoints]
}


type GroundParam struct {
    StageScale float32
    _ [216]byte
}

type StageData struct {
	CollData CollData
	MapHead  MapHead
    GroundParam GroundParam
}

type StageFile = File[StageData]

func (s *StageData) BinRead(r *binread.Reader, args ...Args) error {
	var err error
	var collDataOffset Addr
	var mapHeadOffset Addr
	var groundParamOffset Addr
	for _, args := range args {
		if desc, ok := args["descriptor"].(*descriptor.Descriptor); ok {
			if collDataOffset, err = desc.FindRootOffset("coll_data"); err != nil {
				return err
			}
			if mapHeadOffset, err = desc.FindRootOffset("map_head"); err != nil {
				return err
			}
			if groundParamOffset, err = desc.FindRootOffset("grGroundParam"); err != nil {
				return err
			}
			break
		}
	}

	_, err = r.Seek(collDataOffset.ToSeek(), io.SeekStart)
	if err != nil {
		return err
	}
	err = r.Decode(&s.CollData)
	if err != nil {
		return err
	}

	_, err = r.Seek(mapHeadOffset.ToSeek(), io.SeekStart)
	if err != nil {
		return err
	}
	err = r.Decode(&s.MapHead)
	if err != nil {
		return err
	}

	_, err = r.Seek(groundParamOffset.ToSeek(), io.SeekStart)
	if err != nil {
		return err
	}
	err = r.Decode(&s.GroundParam)
	if err != nil {
		return err
	}

	return nil
}
