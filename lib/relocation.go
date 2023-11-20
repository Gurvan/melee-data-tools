package lib

import (
	"errors"
	"io"
	"sort"

	"github.com/Gurvan/melee-data-tools/binread"
)

type Relocation map[Addr]uint32

func (t *Relocation) BinRead(r *binread.Reader, args ...Args) error {
	reloc := make(map[Addr]uint32)

	var offset Addr
	var count uint32

	for _, args := range args {
		var ok bool
		if offset, ok = args["offset"].(Addr); !ok {
			return errors.New("offset required for parsing relocation table")
		}
		if count, ok = args["count"].(uint32); !ok {
			return errors.New("count required for parsing relocation table")
		}
	}

	_, err := r.Seek(offset.ToSeek(), io.SeekStart)
	if err != nil {
		return err
	}

	offsets := make([]Addr, 0)

	var c uint32 = 0
	for {
		var v Ptr[Addr]
		err = r.Decode(&v)
		if err != nil {
			return err
		}
		offsets = append(offsets, v.GetValue())
		c++
		if c >= count {
			break
		}
	}

	sort.Slice(offsets, func(i, j int) bool { return offsets[i] < offsets[j] })
	for i, offset := range offsets {
		if i == len(offsets)-1 {
			reloc[offset] = 0
			break
		}
		s := uint32(offsets[i+1] - offset)
		if s > 0 {
			reloc[offset] = s
		}
	}

	*t = reloc
	return nil
}
