package fighter

import (
	"fmt"
	"io"
	"strings"

	"github.com/Gurvan/melee-data-tools/binread"
	. "github.com/Gurvan/melee-data-tools/common"
	"github.com/Gurvan/melee-data-tools/descriptor"
	"github.com/Gurvan/melee-data-tools/logger"
)

const ActionSize = 0x18

type Action struct {
	Name            Ptr[NullTerminatedString]
	AnimationOffset Addr
	AnimationSize   uint32
	Subactions      Ptr[[]SubAction]
	Flags           uint32
	_               [4]byte
}

func (a *Action) AfterParse(r *binread.Reader, _ ...Args) error {
	subactions := make([]SubAction, 0)

	before := r.CurrentPosition()

	_, err := r.Seek(a.Subactions.Offset.ToSeek(), io.SeekStart)
	if err != nil {
		return err
	}

subactionloop:
	for {
		subaction, err := GetSubActionType(r)
		if err != nil {
			if _, ok := err.(*SubActionNotImplemented); ok {
				var x [4]byte
				r.Decode(&x)
				// log.Println(err)
				continue
			} else {
				return err
			}
		}

		err = DecodeSubAction(r, subaction)
		if err != nil {
			return err
		}

		subactions = append(subactions, subaction)
		switch subaction.(type) {
		// case *EndOfScript, *GoTo, *SubroutineReturn:
		case *EndOfScript, *SubroutineReturn:
			// fmt.Printf("%T\n", subaction)
			// fmt.Println("BREAK")
			break subactionloop
		default:
			// fmt.Printf("%T\n", subaction)
		}
	}

	a.Subactions.Value = subactions
	_, err = r.Seek(before, io.SeekStart)
	return err
}

type ActionTable []Action

var _ binread.BinReader = (*ActionTable)(nil)

func (t *ActionTable) BinRead(r *binread.Reader, args ...Args) error {
	var count int = 0

	// Get number of actions from the relocation table in the file descriptor.
	// Can't get it if the ActionTable is not parsed as a part of a FighterFile,
	// so do nothing in this case.
	if len(args) == 0 {
		logger.Warning.Println("ActionTable needs to be parsed as a part of FighterFile")
		return nil
	}
	for _, args := range args {
		if desc, ok := args["descriptor"].(*descriptor.Descriptor); ok {
			offset := Addr(r.CurrentPosition())
			count = int(desc.Relocation[offset]) / ActionSize
			break
		} else {
			logger.Warning.Println("ActionTable needs to be parsed as a part of FighterFile")
			return nil
		}
	}

	actions := make([]Action, count)
	err := r.Decode(&actions)
	if err != nil {
		return err
	}

	actionsNames := make(map[string]int)
	for actionIndex, action := range actions {
		name := action.Name.Value.String()

		// Name unnamed actions
		if name == "" {
			name = "Function_" + fmt.Sprint(actionIndex)
			actions[actionIndex].Name.Value = NullTerminatedString(name)
		}

		// Increment names for actions with duplicate names
		if duplicate, ok := actionsNames[name]; ok {
			actionsNames[name] = duplicate + 1
			suffix := "_figatree"
			prefix := strings.TrimSuffix(name, suffix)
			name = prefix + "_" + fmt.Sprintf("%d", duplicate) + suffix
			actions[actionIndex].Name.Value = NullTerminatedString(name)
		} else {
			actionsNames[name] = 1
		}

		// Add subroutines to actions slice
		for _, subaction := range action.Subactions.Value {
			switch s := subaction.(type) {
			case *Subroutine:
				addr := Addr(s.Pointer)
				name := "Subroutine" + addr.String()
				if _, ok := actionsNames[name]; ok {
					continue
				}
				newAction := Action{}
				newAction.Name = Ptr[NullTerminatedString]{Value: NullTerminatedString(name)}
				actionsNames[name] = 1
				newAction.Subactions.Offset = addr
				newAction.AfterParse(r)
				actions = append(actions, newAction)
			default:
			}
		}
	}

	*t = actions
	return nil
}
