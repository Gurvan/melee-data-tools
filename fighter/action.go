package fighter

import (
	"github.com/Gurvan/melee-data-tools/binread"
	. "github.com/Gurvan/melee-data-tools/common"
)

type ActionTable []Action

var _ binread.BinReader = (*ActionTable)(nil)

func (t *ActionTable) BinRead(r *binread.Reader, args ...map[string]interface{}) error {
	// var count int = 0
	var count int = 10 // 318
	for _, args := range args {
		var ok bool
		if count, ok = args["count"].(int); ok {
			break
		}
	}
	// for _, arg := range args {
	//         var ok bool
	//         if count, ok = arg.(int); ok {
	//                 break
	//         }
	// }

	actions := make([]Action, count)
	err := r.Decode(&actions)
	if err != nil {
		return err
	}
	*t = actions

	return nil
}

type Action struct {
	Name            Ptr[NullTerminatedString]
	AnimationOffset Addr
	AnimationSize   uint32
	Subaction       Ptr[uint32] // Ptr[[]SubAction]
	Flags           uint32
	_               [4]byte
}
