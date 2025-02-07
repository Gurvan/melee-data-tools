package binread

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

const debug = false

// const debug = true

func printlnDebug(x ...interface{}) {
	if debug {
		fmt.Println(x...)
	}
}

type Args = map[string]interface{}

type Reader struct {
	io.ReadSeeker
}

func NewReader(r io.ReadSeeker) Reader {
	return Reader{r}
}

func (r *Reader) CurrentPosition() int64 {
	pos, _ := r.Seek(0, io.SeekCurrent)
	return pos
}

func (r *Reader) Peek(n int) ([]byte, error) {
	before := r.CurrentPosition()

	b := make([]byte, n)
	err := r.Decode(&b)
	if err != nil {
		return nil, err
	}
	_, err = r.Seek(before, io.SeekStart)
	return b, err
}

func (r *Reader) ReadAt(offset int64) error {
	return nil
}

func (r *Reader) Size() int64 {
	before := r.CurrentPosition()
	pos, _ := r.Seek(0, io.SeekEnd)
	r.Seek(before, io.SeekStart)
	return pos
}

type BinReader interface {
	BinRead(*Reader, ...Args) error
}

type AfterParser interface {
	AfterParse(*Reader, ...Args) error
}

func (r *Reader) Decode(data any, args ...Args) error {
	return unmarshal(r, data, args...)
}

func (r *Reader) DecodePeek(data any, args ...Args) error {
	before := r.CurrentPosition()

	err := r.Decode(data)
	if err != nil {
		return err
	}
	_, err = r.Seek(before, io.SeekStart)
	return err
}

func unmarshal(r *Reader, data any, args ...Args) error {
	rv := reflect.ValueOf(data)
	rt := reflect.TypeOf(data)

	var v reflect.Value
	var t reflect.Type
	switch {
	case rv.Kind() == reflect.Ptr:
		if rv.IsNil() {
			return nil
		}
		v = rv.Elem()
		t = rt.Elem()
	default:
		v = rv
		t = rt
	}

	if v.Kind() == reflect.Interface {
		v = v.Elem()
		t = rt.Elem()
	}

	var err error

	printlnDebug("Type:", reflect.TypeOf(data))

	interBinReader := reflect.TypeOf((*BinReader)(nil)).Elem()
	if reflect.TypeOf(data).Implements(interBinReader) {
		printlnDebug("Reader")
		err = any(data).(BinReader).BinRead(r, args...)
		if err != nil {
			return err
		}
	} else {
		switch v.Kind() {
		case reflect.Invalid:
			return nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
			printlnDebug("Other Type:", v.Type())
			err = binary.Read(r, binary.BigEndian, data)
			if err != nil {
				return err
			}
		case reflect.Struct:
			printlnDebug("Struct:", v.NumField())
			for i := 0; i < v.NumField(); i++ {
				if bit, ok := t.Field(i).Tag.Lookup("binread"); ok {
					if bit == "-" || bit == "ignore" {
						continue
					}
				}
				field := reflect.New(v.Field(i).Type()).Interface()
				err = unmarshal(r, field, args...)
				if err != nil {
					return err
				}
				printlnDebug("Field:", v.Type().Field(i).Name)
				if v.Field(i).CanSet() {
					v.Field(i).Set(reflect.ValueOf(field).Elem())
				}
			}
		case reflect.Array:
			printlnDebug("Array")
			for i := 0; i < v.Len(); i++ {
				item := reflect.New(v.Index(i).Type()).Interface()
				err = unmarshal(r, item, args...)
				if err != nil {
					return err
				}
				v.Index(i).Set(reflect.ValueOf(item).Elem())
			}
		case reflect.Slice:
			printlnDebug("Slice")
			if v.Len() > 0 {
				for i := 0; i < v.Len(); i++ {
					item := reflect.New(v.Index(i).Type()).Interface()
					err = unmarshal(r, item, args...)
					if err != nil {
						return err
					}
					v.Index(i).Set(reflect.ValueOf(item).Elem())
				}
			}
		default:
		}
	}

	if err != nil {
		return err
	}

	interAfterParser := reflect.TypeOf((*AfterParser)(nil)).Elem()
	if reflect.TypeOf(data).Implements(interAfterParser) {
		printlnDebug("AfterParsing")
		err = any(data).(AfterParser).AfterParse(r, args...)
		if err != nil {
			return err
		}
	}

	return nil
}
