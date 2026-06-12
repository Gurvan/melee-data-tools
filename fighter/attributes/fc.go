package attributes

import "github.com/Gurvan/melee-data-tools/helpers"

type Fc struct {
	_                    uint32  // 0x00
	_                    uint32  // 0x04
	_                    uint32  // 0x08
	_                    uint32  // 0x0C
	_                    uint32  // 0x10
	LaserSpeedMultiplier float32 // 0x14
	Unk0x18              uint32
	Unk0x1C              uint32
	Unk0x20              uint32
	_                    uint32
	_                    uint32
	_                    uint32
	_                    uint32
	_                    uint32
	_                    uint32
	_                    uint32
	_                    uint32
	_                    uint32
	_                    uint32
	_                    uint32
	_                    uint32
	_                    uint32
	_                    uint32
	_                    uint32
	_                    uint32
	_                    uint32
	_                    uint32
	_                    uint32
	_                    uint32
	_                    uint32
	_                    uint32
	_                    uint32
	_                    uint32
	_                    uint32
	_                    uint32
}

func (a Fc) String() string {
	return helpers.PrettyString(a)
}
