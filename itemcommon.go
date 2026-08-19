package mdt

import (
	"github.com/Gurvan/melee-data-tools/item"
	. "github.com/Gurvan/melee-data-tools/lib"
)

// ItemCommonData is the root of the game's common-item file: a block of game-global item
// constants, then the table of items any fighter can come to hold.
type ItemCommonData struct {
	_     Addr // game-global item constants, unread so far
	Items item.CommonItems
}

type ItemCommonFile = File[ItemCommonData]
