package model

import (
	"fmt"

	. "github.com/Gurvan/melee-data-tools/lib"
)

type Flags uint32

func (f Flags) String() string {
	return fmt.Sprintf("0x%X", uint32(f))
}

type Vec3[T any] struct {
	X T
	Y T
	Z T
}

type Joint struct {
	Index       int32 `binread:"ignore"`
	Name        OptionalPtr[NullTerminatedString]
	Flags       Flags
	Child       OptionalPtr[Joint]
	Next        OptionalPtr[Joint]
	_           [4]byte
	Rotation    Vec3[float32]
	Scale       Vec3[float32]
	Translation Vec3[float32]
	_           [8]byte
}
