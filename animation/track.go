package animation

import (
	"github.com/Gurvan/melee-data-tools/animation/track"
	"github.com/Gurvan/melee-data-tools/binread"
	. "github.com/Gurvan/melee-data-tools/common"
)

type Track struct {
	DataLength    uint16
	StartFrame    int16
	Type          track.Type
	ValueFormat   track.DataFormat
	ValueScale    uint32
	TangentFormat track.DataFormat
	TangentScale  uint32
	Keys          track.Keys
}

func (t *Track) BinRead(r *binread.Reader, _ ...Args) error {
	var err error
	err = r.Decode(&t.DataLength)
	if err != nil {
		return err
	}

	err = r.Decode(&t.StartFrame)
	if err != nil {
		return err
	}

	err = r.Decode(&t.Type)
	if err != nil {
		return err
	}

	var valueFlags uint8
	err = r.Decode(&valueFlags)
	if err != nil {
		return err
	}

	var tangentFlags uint8
	err = r.Decode(&tangentFlags)
	if err != nil {
		return err
	}

	t.ValueFormat = track.DataFormat(valueFlags & 0b11100000)
	t.ValueScale = 1 << (valueFlags & 0b00011111)
	t.TangentFormat = track.DataFormat(tangentFlags & 0b11100000)
	t.TangentScale = 1 << (tangentFlags & 0b00011111)

	var unused byte
	err = r.Decode(&unused)
	if err != nil {
		return err
	}

	args := Args{
		"valueScale":    t.ValueScale,
		"valueFormat":   t.ValueFormat,
		"tangentScale":  t.TangentScale,
		"tangentFormat": t.TangentFormat,
		"dataLength":    t.DataLength,
	}
	err = r.Decode(&t.Keys, args)
	if err != nil {
		return err
	}

	return nil
}

type Tracks [][]Track

func (t *Tracks) BinRead(r *binread.Reader, args ...Args) error {
	return nil
}
