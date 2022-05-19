package fighter

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/Gurvan/melee-data-tools/binread"
)

type emptySubAction struct{}

type SubActionNotImplemented struct {
	Id uint8
}

func (e *SubActionNotImplemented) Error() string {
	return fmt.Sprintf("SubAction 0x%X is not implemented", e.Id)
}

// func (s emptySubAction) IsSubAction() bool {
//         return true
// }

type SubAction interface {
	// IsSubAction() bool
}

var _ SubAction = (*emptySubAction)(nil)

func SubActionString(s SubAction) string {
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	str := t.Name() + "{"
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Struct {
			continue
		}
		str += field.Name + ": " + fmt.Sprintf("%v", v.Field(i).Interface())
		if i < t.NumField()-1 {
			str += ", "
		}
	}

	str += "}"
	return str
}

// 0x00
type EndOfScript struct {
	emptySubAction
	_ uint32 `bit:"26"`
}

// 0x01
type SynchronousTimer struct {
	emptySubAction
	Frame uint32 `bit:"26"`
}

// 0x02
type AsynchronousTimer struct {
	emptySubAction
	Frame uint32 `bit:"26"`
}

// 0x13
type SetFlag struct {
	emptySubAction
	Flag  uint32 `bit:"2"`
	Value uint32 `bit:"24"`
}

func subactionTypeSwitch(i uint8) (SubAction, error) {
	switch i {
	case 0x00:
		return &EndOfScript{}, nil
	case 0x01:
		return &SynchronousTimer{}, nil
	case 0x02:
		return &AsynchronousTimer{}, nil
	case 0x13:
		return &SetFlag{}, nil
	default:
		return &emptySubAction{}, &SubActionNotImplemented{Id: i}
	}
}

func GetSubActionType(r *binread.Reader) (SubAction, error) {
	b, err := r.Peek(1)
	if err != nil {
		return emptySubAction{}, err
	}
	return subactionTypeSwitch(uint8(b[0]) >> 2)
}

type Bit = byte

func SplitBytes(byts []byte) []Bit {
	bits := make([]Bit, 0)
	for _, b := range byts {
		for i := 0; i < 8; i++ {
			v := 0
			if int(b)&(2<<i) > 0 {
				v = 1
			}
			bits = append(bits, byte(v))
		}
	}
	return bits
}

func JoinBits(bits []Bit, padto int) []byte {
	byts := make([]byte, 0)
	if padto > len(bits) {
		padding := make([]Bit, padto-len(bits))
		bits = append(padding, bits...)
	} else {
		numbytes := (len(bits)-1)/8 + 1
		padding := make([]Bit, 8*numbytes-len(bits))
		bits = append(padding, bits...)
	}

	var by int
	for i, b := range bits {
		i = i % 8
		if i == 0 {
			by = 0
		}
		by += int(b) * (2 << i)
		if i == 7 {
			byts = append(byts, byte(by))
		}
	}
	return byts
}

func NumBits(s SubAction) int {
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	var totalbits int = 6
	var numbits int
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Struct {
			continue
		}
		if bit, ok := field.Tag.Lookup("bit"); ok {
			fmt.Sscanf(bit, "%d", &numbits)
			totalbits += numbits
			// fmt.Println(field.Name, bit, numbits, totalbits)
			// if v.Field(i).CanSet() {
			//         v.Field(i).Set(reflect.ValueOf(uint32(numbits)))
			// }
		}
	}
	return totalbits
}

func BitRead(bits []Bit, s SubAction) error {
	v := reflect.ValueOf(s)
	// pt := reflect.TypeOf(s)
	// if pv.Kind() != reflect.Ptr {
	//         return errors.New("BitRead should be called with a pointer to a SubAction.")
	//         // v = v.Elem()
	//         // t = t.Elem()
	// }
	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}
	// v := pv.Elem()
	// t := pt.Elem()
	if v.Kind() != reflect.Struct {
		log.Fatal("unexpected type")
	}
	t := v.Type()
	fmt.Println(v)

	var p int = 6
	var numbits int
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Struct {
			continue
		}
		if bit, ok := field.Tag.Lookup("bit"); ok {
			fmt.Sscanf(bit, "%d", &numbits)
			// fmt.Println(field.Name, bit, numbits)
			switch t.Field(i).Type.Kind() {
			case reflect.Uint32:
				byts := JoinBits(bits[p:p+numbits], 32)
				fmt.Println(bits, bits[p:p+numbits], p, p+numbits, byts)
				r := bytes.NewReader(byts)
				var value uint32
				err := binary.Read(r, binary.BigEndian, &value)
				if err != nil {
					return err
				}
				fmt.Printf("Type %T, Value: %v\n", value, value)
				fmt.Println(t.Field(i).Name)
				if t.Field(i).Name != "_" {
					v.Field(i).Set(reflect.ValueOf(value))
				}
				if v.Field(i).CanSet() {
					fmt.Println("Setting:", value)
					v.Field(i).Set(reflect.ValueOf(value))
				}
			default:
			}
			p += numbits
		} else {
			return errors.New(fmt.Sprintf("All SubAction fields should have a `bit` number tag. SubAction: %v | Field: %v", t.Name(), field.Name))
		}
	}
	return nil
}

func DecodeSubAction(r *binread.Reader, s SubAction) error {
	// fmt.Printf("%#+v\n", s)
	numbits := NumBits(s)
	if numbits%32 != 0 {
		return errors.New(fmt.Sprintf("SubAction total bit number +6 should be a multiple of 32. Numbits+6=%d", numbits))
	}

	byts := make([]byte, numbits/8)
	// fmt.Println(numbits/8, r.CurrentPosition())
	err := r.Decode(&byts)
	if err != nil {
		return err
	}

	bits := SplitBytes(byts)
	return BitRead(bits, s)
}
