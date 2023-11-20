package attributes

import (
	"github.com/Gurvan/melee-data-tools/helpers"
	. "github.com/Gurvan/melee-data-tools/lib"
)

type ModelParams struct {
	_             int32
	_             Addr
	_             int32
	_             Addr
	ItemHoldBone  uint8
	ShieldBone    uint8
	TopOfHeadBone uint8
	LeftFootBone  uint8
	RightFootBone uint8
}

func (a ModelParams) String() string {
	return helpers.PrettyString(a)
}
