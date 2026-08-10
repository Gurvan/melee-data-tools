package mdt

import (
	"bytes"
	"os"

	"github.com/Gurvan/melee-data-tools/binread"
	"github.com/Gurvan/melee-data-tools/fighter"
	"github.com/Gurvan/melee-data-tools/fighter/attributes"
	"github.com/Gurvan/melee-data-tools/item"
	. "github.com/Gurvan/melee-data-tools/lib"
	"github.com/Gurvan/melee-data-tools/logger"
	"github.com/Gurvan/melee-data-tools/model"
)

type FighterData struct {
	AttributesCommon  Ptr[attributes.Common]
	AttributesSpecial Ptr[attributes.Special]
	ModelParams       Ptr[attributes.ModelParams]
	ActionTable       Ptr[fighter.ActionTable]
	AnimInterpolation Ptr[fighter.AnimInterpolationTable]
	_                 [8]byte
	PartPoses         Ptr[fighter.PartPoseTable]
	ShieldPose        Ptr[ShieldPose]
	_                 [8]byte
	Dynamics          OptionalPtr[fighter.Dynamics]
	Hurtboxes         Ptr[fighter.Hurtboxes]
	_                 [12]byte
	ItemPickupBoxes   Ptr[fighter.ItemPickupBoxes]
	ECB               Ptr[fighter.ECB]
	Items             item.Items
	_                 [4]byte
	JostleBox         Ptr[fighter.JostleBox]
	_                 [8]byte
	Model             Ptr[model.Joint]
}

type ShieldPose struct {
	Root    OptionalPtr[model.Joint]
	Unknown Addr
}

func (s *ShieldPose) BinRead(r *binread.Reader, args ...Args) error {
	if err := r.Decode(&s.Root, args...); err != nil {
		return err
	}
	return r.Decode(&s.Unknown)
}

func (s *ShieldPose) Joint() *model.Joint {
	if s == nil || s.Root.ValuePtr == nil {
		return nil
	}
	return s.Root.ValuePtr.Child.ValuePtr
}

type FighterFile = File[FighterData]

func (f *FighterData) AfterParse(r *binread.Reader, _ ...Args) error {
	if f.Model.ValuePtr == nil {
		return nil
	}

	model := f.Model.ValuePtr
	name := NullTerminatedString("Metal")
	model.Name.ValuePtr = &name

	return nil
}

func (f *FighterData) ParseAnimationFromFile(animationPath string) error {
	file, err := os.Open(animationPath)
	if err != nil {
		return err
	}
	defer file.Close()
	r := binread.NewReader(file)
	return f.parseAnimation(r)
}

func (f *FighterData) ParseAnimationFromBytes(data []byte) error {
	r := binread.NewReader(bytes.NewReader(data))
	return f.parseAnimation(r)
}

func (f *FighterData) parseAnimation(r binread.Reader) error {
	animationFilesAsBytes, offsets, err := SplitAnimationFile(r)
	if err != nil {
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
			f.ActionTable.GetValue()[i].Animation.SetValue(animationData)
		}
	}
	return nil
}

func (f *FighterData) ParseModelFromFile(modelPath string) error {
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

func (f *FighterData) ParseModelFromBytes(data []byte) error {
	var modelFile ModelFile

	modelData, desc, err := modelFile.ReadFromBytes(data)
	if err != nil {
		return err
	}
	if firstRoot, err := desc.FirstRootName(); err == nil {
		modelData.UpdateName(firstRoot)
	}

	f.Model.SetValue(modelData.Joint)
	return nil
}
