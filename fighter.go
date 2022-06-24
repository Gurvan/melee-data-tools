package main

import (
	"fmt"
	"io"

	"github.com/Gurvan/melee-data-tools/binread"
	. "github.com/Gurvan/melee-data-tools/common"
	"github.com/Gurvan/melee-data-tools/descriptor"
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
	// _   uint32
	_   [16]byte
	ECB Ptr[fighter.ECB]
	// ArticlePointerPtr uint32
	_         uint32
	_         [4]byte
	JostleBox Ptr[fighter.JostleBox]
	_         [12]byte
}

type FighterFile struct {
	Desc descriptor.Descriptor
	Data FighterData
}

func (f *FighterFile) AfterParse(r *binread.Reader, _ ...Args) error {
	var actionCount int = int(f.Desc.Relocation[f.Data.ActionTable.Offset]) / fighter.ActionSize

	fmt.Println(actionCount)

	before := r.CurrentPosition()

	_, err := r.Seek(f.Data.ActionTable.Offset.ToSeek(), io.SeekStart)
	if err != nil {
		return err
	}

	actionArgs := Args{"actionCount": actionCount}
	err = r.Decode(&f.Data.ActionTable.Value, actionArgs)
	if err != nil {
		return err
	}

	_, err = r.Seek(before, io.SeekStart)
	if err != nil {
		return err
	}
	return nil
}
