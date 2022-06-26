package common

import (
	"fmt"

	"github.com/Gurvan/melee-data-tools/binread"
)

type CharacterID uint32

const (
	Mario         CharacterID = 0x00
	Fox           CharacterID = 0x01
	CaptainFalcon CharacterID = 0x02
	DonkeyKong    CharacterID = 0x03
	Kirby         CharacterID = 0x04
	Bowser        CharacterID = 0x05
	Link          CharacterID = 0x06
	Shiek         CharacterID = 0x07
	Ness          CharacterID = 0x08
	Peach         CharacterID = 0x09
	Popo          CharacterID = 0x0A
	Nana          CharacterID = 0x0B
	Pikachu       CharacterID = 0x0C
	Samus         CharacterID = 0x0D
	Yoshi         CharacterID = 0x0E
	Jigglypuff    CharacterID = 0x0F
	Mewtwo        CharacterID = 0x10
	Luigi         CharacterID = 0x11
	Marth         CharacterID = 0x12
	Zelda         CharacterID = 0x13
	YoungLink     CharacterID = 0x14
	DrMario       CharacterID = 0x15
	Falco         CharacterID = 0x16
	Pichu         CharacterID = 0x17
	Watch         CharacterID = 0x18
	Ganondorf     CharacterID = 0x19
	Roy           CharacterID = 0x1A
	GameAndWatch  CharacterID = 0x1B
	WireframeM    CharacterID = 0x1D
	WireframeF    CharacterID = 0x1E
	GigaBowser    CharacterID = 0x1F
	Sandbag       CharacterID = 0x20
)

func (c *CharacterID) BinRead(r *binread.Reader, args ...Args) error {
	return nil

}

func (t CharacterID) String() string {
	switch t {
	case Mario:
		return "Mario"
	case Fox:
		return "Fox"
	case CaptainFalcon:
		return "CaptainFalcon"
	case DonkeyKong:
		return "DonkeyKong"
	case Kirby:
		return "Kirby"
	case Bowser:
		return "Bowser"
	case Link:
		return "Link"
	case Shiek:
		return "Shiek"
	case Ness:
		return "Ness"
	case Peach:
		return "Peach"
	case Popo:
		return "Popo"
	case Nana:
		return "Nana"
	case Pikachu:
		return "Pikachu"
	case Samus:
		return "Samus"
	case Yoshi:
		return "Yoshi"
	case Jigglypuff:
		return "Jigglypuff"
	case Mewtwo:
		return "Mewtwo"
	case Luigi:
		return "Luigi"
	case Marth:
		return "Marth"
	case Zelda:
		return "Zelda"
	case YoungLink:
		return "YoungLink"
	case DrMario:
		return "DrMario"
	case Falco:
		return "Falco"
	case Pichu:
		return "Pichu"
	case Watch:
		return "Watch"
	case Ganondorf:
		return "Ganondorf"
	case Roy:
		return "Roy"
	case GameAndWatch:
		return "GameAndWatch"
	case WireframeM:
		return "WireframeM"
	case WireframeF:
		return "WireframeF"
	case GigaBowser:
		return "GigaBowser"
	case Sandbag:
		return "Sandbag"
	default:
		return fmt.Sprintf("CharacterIDUnimplemented(%d)", t)
	}
}
