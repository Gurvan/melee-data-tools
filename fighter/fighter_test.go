package fighter

import (
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestSubAction(t *testing.T) {
	t.Logf("%v\n", spew.Sdump(subactionTypeSwitch(0)))
	t.Logf("%v\n", spew.Sdump(subactionTypeSwitch(1)))

	// args := map[string]uint32{"a": 1}
	// v, ok := args["a"]
	// t.Logf("%v | %v\n", v, ok)
	// v, ok = args["b"]
	// t.Logf("%v | %v\n", v, ok)

	s := SynchronousTimer{}
	// fmt.Println(s)
	// BitRead(&s)
	// fmt.Println(s)

	f := uint32(1141374976)
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, f)
	// binary.LittleEndian.PutUint32(bs, f)

	fmt.Println(bs)
	fmt.Println(SplitBytes(bs))
	fmt.Println(JoinBits(SplitBytes(bs)))

	fmt.Println(NumBits(s))
	fmt.Println(NumBits(&s))

	// fmt.Println(DecodeSubAction(&s))
}
