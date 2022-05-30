package descriptor

import (
	"errors"
	"fmt"
	"io"
	"sort"

	"github.com/Gurvan/melee-data-tools/binread"
	. "github.com/Gurvan/melee-data-tools/common"
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

type Relocation map[Addr]uint32

func (t *Relocation) BinRead(r *binread.Reader, args ...Args) error {
	reloc := make(map[Addr]uint32)

	var offset Addr
	var count uint32

	for _, args := range args {
		var ok bool
		if offset, ok = args["offset"].(Addr); !ok {
			return errors.New("Offset required for parsing relocation table.")
		}
		if count, ok = args["count"].(uint32); !ok {
			return errors.New("Count required for parsing relocation table.")
		}
	}

	_, err := r.Seek(offset.ToSeek(), io.SeekStart)
	if err != nil {
		return err
	}

	offsets := make([]Addr, count)

	var c uint32 = 0
	for {
		var v Ptr[Addr]
		err = r.Decode(&v)
		if err != nil {
			return err
		}
		offsets[c] = v.Value
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
		reloc[offset] = uint32(offsets[i+1] - offset)
	}

	*t = reloc
	return nil
}

type NamedOffset struct {
	Offset Addr
	Name   NullTerminatedString
}

var _ binread.BinReader = (*NamedOffset)(nil)

func (n *NamedOffset) BinRead(r *binread.Reader, args ...Args) error {
	var stringsOffset uint32

	for _, args := range args {
		if offset, ok := args["offset"].(uint32); ok {
			stringsOffset = offset
		}
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
	ptrArgs := Args{"offset": Addr(stringOffset + before + stringsOffset), "seekfrom": io.SeekStart}
	err = r.Decode(&namePtr, ptrArgs)
	if err != nil {
		return err
	}
	n.Name = namePtr.Value
	return nil
}

type Root = NamedOffset
type Ref = NamedOffset

type Footer struct {
	Relocation Relocation
	Roots      []Root
	Refs       []Ref
}

type Descriptor struct {
	Header
	Footer
}

func (d *Descriptor) BinRead(r *binread.Reader, _ ...Args) error {
	err := r.Decode(&d.Header)
	if err != nil {
		return err
	}

	_, err = r.Seek(d.RelocationOffset.ToSeek(), io.SeekStart)
	if err != nil {
		return err
	}

	relocArgs := Args{"offset": d.RelocationOffset, "count": d.RelocationCount}
	err = r.Decode(&d.Relocation, relocArgs)
	if err != nil {
		return err
	}

	rootRefOffset := 8 * (d.RootCount + d.RefCount)
	rootRefArgs := Args{"offset": rootRefOffset}
	d.Roots = make([]Root, d.RootCount)
	err = r.Decode(&d.Roots, rootRefArgs)
	if err != nil {
		return err
	}
	d.Refs = make([]Ref, d.RefCount)
	err = r.Decode(&d.Refs, rootRefArgs)
	if err != nil {
		return err
	}

	err = d.GoToFirstRoot(r)
	if err != nil {
		return err
	}

	return nil
}

func (d *Descriptor) GoToFirstRoot(r *binread.Reader) error {
	if len(d.Roots) == 0 {
		return errors.New("No roots in struct")
	}
	firstRootOffset := d.Roots[0].Offset
	_, err := r.Seek(firstRootOffset.ToSeek(), io.SeekStart)
	if err != nil {
		return err
	}
	return nil
}
