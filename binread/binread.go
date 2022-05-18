package binread

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"
)

const debug = false

func printDebug(x ...interface{}) {
	if debug {
		fmt.Println(x...)
	}
}

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

type BinReader interface {
	BinRead(*Reader, ...interface{}) error
}

func (r *Reader) Decode(data any, args ...interface{}) error {
	return unmarshal(r, data, args...)
}

func unmarshal(r *Reader, data any, args ...interface{}) error {
	rv := reflect.ValueOf(data)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("Not a pointer")
	}

	var err error

	// fmt.Println("Type:", reflect.TypeOf(data))
	printDebug("Type:", reflect.TypeOf(data))

	inter := reflect.TypeOf((*BinReader)(nil)).Elem()
	if reflect.TypeOf(data).Implements(inter) {
		// fmt.Println("Reader")
		printDebug("Reader")
		err = any(data).(BinReader).BinRead(r, args...)
	} else {
		v := rv.Elem()
		switch v.Kind() {
		case reflect.Struct:
			// fmt.Println("Struct")
			printDebug("Struct")
			for i := 0; i < v.NumField(); i++ {
				field := reflect.New(v.Field(i).Type()).Interface()
				err = unmarshal(r, field, args...)
				v.Field(i).Set(reflect.ValueOf(field).Elem())
				if err != nil {
					return err
				}
			}
		case reflect.Array:
			// fmt.Println("Array")
			printDebug("Array")
			for i := 0; i < v.Len(); i++ {
				item := reflect.New(v.Index(i).Type()).Interface()
				err = unmarshal(r, item, args...)
				v.Index(i).Set(reflect.ValueOf(item).Elem())
				if err != nil {
					return err
				}
			}
		case reflect.Slice:
			printDebug("Slice")
			if v.Len() > 0 {
				for i := 0; i < v.Len(); i++ {
					item := reflect.New(v.Index(i).Type()).Interface()
					err = unmarshal(r, item, args...)
					v.Index(i).Set(reflect.ValueOf(item).Elem())
					if err != nil {
						return err
					}
				}
			}
			// return errors.New("Slices are not suppported")
		default:
			// fmt.Println("Other Type:", v.Type())
			printDebug("Other Type:", v.Type())
			err = binary.Read(r, binary.BigEndian, data)
		}
	}

	if err != nil {
		return err
	}

	return nil

}
