package fighter

import (
	"github.com/Gurvan/melee-data-tools/binread"
	. "github.com/Gurvan/melee-data-tools/lib"
	"github.com/Gurvan/melee-data-tools/logger"
)

const AnimInterpolationSize = 0x02

// AnimInterpolation is one action's entry in the table at fighter data offset 0x10.
type AnimInterpolation struct {
	// Frames is how many frames the game slides the fighter's pose into the action over
	// instead of snapping straight to its animation, and is 0 for most actions. The action
	// a move starts uses it whenever the move does not ask for a length of its own.
	//
	// It is set per character per action: the roster agrees on 6 frames for the idle, the
	// walks and the run, while Peach takes 30 settling into her idle and 10 easing into
	// every fall, and Marth takes 30 dropping into his crouch and 4 into each side special
	// angle.
	Frames uint8
	// Unknown looks to be dynamics related, though what it selects has not been worked out.
	// Only the characters with cloth set it - Peach and Zelda on most of their actions,
	// Marth, Roy and Ganondorf on a handful - and almost only on actions that also set the
	// DisableDynamics action flag (Peach 252 of 255, Zelda 93 of 93), never the other way
	// round. Its values are small and dense per character (Peach 0-85, Zelda 0-35, Roy 0-5,
	// Marth 0-4, Ganondorf 0-1) and group actions that share a pose rather than a length:
	// every idle takes one value, every ground damage another. It indexes no block in the
	// fighter data file or the model file, so it is not a plain table offset.
	Unknown uint8
}

// AnimInterpolationTable runs parallel to the ActionTable - entry i belongs to action i.
type AnimInterpolationTable []AnimInterpolation

var _ binread.BinReader = (*AnimInterpolationTable)(nil)

func (t *AnimInterpolationTable) BinRead(r *binread.Reader, args ...Args) error {
	count, ok := relocEntryCount(r, args, AnimInterpolationSize)
	if !ok {
		logger.Warning.Println("AnimInterpolationTable needs to be parsed as a part of FighterFile")
		return nil
	}

	entries := make([]AnimInterpolation, count)
	if err := r.Decode(&entries); err != nil {
		return err
	}

	*t = entries
	return nil
}

// Frames returns how many frames action i is interpolated into, and 0 for an index the table
// does not cover. The action table is the longer of the two: the parser appends a character's
// subroutines to it, and the file's per-action tables stop before those.
func (t AnimInterpolationTable) Frames(i int) uint8 {
	if i < 0 || i >= len(t) {
		return 0
	}
	return t[i].Frames
}
