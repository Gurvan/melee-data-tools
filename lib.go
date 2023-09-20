package mdt

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Gurvan/melee-data-tools/binread"
	. "github.com/Gurvan/melee-data-tools/lib"
	"github.com/Gurvan/melee-data-tools/descriptor"
)

type File[T any] struct {
	Desc descriptor.Descriptor
	Data T
}

func (f *File[T]) BinRead(r *binread.Reader, args ...Args) error {
	var err error
	err = r.Decode(&f.Desc)
	if err != nil {
		return err
	}

	var argsForData map[string]interface{}
	if len(args) > 0 {
		argsForData = args[0]
		argsForData["descriptor"] = &f.Desc
	} else {
		argsForData = Args{"descriptor": &f.Desc}
	}

	err = r.Decode(&f.Data, argsForData)
	return nil
}

func (f *File[T]) AfterParse(r *binread.Reader, _ ...Args) error {
	var firstRoot string
	var err error

	// Check for correct file format
	switch any(f.Data).(type) {
	case FighterData[any]:
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
	}

	return nil
}

func (f *File[T]) ReadFromFile(path string) (T, descriptor.Descriptor, error) {
	fileData, err := os.Open(path)
	defer fileData.Close()
	if err != nil {
		return f.Data, f.Desc, err
	}
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
