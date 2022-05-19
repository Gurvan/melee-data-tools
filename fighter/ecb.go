package fighter

import "github.com/Gurvan/melee-data-tools/helpers"

type ECB struct {
	Joints          [6]int16
	Multiplier      float32
	HorizontalScale float32
	VecticalOffset  float32
	VerticalScale   float32
}

func (e ECB) String() string {
	return helpers.PrettyString(e)
}
