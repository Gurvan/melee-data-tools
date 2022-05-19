package fighter

import (
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
		}
	}

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
	Subaction       Ptr[uint32] // Ptr[[]SubAction]
	Flags           uint32
	_               [4]byte
}
