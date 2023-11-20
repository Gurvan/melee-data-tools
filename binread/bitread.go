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

func JoinBits(bits []Bit, padto int, signed bool) []byte {
	signBit := bits[0]
	byts := make([]byte, 0)
	if padto > len(bits) {
		padding := make([]Bit, padto-len(bits))
		if signed {
			for i := range padding {
				padding[i] = signBit
			}
		}
		bits = append(padding, bits...)
	} else {
		numbytes := (len(bits)-1)/8 + 1
		padding := make([]Bit, 8*numbytes-len(bits))
		if signed {
			for i := range padding {
				padding[i] = signBit
			}
		}
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
			n, _ := fmt.Sscanf(bit, "%d", &numbits)
			if n == 1 {
				totalbits += numbits
			}
		}
	}
	return totalbits, nil
}

func BitRead(bits []Bit, s any, startIndex int) error {
	readBits := func(value any, p, numbits int, signed bool) (any, error) {
		numbytes := reflect.TypeOf(value).Elem().Size()
		byts := JoinBits(bits[p:p+numbits], int(8*numbytes), signed)
		r := bytes.NewReader(byts)
		err := binary.Read(r, binary.BigEndian, value)
		return any(value), err
	}
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr {
		return errors.New("BitRead should be called with a pointer to a struct.")
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
			if bit == "-" || bit == "ignore" {
				continue
			}
			fmt.Sscanf(bit, "%d", &numbits)

			var value any
			var err error
			switch field.Type.Kind() {
			case reflect.Uint8:
				var valueLocal uint8
				value, err = readBits(&valueLocal, p, numbits, false)
			case reflect.Int8:
				var valueLocal int8
				value, err = readBits(&valueLocal, p, numbits, false)
			case reflect.Uint32:
				var valueLocal uint32
				value, err = readBits(&valueLocal, p, numbits, false)
			case reflect.Int32:
				var valueLocal int32
				value, err = readBits(&valueLocal, p, numbits, true)
			case reflect.Bool:
				var valueLocal bool
				value, err = readBits(&valueLocal, p, numbits, false)
			default:
				panic(fmt.Sprintf("Type %v not implemented for Bitread", field.Type))
			}
			if err != nil {
				panic(err)
			}
			if v.Field(i).CanSet() {
				// v.Field(i).Set(reflect.ValueOf(value).Elem().Convert(field.Type))
				v.Field(i).Set(reflect.ValueOf(value).Elem().Convert(field.Type))
			}
			p += numbits
		} else {
			logger.Error.Printf("All fields should have a `bit` number tag. Struct: %v | Field: %v", t.Name(), field.Name)
			return errors.New(fmt.Sprintf("All fields should have a `bit` number tag. Struct: %v | Field: %v", t.Name(), field.Name))
		}
	}
	return nil
}
