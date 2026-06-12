package attributes

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/Gurvan/melee-data-tools/binread"

	"github.com/Gurvan/melee-data-tools/descriptor"
	. "github.com/Gurvan/melee-data-tools/lib"
	"github.com/Gurvan/melee-data-tools/logger"
)

type UnimplementedSpecialAttributes struct{}

type Special struct {
	any
}

func fighterSwitch(firstRoot string) (any, error) {
	name := strings.TrimPrefix(firstRoot, "ftData")
	if name == firstRoot {
		return nil, errors.New(fmt.Sprintf("File first root %s does not belong to fighter data file.\n", firstRoot))
	}
	switch name {
	case "Captain":
		return any(&Ca{}), nil
	case "Falco":
		return any(&Fc{}), nil
	case "Purin":
		return any(&Pr{}), nil
	case "Mars":
		return any(&Ms{}), nil
	case "Peach":
		return any(&Pe{}), nil
	case "Seak":
		return any(&Sk{}), nil
	case "Fox":
		return any(&Fx{}), nil
	default:
		return any(&UnimplementedSpecialAttributes{}), nil
	}
}

func (a *Special) BinRead(r *binread.Reader, args ...Args) error {
	var firstRoot string
	var err error
	for _, args := range args {
		if desc, ok := args["descriptor"].(*descriptor.Descriptor); ok {
			if firstRoot, err = desc.FirstRootName(); err != nil {
				return err
			}
			break
		} else {
			logger.Warning.Println("SpecialAttributes needs to be parsed as a part of FighterFile")
			return nil
		}
	}

	var attr any
	if attr, err = fighterSwitch(firstRoot); err != nil {
		return err
	}
	if err = r.Decode(attr); err != nil {
		logger.Error.Fatal(err)
		return err
	}
	*a = Special{reflect.ValueOf(attr).Elem().Interface()}
	return nil
}
