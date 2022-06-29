package attributes

import (
	"github.com/Gurvan/melee-data-tools/binread"
	"github.com/Gurvan/melee-data-tools/common"
	"github.com/Gurvan/melee-data-tools/helpers"
)

type ThrowFlags struct {
	DThrowIsWeightIndependant bool `bit:"1"`
	BThrowIsWeightIndependant bool `bit:"1"`
	UThrowIsWeightIndependant bool `bit:"1"`
	FThrowIsWeightIndependant bool `bit:"1"`
}

func (f *ThrowFlags) BinRead(r *binread.Reader, args ...common.Args) error {
	var err error

	byts := make([]byte, 1)
	err = r.Decode(&byts)
	if err != nil {
		return err
	}

	bits := binread.SplitBytes(byts)
	err = binread.BitRead(bits, f, 4)
	return err
}

type Common struct {
	WalkSpeedRelative               float32
	WalkSpeedFixed                  float32
	WalkSpeedMax                    float32
	WalkAnimSpeedSlow               float32
	WalkAnimSpeedMid                float32
	WalkAnimSpeedFast               float32
	Traction                        float32
	DashSpeedStart                  float32
	DashRunSpeedRelative            float32
	DashRunSpeedFixed               float32
	DashRunSpeedMax                 float32
	RunAnimSpeed                    float32
	RunBrakeMaxNumFrames            float32
	LandingSpeedMax                 float32
	JumpSquatNumFrames              float32
	JumpSpeedStartRelative          float32
	JumpFullSpeedVertical           float32
	JumpGroundSpeedConservation     float32
	JumpStartSpeedMax               float32
	JumpShortSpeedVertical          float32
	JumpAirStartSpeedRelative       float32
	JumpAirSpeedVerticalMultiplier  float32
	JumpNum                         uint32
	Gravity                         float32
	FallSlowSpeedVerticalMax        float32
	AirSpeedRelative                float32
	AirSpeedFixed                   float32
	AirSpeedMax                     float32
	AirFriction                     float32
	FallFastSpeedVertical           float32
	_                               float32 // tilt_turn_forced_velocity
	Jab2FrameWindow                 float32
	Jab3FrameWindow                 float32
	TurnSlowDirectionFlipFrameIndex float32
	Weight                          float32
	ModelScale                      float32
	ShieldDefaultSize               float32
	ShieldBrokenStartSpeedVertical  float32
	JabMultiFrameWindow             uint32
	_                               float32 // clank_speed_multiplier
	_                               float32 // hit_by_item_flag
	_                               int32   // unknown
	LedgeJumpStartSpeedHorizontal   float32
	LedgeJumpStartSpeedVertical     float32
	_                               float32 // item_throw_velocity
	_                               float32 // item_throw_damage_scale
	SpecialSGroundSpeedMultiplier   float32
	_                               float32
	_                               float32
	_                               float32
	_                               float32
	_                               float32
	_                               float32
	_                               float32
	_                               float32
	_                               float32
	_                               float32
	LandingLagDefault               float32
	LandingLagNAir                  float32
	LandingLagFAir                  float32
	LandingLagBAir                  float32
	LandingLagUAir                  float32
	LandingLagDAir                  float32
	_                               float32
	TechWallStartSpeedHorizontal    float32
	WallJumpStartSpeedHorizontal    float32
	WallJumpStartSpeedVertical      float32
	_                               float32 // ceiling_tech_x_direction
	_                               float32 // items related
	_                               float32
	_                               float32
	_                               float32
	_                               float32
	_                               float32
	_                               float32
	_                               float32
	_                               float32
	_                               float32
	_                               float32
	_                               float32
	_                               float32
	_                               float32
	_                               float32
	_                               float32
	_                               float32
	_                               float32
	_                               float32
	_                               float32
	_                               float32
	_                               float32
	_                               float32
	_                               float32
	_                               float32
	_                               float32
	_                               float32
	_                               float32
	ThrowFlags                      ThrowFlags
	// ThrowWeightDependant uint8 // F=1,U=2,B=4,D=8 / throw is weight dependant is bit is 0
}

func (a Common) String() string {
	return helpers.PrettyString(a)
}
