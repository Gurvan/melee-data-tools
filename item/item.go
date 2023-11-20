package item

import (
	"io"

	"github.com/Gurvan/melee-data-tools/binread"
	. "github.com/Gurvan/melee-data-tools/lib"
	"github.com/Gurvan/melee-data-tools/model"
	"github.com/Gurvan/melee-data-tools/subaction"
)

type Item struct {
	Attributes Ptr[Attributes]
	_          Addr
	Hurtboxes  SizedArray[Hurtbox]
	States     States
	Model      Ptr[Model]
	_          Addr // Dynamics
}

type ItemFlags struct {
	IsHeavy  bool  `bit:"1"`
	_        uint8 `bit:"4"`
	HoldKind uint8 `bit:"7"`
}

func (f *ItemFlags) BinRead(r *binread.Reader, args ...Args) error {
	var err error

	byts := make([]byte, 4)
	err = r.Decode(&byts)
	if err != nil {
		return err
	}

	bits := binread.SplitBytes(byts)
	err = binread.BitRead(bits, f, 0)
	return err
}

type ECB struct {
	Top    float32
	Bottom float32
	Left   float32
	Right  float32
}

type Attributes struct {
	Flags               ItemFlags
	ThrowSpeedMultipler float32
	_                   float32
	SpinSpeed           float32
	FallAcceleration    float32
	MaxFallSpeed        float32
	_                   [10]float32 // Unk0x18
	ECB                 ECB
	_                   [4]float32 // Unk0x50
	ModelScale          float32
	_                   [8]int32 // SFXs
}

type Hurtbox struct {
	BoneIndex uint32
	BaseX     float32
	BaseY     float32
	BaseZ     float32
	TipX      float32
	TipY      float32
	TipZ      float32
	Radius    float32
}

const StateSize = 0x10

type State struct {
	_          [3]Addr
	Subactions subaction.SubActions
}

type States []State

func (s *States) BinRead(r *binread.Reader, args ...Args) error {
	var err error
	var offset Addr

	err = r.Decode(&offset)
	if err != nil {
		return err
	}

	var numElem int
	for _, args := range args {
		if reloc, ok := args["relocation"].(*Relocation); ok && reloc != nil {
			numElem = int((*reloc)[offset]) / StateSize
			// fmt.Println("Num elem states:", numElem, elemSize, offset, statesSize, ok)
		}
	}
	before := r.CurrentPosition()
	_, err = r.Seek(offset.ToSeek(), io.SeekStart)
	if err != nil {
		return err
	}

	states := make([]State, 0)
	for i := 0; i < numElem; i++ {
		var state State
		err = r.Decode(&state)
		if err != nil {
			return err
		}
		states = append(states, state)
	}
	_, err = r.Seek(before, io.SeekStart)
	if err != nil {
		return err
	}

	*s = states
	return nil
}

type Model struct {
	Joint        Ptr[model.Joint]
	BoneCount    int32
	BoneAttachID int32
	BitField     int32
}
