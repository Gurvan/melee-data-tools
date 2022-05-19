package fighter

import (
	"encoding/binary"
	"fmt"
	"reflect"
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
	// fmt.Println(JoinBits(SplitBytes(bs)))

	fmt.Println(NumBits(s))
	fmt.Println(NumBits(&s))

	// fmt.Println(DecodeSubAction(&s))
}

func TestInterface(t *testing.T) {
	type Sub interface{}

	type Timer struct {
		Frame uint32
	}

	var _ Sub = (*Timer)(nil)

	MakeTimer := func() Sub {
		return &Timer{}
	}

	SetFrame := func(s Sub) {
		v := reflect.ValueOf(s)
		fmt.Println(v.Kind())
		v = v.Elem()
		fmt.Println(v.Kind())
		t := reflect.TypeOf(s).Elem()
		fmt.Println(t)
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fmt.Printf("%#v\n", field)
			fmt.Println(v.Field(i).CanSet())
			v.Field(i).Set(reflect.ValueOf(uint32(3)))
		}

	}

	// timer := &Timer{}
	timer := MakeTimer()
	SetFrame(timer)
	fmt.Println(timer)
}
