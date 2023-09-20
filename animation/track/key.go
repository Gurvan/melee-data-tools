package track

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/Gurvan/melee-data-tools/binread"
	. "github.com/Gurvan/melee-data-tools/lib"
)

type Key struct {
	Frame         float32
	Value         float32
	Tangent       float32
	Interpolation Interpolation
}

type Keys []Key

var _ binread.BinReader = (*Keys)(nil)

func (k *Keys) BinRead(r *binread.Reader, args ...Args) error {
	errorReadKeys := func(s string) error {
		return errors.New(fmt.Sprintf("When reading animation Keys, Args should contain %s", s))
	}
	var valueFormat DataFormat
	var valueScale uint32
	var tangentFormat DataFormat
	var tangentScale uint32
	var dataLength uint16
	if len(args) == 0 {
		return errorReadKeys("\"valueScale\", \"valueFormat\", \"tangentScale\", \"tangentFormat\", \"dataLength\"")
	}
	for _, args := range args {
		var ok bool
		if valueFormat, ok = args["valueFormat"].(DataFormat); !ok {
			return errorReadKeys("\"valueFormat\"")
		}
		if valueScale, ok = args["valueScale"].(uint32); !ok {
			return errorReadKeys("\"valueScale\"")
		}
		if tangentFormat, ok = args["tangentFormat"].(DataFormat); !ok {
			return errorReadKeys("\"tangentFormat\"")
		}
		if tangentScale, ok = args["tangentScale"].(uint32); !ok {
			return errorReadKeys("\"tangentScale\"")
		}
		if dataLength, ok = args["dataLength"].(uint16); !ok {
			return errorReadKeys("\"dataLength\"")
		}
	}
	readPacked := func(r *binread.Reader) (uint16, error) {
		var valueHalf uint8
		err := r.Decode(&valueHalf)
		if err != nil {
			return 0, err
		}
		value := uint16(valueHalf)
		if valueHalf>>7 == 0 {
			return value, nil
		}
		err = r.Decode(&valueHalf)
		if err != nil {
			return 0, err
		}
		return value&0x7F | (uint16(valueHalf) << 7), nil
	}
	getNumKeysAndInterp := func(value uint16) (uint16, Interpolation) {
		n := value>>4 + 1
		interp := Interpolation(value & 0x0F)
		return n, interp
	}
	readFloat := func(r *binread.Reader, format DataFormat, scale uint32) (float32, error) {
		var err error
		switch format {
		case FLOAT:
			var value float32
			err = binary.Read(r, binary.LittleEndian, &value)
			return value, err
		case S16:
			var value int16
			err = binary.Read(r, binary.LittleEndian, &value)
			return float32(value) / float32(scale), err
		case U16:
			var value uint16
			err = binary.Read(r, binary.LittleEndian, &value)
			return float32(value) / float32(scale), err
		case S8:
			var value int8
			err = binary.Read(r, binary.LittleEndian, &value)
			return float32(value) / float32(scale), err
		case U8:
			var value uint8
			err = binary.Read(r, binary.LittleEndian, &value)
			return float32(value) / float32(scale), err
		default:
			return 0, NewDataFormatError(format)
		}
	}

	// Follow pointer
	var addr Addr
	err := r.Decode(&addr)
	if err != nil {
		return err
	}

	before := r.CurrentPosition()
	_, err = r.Seek(addr.ToSeek(), io.SeekStart)
	if err != nil {
		return err
	}

	keys := make([]Key, 0)
	var frame float32 = 0
	readPosAtStart := r.CurrentPosition()
	for {
		if r.CurrentPosition()-readPosAtStart >= int64(dataLength) {
			break
		}
		value, err := readPacked(r)
		if err != nil {
			return err
		}
		n, interp := getNumKeysAndInterp(value)
		if interp == NoneInterpolation {
			break
		}
		for i := 0; i < int(n); i++ {
			key := Key{}
			key.Interpolation = interp
			key.Frame = frame

			// Values
			switch interp {
			case NoneInterpolation, SLP:
			case CON, LIN, SPL0, SPL, KEY:
				key.Value, err = readFloat(r, valueFormat, valueScale)
			default:
				return NewInterpolationError(interp)
			}
			if err != nil {
				return err
			}

			// Tangents
			switch interp {
			case NoneInterpolation, CON, LIN, SPL0, KEY:
			case SPL, SLP:
				key.Tangent, err = readFloat(r, tangentFormat, tangentScale)
			default:
				return NewInterpolationError(interp)
			}
			if err != nil {
				return err
			}

			// Frame increment
			var time uint16
			switch interp {
			case NoneInterpolation, SLP, KEY:
			case CON, LIN, SPL0, SPL:
				time, err = readPacked(r)
			default:
				return NewInterpolationError(interp)
			}
			if err != nil {
				return err
			}

			frame += float32(time)
			keys = append(keys, key)
		}
	}

	*k = keys
	_, err = r.Seek(before, io.SeekStart)
	return err
}
