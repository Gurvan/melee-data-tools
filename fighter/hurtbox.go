package fighter

import (
	"io"

	"github.com/Gurvan/melee-data-tools/binread"
	. "github.com/Gurvan/melee-data-tools/common"
	"github.com/Gurvan/melee-data-tools/fighter/hurtbox"
)

type Hurtbox struct {
	BoneIndex uint32
	Position  hurtbox.Position
	Grabbable hurtbox.Grabbable
	BaseX     float32
	BaseY     float32
	BaseZ     float32
	TipX      float32
	TipY      float32
	TipZ      float32
	Radius    float32
}

type Hurtboxes []Hurtbox

var _ binread.BinReader = (*Hurtboxes)(nil)

func (t *Hurtboxes) BinRead(r *binread.Reader, _ ...Args) error {
	var count uint32

	err := r.Decode(&count)
	if err != nil {
		return err
	}

	// Follow Ptr
	before := r.CurrentPosition()

	var addr Addr
	err = r.Decode(&addr)
	if err != nil {
		return err
	}
	_, err = r.Seek(addr.ToSeek(), io.SeekStart)
	if err != nil {
		return err
	}

	// Parse hurtboxes
	hurtboxes := make([]Hurtbox, count)
	err = r.Decode(&hurtboxes)
	if err != nil {
		return err
	}

	*t = hurtboxes
	_, err = r.Seek(before, io.SeekStart)
	return err
}
