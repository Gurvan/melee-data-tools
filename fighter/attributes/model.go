package attributes

import (
	"github.com/Gurvan/melee-data-tools/helpers"
)

type ModelParams struct {
	_             [4]byte
	_             [4]byte
	_             [4]byte
	_             [4]byte
	ItemHoldBone  byte
	ShieldBone    byte
	TopOfHeadBone byte
	LeftFootBone  byte
	RightFootBone byte
}

func (a ModelParams) String() string {
	return helpers.PrettyString(a)
}
