package mdt

import (
	"io"

	"github.com/Gurvan/melee-data-tools/binread"
	"github.com/Gurvan/melee-data-tools/descriptor"
	. "github.com/Gurvan/melee-data-tools/lib"
)

type StageData struct {
	CollData CollData
}

type StageFile = File[StageData]

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

func (s *StageData) BinRead(r *binread.Reader, args ...Args) error {
	var err error
	var collDataOffset Addr
	for _, args := range args {
		if desc, ok := args["descriptor"].(*descriptor.Descriptor); ok {
			if collDataOffset, err = desc.FindRootOffset("coll_data"); err != nil {
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

	return nil
}
