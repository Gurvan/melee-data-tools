package lib

import (
	"fmt"
	"io"
	"reflect"

	"github.com/Gurvan/melee-data-tools/binread"
	"github.com/Gurvan/melee-data-tools/logger"
)

type Args = binread.Args

type Addr uint32

var _ binread.BinReader = (*Addr)(nil)

func (p Addr) String() string {
	return fmt.Sprintf("0x%X", uint32(p))
}

func (p *Addr) BinRead(r *binread.Reader, args ...Args) error {
	var v uint32
	err := r.Decode(&v, args...)
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

func (s NullTerminatedString) String() string {
	return string(s)
}

func (s *NullTerminatedString) BinRead(r *binread.Reader, _ ...Args) error {
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
	Offset   Addr
	ValuePtr *T
}

func (p *Ptr[T]) GetValue() T {
	if p.ValuePtr == nil {
		p.ValuePtr = new(T)
		logger.Warning.Printf("Try to deref pointer of type %T. Returning %T default value instead.", *p.ValuePtr, *p.ValuePtr)
	}
	return *p.ValuePtr
}

func (p *Ptr[T]) SetValue(v T) {
	p.ValuePtr = &v
}

func IsPtr(x interface{}) bool {
	t := reflect.TypeOf(x)
	return t.Kind() == reflect.Struct &&
		t.NumField() == 2 &&
		t.Field(0).Name == "Offset" &&
		t.Field(0).Type.Kind() == reflect.Uint32 &&
		t.Field(1).Name == "ValuePtr"
}

var _ binread.BinReader = (*Ptr[uint32])(nil)

func (p *Ptr[T]) BinRead(r *binread.Reader, args ...Args) error {
	seekFrom := io.SeekStart
	ok := false
	for _, args := range args {
		if p.Offset, ok = args["offset"].(Addr); ok {
			if seek, ok := args["seekfrom"].(int); ok {
				seekFrom = seek
			}
			break
		}
	}

	if !ok {
		err := r.Decode(&p.Offset)
		if err != nil {
			return err
		}
	}

	before := r.CurrentPosition()
	_, err := r.Seek(p.Offset.ToSeek(), seekFrom)
	if err != nil {
		return err
	}

	p.ValuePtr = new(T)

	err = r.Decode(p.ValuePtr, args...)
	if err != nil {
		return err
	}

	_, err = r.Seek(before, io.SeekStart)
	if err != nil {
		return err
	}
	return nil
}

type OptionalPtr[T any] Ptr[T]

func (p *OptionalPtr[T]) BinRead(r *binread.Reader, args ...Args) error {
	seekFrom := io.SeekStart
	ok := false
	for _, args := range args {
		if p.Offset, ok = args["offset"].(Addr); ok {
			if seek, ok := args["seekfrom"].(int); ok {
				seekFrom = seek
			}
			break
		}
	}

	if !ok {
		err := r.Decode(&p.Offset)
		if err != nil {
			return err
		}
	}

	if p.Offset == Addr(0x20) {
		return nil
	}

	before := r.CurrentPosition()
	_, err := r.Seek(p.Offset.ToSeek(), seekFrom)
	if err != nil {
		return err
	}

	p.ValuePtr = new(T)

	err = r.Decode(p.ValuePtr, args...)
	if err != nil {
		return err
	}

	_, err = r.Seek(before, io.SeekStart)
	if err != nil {
		return err
	}
	return nil
}

// type Array[T any] []T
//
// func (a *Array[T]) BinRead(r *binread.Reader, args ...Args) error {
// 	var err error
// 	var offset Addr
//
// 	err = r.Decode(&offset)
// 	if err != nil {
// 		return err
// 	}
//
// 	var numElem int
// 	for _, args := range args {
// 		if reloc, ok := args["relocation"].(*Relocation); ok && reloc != nil {
// 			// This section is probably incorrect.
// 			// We still need to make the relcation table work properly.
// 			var t T
// 			elemSize := unsafe.Sizeof(t)
// 			if IsPtr(t) {
// 				elemSize = 4
// 			}
// 			numElem = int((*reloc)[offset] / uint32(elemSize))
// 			fmt.Println("Num elem:", numElem)
// 		}
// 	}
//
// 	before := r.CurrentPosition()
// 	_, err = r.Seek(offset.ToSeek(), io.SeekStart)
// 	if err != nil {
// 		return err
// 	}
//
// 	elems := make([]T, 0)
// 	for i := 0; i < numElem; i++ {
// 		var t T
// 		err = r.Decode(&t, args...)
// 		if err != nil {
// 			return err
// 		}
// 		elems = append(elems, t)
// 	}
//
// 	_, err = r.Seek(before, io.SeekStart)
// 	if err != nil {
// 		return err
// 	}
//
// 	*a = elems
// 	return nil
// }

type SizedArray[T any] struct {
	Data []T
	Size uint32
}

func (a *SizedArray[T]) BinRead(r *binread.Reader, _ ...Args) error {
	var err error
	var offset Addr

	err = r.Decode(&offset)
	if err != nil {
		return err
	}

	if offset == Addr(0x20) {
		a.Size = 0
		a.Data = make([]T, a.Size)
		return nil
	}

	err = r.Decode(&a.Size)
	if err != nil {
		return err
	}

	before := r.CurrentPosition()
	_, err = r.Seek(offset.ToSeek(), io.SeekStart)
	if err != nil {
		return err
	}

	data := make([]T, a.Size)
	err = r.Decode(&data)
	if err != nil {
		return err
	}
	a.Data = data

	_, err = r.Seek(before, io.SeekStart)
	if err != nil {
		return err
	}

	return nil
}
