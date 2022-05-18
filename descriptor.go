package main

import (
	"fmt"
	"io"

	"github.com/Gurvan/melee-data-tools/binread"
)

type Version [4]byte

func (v Version) String() string {
	var b [4]byte = v
	return fmt.Sprintf("%c%c%c%c", b[0], b[1], b[2], b[3])
}

type Header struct {
	FileSize         uint32
	RelocationOffset Addr
	RelocationCount  uint32
	RootCount        uint32
	RefCount         uint32
	Version          Version
}

func (h *Header) GetStringsOffset() *Addr {
	return h.RelocationOffset.Add(4*h.RelocationCount + 8*(h.RootCount+h.RefCount))
}

type NamedOffset struct {
	Offset Addr
	Name   NullTerminatedString
}

var _ binread.BinReader = (*NamedOffset)(nil)

func (n *NamedOffset) BinRead(r *binread.Reader, args ...interface{}) error {
	var stringsOffset uint32
	if offset, ok := args[0].(uint32); ok {
		stringsOffset = offset
	}

	before := uint32(r.CurrentPosition())

	err := r.Decode(&n.Offset)
	if err != nil {
		return err
	}

	var stringOffset uint32
	err = r.Decode(&stringOffset)
	if err != nil {
		return err
	}

	var namePtr Ptr[NullTerminatedString]
	err = r.Decode(&namePtr, Addr(stringOffset+before+stringsOffset), io.SeekStart)
	if err != nil {
		return err
	}
	n.Name = namePtr.Value
	return nil
}

type Root = NamedOffset
type Ref = NamedOffset

type Footer struct {
	Relocation Addr
	Roots      []Root
	Refs       []Ref
}

type Descriptor struct {
	Header
	Footer
}

func (d *Descriptor) BinRead(r *binread.Reader, _ ...interface{}) error {
	err := r.Decode(&d.Header)
	if err != nil {
		return err
	}

	_, err = r.Seek(d.RelocationOffset.ToSeek(), io.SeekStart)
	if err != nil {
		return err
	}

	refRootsOffset := 8 * (d.RootCount + d.RefCount)

	reloc := make([]Addr, d.RelocationCount)
	err = r.Decode(&reloc)
	if err != nil {
		return err
	}

	d.Roots = make([]Root, d.RootCount)
	err = r.Decode(&d.Roots, refRootsOffset)
	if err != nil {
		return err
	}

	d.Refs = make([]Ref, d.RefCount)
	err = r.Decode(&d.Refs, refRootsOffset)
	if err != nil {
		return err
	}

	return nil
}
