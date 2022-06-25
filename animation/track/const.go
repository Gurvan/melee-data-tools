package track

import (
	"errors"
	"fmt"
)

type Type uint8

const (
	NoneType Type = iota
	ROTX
	ROTY
	ROTZ
	PATH
	TRAX
	TRAY
	TRAZ
	SCAX
	SCAY
	SCAZ
	NODE
)

func (t Type) String() string {
	switch t {
	case NoneType:
		return "NoneType"
	case ROTX:
		return "ROTX"
	case ROTY:
		return "ROTY"
	case ROTZ:
		return "ROTZ"
	case PATH:
		return "PATH"
	case TRAX:
		return "TRAX"
	case TRAY:
		return "TRAY"
	case TRAZ:
		return "TRAZ"
	case SCAX:
		return "SCAX"
	case SCAY:
		return "SCAY"
	case SCAZ:
		return "SCAZ"
	case NODE:
		return "NODE"
	default:
		return fmt.Sprintf("TrackTypeUnimplemented(%d)", t)
	}
}

type Interpolation uint8

const (
	NoneInterpolation Interpolation = iota
	CON
	LIN
	SPL0
	SPL
	SLP
	KEY
)

func (i Interpolation) String() string {
	switch i {
	case CON:
		return "CON"
	case LIN:
		return "LIN"
	case SPL0:
		return "SPL0"
	case SPL:
		return "SPL"
	case SLP:
		return "SLP"
	case KEY:
		return "KEY"
	default:
		return fmt.Sprintf("InterpolationUnimplemented(%d)", i)
	}
}

func NewInterpolationError(interp Interpolation) error {
	return errors.New(fmt.Sprintf("Animation key interpolation %d not implemented", interp))
}

type DataFormat uint8

const (
	FLOAT DataFormat = 0x0
	S16   DataFormat = 0x20
	U16   DataFormat = 0x40
	S8    DataFormat = 0x60
	U8    DataFormat = 0x80
)

func (d DataFormat) String() string {
	switch d {
	case FLOAT:
		return "FLOAT"
	case S16:
		return "S16"
	case U16:
		return "U16"
	case S8:
		return "S8"
	case U8:
		return "U8"
	default:
		return fmt.Sprintf("DataFormatUnimplemented(%d)", d)
	}
}

func NewDataFormatError(format DataFormat) error {
	return errors.New(fmt.Sprintf("Animation key data format %d not implemented", format))
}
