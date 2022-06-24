package fighter

import (
	"io"

	"github.com/Gurvan/melee-data-tools/binread"
	. "github.com/Gurvan/melee-data-tools/common"
)

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

	// count = 10

	actions := make([]Action, count)
	// unnamedActionsMap := make(map[Addr]*Action)
	// err := r.Decode(&actions)
	// if err != nil {
	//         return err
	// }

	// beforeActions := r.CurrentPosition()
	// var firstUnnamed int64 = 0
	for i := 0; i < count; i++ {
		before := r.CurrentPosition()
		err := r.Decode(&actions[i])
		if err != nil {
			return err
		}
		// actions[i].ActionOffset = Addr(uint32(before) - uint32(beforeActions))
		// actions[i].ActionOffset = actions[i].Subactions.Offset
		actions[i].ActionOffset = Addr(uint32(before))
		// if actions[i].Name.Value == "" {
		// if firstUnnamed == 0 {
		//         firstUnnamed = r.CurrentPosition()
		// }
		// unnamedActionsMap[actions[i].ActionOffset] = &actions[i]
		// }
		// fmt.Printf("%s | %s | 0x%X\n", actions[i].Name, actions[i].ActionOffset, before)
	}

	// fmt.Println(unnamedActionsMap)
	// fmt.Println("First: ", Addr(firstUnnamed))
	// for addr := range unnamedActionsMap {
	//         fmt.Println(addr)
	// }

	// fmt.Println("#########################")

	newActionsNames := make(map[string]bool)
	for _, action := range actions {
		// fmt.Println(action.Name.Value)
		for _, subaction := range action.Subactions.Value {
			switch s := subaction.(type) {
			case *Subroutine:
				addr := Addr(s.Pointer)
				name := "Subroutine" + addr.String()
				if _, ok := newActionsNames[name]; ok {
					continue
				}
				newAction := Action{}
				newAction.Name = Ptr[NullTerminatedString]{Value: NullTerminatedString(name)}
				newActionsNames[name] = true
				newAction.Subactions.Offset = addr
				newAction.AfterParse(r)
				actions = append(actions, newAction)

				// _, err := r.Seek(int64(s.Pointer+20), io.SeekStart)
				// if err != nil {
				//         return err
				// }
				// var addr Addr
				// err = r.Decode(&addr)
				// if err != nil {
				//         return err
				// }
				// fmt.Println(addr)
				// if a, ok := unnamedActionsMap[addr]; ok {
				//         a.Name.Value = common.NullTerminatedString("Subroutine" + addr.String())
				//         fmt.Println("\t", a.Name.Value)
				// }

			default:
			}
		}
	}

	*t = actions

	return nil
}

const ActionSize = 0x18

type Action struct {
	Name            Ptr[NullTerminatedString]
	AnimationOffset Addr
	AnimationSize   uint32
	Subactions      Ptr[[]SubAction]
	Flags           uint32
	// _               [4]byte
	ActionOffset Addr // Conveniently this field is always 0 in the file so it can be replaced by the offset
}

func (a *Action) AfterParse(r *binread.Reader, _ ...Args) error {
	// fmt.Println(a.Name.Value)
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
		case *EndOfScript, *GoTo, *SubroutineReturn:
			// fmt.Printf("%T\n", subaction)
			// fmt.Println("BREAK")
			break subactionloop
		default:
			// fmt.Printf("%T\n", subaction)
		}
	}

	// subactions = append(subactions, subaction)
	a.Subactions.Value = subactions
	_, err = r.Seek(before, io.SeekStart)
	return err
}
