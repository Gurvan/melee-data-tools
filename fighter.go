package mdt

import (
	"os"

	"github.com/Gurvan/melee-data-tools/binread"
	. "github.com/Gurvan/melee-data-tools/common"
	"github.com/Gurvan/melee-data-tools/fighter"
	"github.com/Gurvan/melee-data-tools/fighter/attributes"
	"github.com/Gurvan/melee-data-tools/logger"
	"github.com/Gurvan/melee-data-tools/model"
)

type FighterData[AS any] struct {
	AttributesCommon  Ptr[attributes.Common]
	AttributesSpecial Ptr[attributes.Special]
	ModelParams       Ptr[attributes.ModelParams]
	ActionTable       Ptr[fighter.ActionTable]
	_                 [32]byte
	Hurtboxes         Ptr[fighter.Hurtboxes]
	_                 [16]byte
	ECB               Ptr[fighter.ECB]
	_                 uint32 // ArticlePointerPtr uint32
	_                 [4]byte
	JostleBox         Ptr[fighter.JostleBox]
	_                 [8]byte
	Model             Ptr[model.Joint]
}

type FighterFile = File[FighterData[any]]

func (f *FighterData[any]) AfterParse(r *binread.Reader, _ ...Args) error {
	if f.Model.ValuePtr == nil {
		return nil
	}

	model := f.Model.ValuePtr
	name := NullTerminatedString("Metal")
	model.Name.ValuePtr = &name

	return nil
}

func (f *FighterData[AS]) ParseAnimation(animationPath string) error {
	var err error
	var file *os.File
	var animationFilesAsBytes [][]byte
	var offsets []int64

	if file, err = os.Open(animationPath); err != nil {
		return err
	}
	defer file.Close()
	r := binread.NewReader(file)

	if animationFilesAsBytes, offsets, err = SplitAnimationFile(r); err != nil {
		return err
	}

	var animationFile AnimationFile
	var animationData AnimationData

	animations := make(map[int64]AnimationData)
	for i, fileData := range animationFilesAsBytes {
		if animationData, _, err = animationFile.ReadFromBytes(fileData); err != nil {
			return err
		}
		animations[offsets[i]] = animationData
	}

	for i, action := range f.ActionTable.GetValue() {
		if action.AnimationSize == 0 {
			continue
		}
		offset := action.Animation.Offset.ToSeek() - 0x20
		if animationData, ok := animations[offset]; !ok {
			logger.Warning.Printf("Animation 0x%X not found.", offset)
		} else {
			// f.ActionTable.Value[i].Animation.Value = animationData
			f.ActionTable.GetValue()[i].Animation.SetValue(animationData)
		}
	}
	return nil
}

func (f *FighterData[AS]) ParseModel(modelPath string) error {
	var modelFile ModelFile

	modelData, desc, err := modelFile.ReadFromFile(modelPath)
	if err != nil {
		return err
	}
	if firstRoot, err := desc.FirstRootName(); err == nil {
		modelData.UpdateName(firstRoot)
	}

	f.Model.SetValue(modelData.Joint)
	return nil
}
