package binread

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"testing"
)

type CustomUInt16 uint16

func (s *CustomUInt16) BinRead(r io.ReadSeeker) error {
	err := binary.Read(r, binary.BigEndian, &s)
	if err != nil {
		return err
	}
	return nil
}

type A struct {
	V CustomUInt16
}

var _ BinReader = (*A)(nil)

type B struct {
	A [2]A
}

func (s *A) BinRead(r *Reader, _ ...interface{}) error {
	err := binary.Read(r, binary.BigEndian, &s.V)
	if err != nil {
		return err
	}
	return nil
}

func TestBinRead(t *testing.T) {
	// data := []byte{0x09, 0x0, 0x2, 0x1}
	// data := []byte{0x0, 0x0, 0x0, 0x09, 0x06}
	// reader := Reader{r: bytes.NewReader(data)}

	// s := A{}
	// err := reader.Decode(&s)
	// if err != nil {
	//         fmt.Println("Err:", err)
	// }
	// fmt.Println(s)

	// var u uint8 = 0
	// reader.Decode(&u)
	// fmt.Println(u)

	data := []byte{0x0, 0x0, 0x0, 0x09, 0x06}
	reader := Reader{bytes.NewReader(data)}

	w := B{}
	err := reader.Decode(&w)
	if err != nil {
		fmt.Println("Err:", err)
	}
	fmt.Println(w)
}
