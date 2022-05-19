package attributes

import "github.com/Gurvan/melee-data-tools/helpers"

type Ca struct {
	NSpecialAnglingStickYMax             float32
	NSpecialAnglingStickYMin             float32
	_                                    float32
	_                                    float32
	_                                    float32
	_                                    float32
	_                                    float32
	_                                    float32
	_                                    float32
	_                                    float32
	_                                    float32
	_                                    float32
	_                                    float32
	_                                    float32
	_                                    float32
	_                                    float32
	USpecialAirSpeedRelativeMultiplier   float32
	USpecialAirSpeedMaxMultiplier        float32
	_                                    float32
	USpecialLandingLag                   float32
	_                                    float32
	_                                    float32
	USpecialFlipDirectionStickXThreshold float32
	_                                    float32
	_                                    float32
	_                                    float32
	_                                    float32
	_                                    float32
	_                                    float32
	_                                    float32
	_                                    float32
	_                                    float32
	_                                    float32
	_                                    float32
	_                                    float32
}

func (a Ca) String() string {
	return helpers.PrettyString(a)
}
