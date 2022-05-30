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
	err := r.Decode(&actions)
	if err != nil {
		return err
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
	_               [4]byte
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
