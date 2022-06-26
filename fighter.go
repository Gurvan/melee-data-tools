package mdt

import (
	. "github.com/Gurvan/melee-data-tools/common"
	"github.com/Gurvan/melee-data-tools/fighter"
	"github.com/Gurvan/melee-data-tools/fighter/attributes"
)

type FighterData[AS any] struct {
	AttributesCommon  Ptr[attributes.Common]
	AttributesSpecial Ptr[attributes.SpecialAttributes]
	_                 [4]byte
	ActionTable       Ptr[fighter.ActionTable]
	_                 [32]byte
	Hurtboxes         Ptr[fighter.Hurtboxes]
	_                 [16]byte
	ECB               Ptr[fighter.ECB]
	_                 uint32 // ArticlePointerPtr uint32
	_                 [4]byte
	JostleBox         Ptr[fighter.JostleBox]
	_                 [12]byte
}

type FighterFile = File[FighterData[any]]
