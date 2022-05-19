package main

import (
	. "github.com/Gurvan/melee-data-tools/common"
	"github.com/Gurvan/melee-data-tools/fighter"
	"github.com/Gurvan/melee-data-tools/fighter/attributes"
)

type FighterData struct {
	AttributesCommon  Ptr[attributes.Common]
	AttributesSpecial Ptr[attributes.Ca]
	_                 [4]byte
	ActionTable       Ptr[fighter.ActionTable]
	_                 uint32
	_                 [32]byte
	// Hurtboxes uint32
	_   uint32
	_   [16]byte
	ECB Ptr[fighter.ECB]
	// ArticlePointerPtr uint32
	_         uint32
	_         [4]byte
	JostleBox Ptr[fighter.JostleBox]
	_         [12]byte
}

type FighterFile struct {
	Desc Descriptor
	Data FighterData
}
