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

type CollVertex struct {
    X float32
    Y float32
}

type CollData struct {
	Vertices SizedArray[CollVertex]
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
