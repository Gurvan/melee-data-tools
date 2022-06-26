package binread

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"

	"github.com/Gurvan/melee-data-tools/logger"
)

type Bit = byte

func SplitBytes(byts []byte) []Bit {
	bits := make([]Bit, 0)
	for _, b := range byts {
		for i := 0; i < 8; i++ {
			v := 0
			if b&(1<<(7-i)) > 0 {
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
		by += int(b) * (1 << (7 - i))
		if i == 7 {
			byts = append(byts, byte(by))
		}
	}
	return byts
}

func NumBits(s any) (int, error) {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr {
		return 0, errors.New("BitRead should be called with a pointer to a struct.")
	}

	v = v.Elem()
	t := reflect.TypeOf(s).Elem()
	var totalbits int = 6
	var numbits int
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous {
			continue
		}
		if bit, ok := field.Tag.Lookup("bit"); ok {
			fmt.Sscanf(bit, "%d", &numbits)
			totalbits += numbits
		}
	}
	return totalbits, nil
}

func BitRead(bits []Bit, s any, startIndex int) error {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr {
		return errors.New("BitRead should be called with a pointer to a SubAction.")
	}

	v = v.Elem()
	t := reflect.TypeOf(s).Elem()

	var p int = startIndex
	var numbits int
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous {
			continue
		}
		if bit, ok := field.Tag.Lookup("bit"); ok {
			fmt.Sscanf(bit, "%d", &numbits)
			byts := JoinBits(bits[p:p+numbits], 32)
			r := bytes.NewReader(byts)

			var value any
			var err error
			switch field.Type.Kind() {
			case reflect.Uint32:
				var valueLocal uint32
				err = binary.Read(r, binary.BigEndian, &valueLocal)
				value = any(valueLocal)
			case reflect.Int32:
				var valueLocal int32
				err = binary.Read(r, binary.BigEndian, &valueLocal)
				value = any(valueLocal)
			case reflect.Bool:
				var valueLocal bool
				err = binary.Read(r, binary.BigEndian, &valueLocal)
				value = any(valueLocal)
			}
			if err != nil {
				return err
			}
			if v.Field(i).CanSet() {
				v.Field(i).Set(reflect.ValueOf(value).Convert(field.Type))
			}
			p += numbits
		} else {
			logger.Error.Printf("All fields should have a `bit` number tag. Struct: %v | Field: %v", t.Name(), field.Name)
			return errors.New(fmt.Sprintf("All fields should have a `bit` number tag. Struct: %v | Field: %v", t.Name(), field.Name))
		}
	}
	return nil
}
