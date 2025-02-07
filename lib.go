package mdt

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Gurvan/melee-data-tools/binread"
	"github.com/Gurvan/melee-data-tools/descriptor"
	. "github.com/Gurvan/melee-data-tools/lib"
)

type File[T any] struct {
	Desc descriptor.Descriptor
	Data T
}

func (f *File[T]) BinRead(r *binread.Reader, args ...Args) error {
	err := r.Decode(&f.Desc)
	if err != nil {
		return err
	}

	var argsForData map[string]interface{}
	if len(args) > 0 {
		argsForData = args[0]
		argsForData["descriptor"] = &f.Desc
		argsForData["relocation"] = &f.Desc.Relocation
	} else {
		argsForData = Args{"descriptor": &f.Desc, "relocation": &f.Desc.Relocation}
	}

	return r.Decode(&f.Data, argsForData)
}

func (f *File[T]) AfterParse(r *binread.Reader, _ ...Args) error {
	var firstRoot string
	var err error

	// Check for correct file format
	switch any(f.Data).(type) {
	case FighterData:
		if firstRoot, err = f.Desc.FirstRootName(); err != nil {
			return err
		}
		if !strings.HasPrefix(firstRoot, "ftData") {
			return errors.New(fmt.Sprintf("File first root %s does not belong to fighter data file.\n", firstRoot))
		}
	case AnimationData:
		if firstRoot, err = f.Desc.FirstRootName(); err != nil {
			return err
		}
		if !strings.HasSuffix(firstRoot, "_figatree") {
			return errors.New(fmt.Sprintf("File first root %s does not belong to fighter animation file.\n", firstRoot))
		}
	case ModelData:
		if firstRoot, err = f.Desc.FirstRootName(); err != nil {
			return err
		}
		if !strings.HasSuffix(firstRoot, "_joint") {
			return errors.New(fmt.Sprintf("File first root %s does not belong to fighter model file.\n", firstRoot))
		}
	case CommonData:
		if firstRoot, err = f.Desc.FirstRootName(); err != nil {
			return err
		}
		if !strings.HasPrefix(firstRoot, "ftLoadCommonData") {
			return errors.New(fmt.Sprintf("File first root %s does not belong to player common file.\n", firstRoot))
		}
	case StageData:
		if _, err = f.Desc.FindRootOffset("coll_data"); err != nil {
			return errors.New(fmt.Sprintf("Couldn't parse file as stage file. Error: %s\n", err))
		}
		if _, err = f.Desc.FindRootOffset("map_head"); err != nil {
			return errors.New(fmt.Sprintf("Couldn't parse file as stage file. Error: %s\n", err))
		}
		if _, err = f.Desc.FindRootOffset("grGroundParam"); err != nil {
			return errors.New(fmt.Sprintf("Couldn't parse file as stage file. Error: %s\n", err))
		}
	}

	return nil
}

func (f *File[T]) ReadFromFile(path string) (T, descriptor.Descriptor, error) {
	fileData, err := os.Open(path)
	if err != nil {
		return f.Data, f.Desc, err
	}
	defer fileData.Close()
	reader := binread.NewReader(fileData)
	if err = reader.Decode(f); err != nil {
		return f.Data, f.Desc, err
	}
	return f.Data, f.Desc, nil
}

func (f *File[T]) ReadFromBytes(fileData []byte) (T, descriptor.Descriptor, error) {
	reader := binread.NewReader(bytes.NewReader(fileData))
	if err := reader.Decode(f); err != nil {
		return f.Data, f.Desc, err
	}
	return f.Data, f.Desc, nil
}
