package hurtbox

import "fmt"

type Position uint32

func (h Position) String() string {
	switch h {
	case Low:
		return "Low"
	case Mid:
		return "Mid"
	case High:
		return "High"
	default:
		return fmt.Sprintf("HurtboxPositionUnimplemented(%d)", h)
	}
}

const (
	Low Position = iota
	Mid
	High
)

type Grabbable uint32

func (g Grabbable) String() string {
	switch g {
	case No:
		return "No"
	case Yes:
		return "Yes"
	default:
		return fmt.Sprintf("HurtboxGrabbableUnimplemented(%d)", g)
	}
}

const (
	No Grabbable = iota
	Yes
)
