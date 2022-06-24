package attributes

import "github.com/Gurvan/melee-data-tools/helpers"

type Ca struct {
	NSpecialAnglingStickYMax             float32
	NSpecialAnglingStickYMin             float32
	NSpecialMaxAngleDeg                  float32
	NSpecialAnglingMagnitude             float32
	NSpecialAirSpeedMult                 float32
	_                                    float32
	SSpecialAirGravity                   float32
	SSpecialAirFallSpeedMax              float32
	_                                    float32
	_                                    float32
	_                                    float32
	_                                    float32
	_                                    float32
	_                                    float32
	SSpecialLandingLag                   float32
	SSpecialHitLandingLag                float32
	USpecialAirSpeedRelativeMultiplier   float32
	USpecialAirSpeedMaxMultiplier        float32
	_                                    float32
	USpecialLandingLag                   float32
	_                                    float32
	_                                    float32
	USpecialFlipDirectionStickXThreshold float32
	_                                    float32
	USpecialGravity                      float32
	_                                    float32
	_                                    float32
	_                                    float32
	_                                    float32
	DSpecialAfterHitSpeedMult            float32
	DSpecialMaxNumHitForSpeedMult        uint32
	DSpecialGroundEndAnimSpeed           float32
	DSpecielAirEndGroundAnimSpeed        float32
	DSpecialGroundEndTractionMult        float32
	DSpecialAirEndFrictionMult           float32
}

func (a Ca) String() string {
	return helpers.PrettyString(a)
}
