package fighter

import (
	"fmt"
	"io"
	"strings"

	"github.com/Gurvan/melee-data-tools/binread"
	. "github.com/Gurvan/melee-data-tools/common"
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
	for _, args := range args {
		var ok bool
		if count, ok = args["actionCount"].(int); ok {
			break
		} else {
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
