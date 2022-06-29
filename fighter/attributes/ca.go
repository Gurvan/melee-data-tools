package attributes

import "github.com/Gurvan/melee-data-tools/helpers"

type Ca struct {
	SpecialNAnglingStickYMax             float32
	SpecialNAnglingStickYMin             float32
	SpecialNMaxAngleDeg                  float32
	SpecialNAnglingMagnitude             float32
	SpecialNAirSpeedMult                 float32
	_                                    float32
	SpecialSAirGravity                   float32
	SpecialSAirFallSpeedMax              float32
	_                                    float32
	_                                    float32
	_                                    float32
	_                                    float32
	_                                    float32
	_                                    float32
	SpecialSLandingLag                   float32
	SpecialSHitLandingLag                float32
	SpecialUAirSpeedRelativeMultiplier   float32
	SpecialUAirSpeedMaxMultiplier        float32
	_                                    float32
	SpecialULandingLag                   float32
	_                                    float32
	_                                    float32
	SpecialUFlipDirectionStickXThreshold float32
	_                                    float32
	SpecialUGravity                      float32
	_                                    float32
	_                                    float32
	_                                    float32
	_                                    float32
	SpecialDAfterHitSpeedMult            float32
	SpecialDMaxNumHitForSpeedMult        uint32
	SpecialDGroundEndAnimSpeed           float32
	DSpecielAirEndGroundAnimSpeed        float32
	SpecialDGroundEndTractionMult        float32
	SpecialDAirEndFrictionMult           float32
}

func (a Ca) String() string {
	return helpers.PrettyString(a)
}
