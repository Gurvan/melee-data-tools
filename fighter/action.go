package fighter

import (
	"fmt"
	"io"
	"regexp"

	"github.com/Gurvan/melee-data-tools/binread"
	. "github.com/Gurvan/melee-data-tools/common"
	"github.com/Gurvan/melee-data-tools/descriptor"
	"github.com/Gurvan/melee-data-tools/logger"
)

type ActionFlags struct {
	UseAnimBasedPhysics        bool        `bit:"1"`
	LoopAnimation              bool        `bit:"1"`
	_                          bool        `bit:"1"`
	_                          bool        `bit:"1"`
	DisableDynamics            bool        `bit:"1"`
	_                          bool        `bit:"1"`
	TransNAffectedByModelScale bool        `bit:"1"`
	_                          uint32      `bit:"3"`
	_                          uint32      `bit:"13"`
	DisableBlendOnBoneIndex    uint32      `bit:"3"`
	CharacterID                CharacterID `bit:"6"`
}

func (f *ActionFlags) BinRead(r *binread.Reader, args ...Args) error {
	var err error

	byts := make([]byte, 4)
	err = r.Decode(&byts)
	if err != nil {
		return err
	}

	bits := binread.SplitBytes(byts)
	err = binread.BitRead(bits, f, 0)
	return err
}

const ActionSize = 0x18

type Action struct {
	Name          Ptr[NullTerminatedString]
	Animation     Ptr[any]
	AnimationSize uint32
	Subactions    Ptr[[]SubAction]
	Flags         ActionFlags
	_             [4]byte
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
		case *EndOfScript, *GoTo, *SubroutineReturn:
			break subactionloop
		default:
		}
	}

	// a.Subactions.Value = subactions
	a.Subactions.SetValue(subactions)
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

	subroutineIndex := 0
	actionsNames := make(map[string]int)
	subroutineByAddr := make(map[Addr]string)
	for actionIndex, action := range actions {
		name := action.Name.GetValue().String()

		// Name unnamed actions
		if name == "" {
			name = "Function_" + fmt.Sprint(actionIndex)
			// actions[actionIndex].Name.Value = NullTerminatedString(name)
			actions[actionIndex].Name.SetValue(NullTerminatedString(name))
		}

		name = truncActionName(name)

		// Increment names for actions with duplicate names
		if duplicate, ok := actionsNames[name]; ok {
			actionsNames[name] = duplicate + 1
			// suffix := "_figatree"
			// prefix := strings.TrimSuffix(name, suffix)
			// name = prefix + "_" + fmt.Sprintf("%d", duplicate) + suffix
			name = name + "_" + fmt.Sprintf("%d", duplicate)
			actions[actionIndex].Name.SetValue(NullTerminatedString(name))
		} else {
			actionsNames[name] = 1
		}

		handleSubaction := func(pointer uint32) string {
			addr := Addr(pointer + 0x20)
			// addr := Addr(s.Target)
			name := ""
			if n, ok := subroutineByAddr[addr]; ok {
				// log.Printf("[Subroutine_%d] Subroutine Fetched. Addr: %x\n", subroutineIndex, pointer+0x20)
				name = n
			} else {
				// log.Printf("[Subroutine_%d] Subroutine Add. Addr: %x\n", subroutineIndex, pointer+0x20)
				name = "Subroutine_" + fmt.Sprint(subroutineIndex)
				subroutineByAddr[addr] = name
				subroutineIndex++
			}
			// name := "Subroutine" + addr.String()
			if _, ok := actionsNames[name]; ok {
				return name
			}
			newAction := Action{}
			newAction.Name = Ptr[NullTerminatedString]{}
			newAction.Name.SetValue(NullTerminatedString(name))
			actionsNames[name] = 1
			newAction.Subactions.Offset = addr
			newAction.AfterParse(r)
			actions = append(actions, newAction)
			return name
		}

		// Add subroutines to actions slice
		for _, subaction := range action.Subactions.GetValue() {
			switch s := subaction.(type) {
			case *Subroutine:
				s.PointerName = handleSubaction(s.Pointer)
			case *GoTo:
				s.PointerName = handleSubaction(s.Pointer)
			default:
			}
		}
	}

	*t = actions
	return nil
}

func truncActionName(name string) string {
	var re *regexp.Regexp
	var substrings []string
	// substrings := re.FindStringSubmatch(name)
	// if len(substrings) > 0 {
	//         return substrings[len(substrings)-1]
	// }

	re, _ = regexp.Compile(".*_Share_ACTION_(.*)+_figatree")
	substrings = re.FindStringSubmatch(name)
	if len(substrings) == 0 {
		// return ""
		return name
	}
	return substrings[len(substrings)-1]
}
