package attributes

import "github.com/Gurvan/melee-data-tools/helpers"

type Common struct {
	WalkSpeedRelative               float32
	WalkSpeedFixed                  float32
	WalkSpeedMax                    float32
	WalkAnimSpeedSlow               float32
	WalkAnimSpeeMid                 float32
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
	ThrowWeightDependant            uint8 // F=1,U=2,B=4,D=8 / throw is weight dependant is bit is 0
}

func (a Common) String() string {
	return helpers.PrettyString(a)
}

// pub tilt_turn_forced_velocity: f32,
// pub clank_speed_multiplier: f32,
// pub hit_by_item_flag: f32,
// pub unknown0x_a4: i32,
// pub item_throw_velocity: f32,
// pub item_throw_damage_scale: f32,
// pub run_side_special_momentum: f32,
// pub egg_size: f32,
// pub egg_hurtbox: f32,
// pub egg_hurtbox_x: f32,
// pub egg_hurtbox_y: f32,
// pub egg_hurtbox_z: f32,
// pub unknown0x_d0: f32,
// pub unknown0x_d4: f32,
// pub egg_hurtbox_radius: f32,
// pub capture_neutral_special_absorb_damage: f32,
// pub capture_neutral_special_damage: f32,
// pub victory_screen_model_scale: f32,
// pub ceiling_tech_x_direction: f32,
// pub unknown0x110: f32,
// pub left_bunny_hood_x: f32,
// pub left_bunny_hood_y: f32,
// pub left_bunny_hood_z: f32,
// pub right_bunny_hood_x: f32,
// pub right_bunny_hood_y: f32,
// pub right_bunny_hood_z: f32,
// pub bunny_hood_size: f32,
// pub flower_x: f32,
// pub flower_y: f32,
// pub flower_z: f32,
// pub flower_size: f32,
// pub screw_attack_upward_knockback: f32,
// pub screw_attack_effect_size: f32,
// pub unknown0x148: f32,
// pub bubble_ratio: f32,
// pub freeze_offset1: f32,
// pub freeze_offset2: f32,
// pub freeze_escape_height: f32,
// pub freeze_escape_x_momentum: f32,
// pub frozen_size: f32,
// pub warp_star_hitbox_scaling: f32,
// pub unknown0x168: f32,
// pub camera_zoom_target_bone: i32,
// pub magnified_x_sway: f32,
// pub magnified_y_sway: f32,
// pub magnified_z_sway: f32,
// pub footstool_y_offset: f32,
