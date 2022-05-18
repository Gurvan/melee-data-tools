package main

import (
	"fmt"
	"io"

	"github.com/Gurvan/melee-data-tools/binread"
)

type Addr uint32

var _ binread.BinReader = (*Addr)(nil)

func (p Addr) String() string {
	return fmt.Sprintf("0x%X", uint32(p))
}

func (p *Addr) BinRead(r *binread.Reader, _ ...interface{}) error {
	var v uint32
	err := r.Decode(&v)
	if err != nil {
		return err
	}
	*p = Addr(v + 0x20)
	return nil
}

func (p *Addr) ToSeek() int64 {
	return int64(*p)
}

func (p *Addr) Add(x uint32) *Addr {
	a := *p + Addr(x)
	return &a
}

type NullTerminatedString string

var _ binread.BinReader = (*NullTerminatedString)(nil)

func (s *NullTerminatedString) BinRead(r *binread.Reader, _ ...interface{}) error {
	bs := make([]byte, 0)
	for {
		var b byte
		err := r.Decode(&b)
		if err != nil {
			return err
		}

		if b == 0x00 {
			break
		}

		bs = append(bs, b)
	}
	*s = NullTerminatedString(string(bs))
	return nil
}

type Ptr[T any] struct {
	Offset Addr
	Value  T
}

var _ binread.BinReader = (*Ptr[uint32])(nil)

func (p *Ptr[T]) BinRead(r *binread.Reader, args ...interface{}) error {
	fmt.Println("Args:", args)
	fmt.Println(args[0].(Addr))
	seekFrom := io.SeekStart
	if offset, ok := args[0].(Addr); ok {
		p.Offset = offset
		fmt.Println(offset)
		if seek, ok := args[1].(int); ok {
			seekFrom = seek
			fmt.Println("SeekFrom:", seek)
		}
	} else {
		err := r.Decode(&p.Offset)
		if err != nil {
			return err
		}
	}

	before, err := r.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}

	_, err = r.Seek(p.Offset.ToSeek(), seekFrom)
	if err != nil {
		return err
	}

	err = r.Decode(&p.Value)

	if err != nil {
		return err
	}

	_, err = r.Seek(before, io.SeekStart)
	if err != nil {
		return err
	}
	return nil
}
