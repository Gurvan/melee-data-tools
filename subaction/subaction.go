package subaction

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/Gurvan/melee-data-tools/binread"
	. "github.com/Gurvan/melee-data-tools/lib"
)

type emptySubAction struct{}

func (emptySubAction) isSubaction() {}

type SubActionNotImplemented struct {
	Id     uint8
	IsItem bool
}

func (e *SubActionNotImplemented) Error() string {
	if e.IsItem {
		return fmt.Sprintf("Item SubAction 0x%X is not implemented", e.Id)
	}
	return fmt.Sprintf("Fighter SubAction 0x%X is not implemented", e.Id)
}

type SubAction interface {
	isSubaction()
}

var _ SubAction = (*emptySubAction)(nil)

func SubActionString(s SubAction) string {
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	str := t.Name() + "{"
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous {
			continue
		}
		if field.Name == "_" {
			continue
		}
		str += field.Name + ": " + fmt.Sprintf("%v", v.Field(i).Interface())
		if i < t.NumField()-1 {
			str += ", "
		}
	}

	// str = str[:len(str)-2] + "}"
	str += "}"
	str = strings.Replace(str, ", }", "}", 1)
	return str
}

// 0x00
type EndOfScript struct {
	// emptySubAction
	_ uint32 `bit:"26"`
}

func (EndOfScript) isSubaction() {}

// 0x01
type SynchronousTimer struct {
	// emptySubAction
	Frame uint32 `bit:"26"`
}

func (SynchronousTimer) isSubaction() {}

// 0x02
type AsynchronousTimer struct {
	// emptySubAction
	Frame uint32 `bit:"26"`
}

func (AsynchronousTimer) isSubaction() {}

// 0x03
type SetLoop struct {
	Count uint32 `bit:"26"`
}

func (SetLoop) isSubaction() {}

// 0x04
type ExecuteLoop struct {
	_ uint32 `bit:"26"`
}

func (ExecuteLoop) isSubaction() {}

// 0x05
type Subroutine struct {
	Target      uint32 `bit:"26"`
	Pointer     uint32 `bit:"32"`
	PointerName string `bit:"ignore"`
}

func (Subroutine) isSubaction() {}

// 0x06
type SubroutineReturn struct {
	_ uint32 `bit:"26"`
}

func (SubroutineReturn) isSubaction() {}

// 0x07
type GoTo struct {
	Target      uint32 `bit:"26"`
	Pointer     uint32 `bit:"32"`
	PointerName string `bit:"ignore"`
}

func (GoTo) isSubaction() {}

// 0x08
type SetTimerAnimation struct {
	_ uint32 `bit:"26"`
}

func (SetTimerAnimation) isSubaction() {}

// 0x09
type Unknown0x09 struct {
	_ uint32 `bit:"26"`
}

func (Unknown0x09) isSubaction() {}

// 0x0A
type GraphicEffect struct {
	BoneId               uint32 `bit:"8"`
	UseCommonBoneIDs     bool   `bit:"1"`
	DestroyOnStateChange bool   `bit:"1"`
	_                    uint32 `bit:"16"`
	GFXID                uint32 `bit:"16"`
	_                    uint32 `bit:"16"`
	Z                    int32  `bit:"16"`
	Y                    int32  `bit:"16"`
	X                    int32  `bit:"16"`
	RangeZ               int32  `bit:"16"`
	RangeY               int32  `bit:"16"`
	RangeX               int32  `bit:"16"`
	// _                    uint32 `bit:"8"`
}

func (GraphicEffect) isSubaction() {}

// 0x0B
type CreateHitbox struct {
	ID                  uint32 `bit:"3"`
	HitGroup            uint32 `bit:"3"`
	_                   uint32 `bit:"1"`
	BoneId              uint32 `bit:"8"`
	UseCommonBoneIDs    bool   `bit:"1"`
	Damage              uint32 `bit:"10"`
	Size                uint32 `bit:"16"`
	Z                   int32  `bit:"16"`
	Y                   int32  `bit:"16"`
	X                   int32  `bit:"16"`
	Angle               uint32 `bit:"9"`
	KnockbackGrowth     uint32 `bit:"9"`
	SetKnockback        uint32 `bit:"9"`
	ItemHitInteraction  bool   `bit:"1"`
	IgnoreThrownFighter bool   `bit:"1"`
	IgnoreFighterScale  bool   `bit:"1"`
	Clank               bool   `bit:"1"`
	Rebound             bool   `bit:"1"`
	BaseKnockback       uint32 `bit:"9"`
	Element             uint32 `bit:"5"`
	ShieldDamage        int32  `bit:"8"`
	HitSFXSeverity      uint32 `bit:"3"`
	HitSFXKind          uint32 `bit:"5"`
	HitGroundedFighters bool   `bit:"1"`
	HitAirBorneFighters bool   `bit:"1"`
}

func (CreateHitbox) isSubaction() {}

// 0x0C
type UpdateHitboxDamage struct {
	ID     uint32 `bit:"3"`
	Damage uint32 `bit:"23"`
}

func (UpdateHitboxDamage) isSubaction() {}

// 0x0D
type UpdateHitboxSize struct {
	ID   uint32 `bit:"3"`
	Size uint32 `bit:"23"`
}

func (UpdateHitboxSize) isSubaction() {}

// 0x0E
type SetHitboxHitGroundAir struct {
	ID    uint32 `bit:"24"`
	Type  bool   `bit:"1"`
	Value bool   `bit:"1"`
}

func (SetHitboxHitGroundAir) isSubaction() {}

// 0x0F
type RemoveHitbox struct {
	ID uint32 `bit:"26"`
}

func (RemoveHitbox) isSubaction() {}

// 0x10
type RemoveAllHitboxes struct {
	_ uint32 `bit:"26"`
}

func (RemoveAllHitboxes) isSubaction() {}

// 0x11
type SoundEffect struct {
	Behavior uint32 `bit:"8"`
	_        uint32 `bit:"18"`
	SFXID    uint32 `bit:"32"`
	_        uint32 `bit:"16"`
	Volume   uint32 `bit:"8"`
	_        uint32 `bit:"8"`
}

func (SoundEffect) isSubaction() {}

// 0x12
type RandomSmashSFX struct {
	Value uint32 `bit:"26"`
}

func (RandomSmashSFX) isSubaction() {}

// 0x13
type SetFlag struct {
	Flag  uint32 `bit:"2"`
	Value uint32 `bit:"24"`
}

func (SetFlag) isSubaction() {}

// 0x14
type AllowReverseDirectionAndRapidJabEnd struct {
	Value uint32 `bit:"26"`
}

func (AllowReverseDirectionAndRapidJabEnd) isSubaction() {}

// 0x15
type SetFlag0x2210to10 struct {
	_ uint32 `bit:"26"`
}

func (SetFlag0x2210to10) isSubaction() {}

// 0x16
type SetFlag0x2210to20 struct {
	_ uint32 `bit:"26"`
}

func (SetFlag0x2210to20) isSubaction() {}

// 0x17
type AllowInterrupt struct {
	_ uint32 `bit:"26"`
}

func (AllowInterrupt) isSubaction() {}

// 0x18
type ProjectileFlag struct {
	Flag uint32 `bit:"26"`
}

func (ProjectileFlag) isSubaction() {}

// 0x19
type SetJumpState struct {
	State uint32 `bit:"26"`
}

func (SetJumpState) isSubaction() {}

// 0x1A
type SetBodyCollisionState struct {
	State uint32 `bit:"26"`
}

func (SetBodyCollisionState) isSubaction() {}

// 0x1B
type SetAllBoneCollisionState struct {
	State uint32 `bit:"26"`
}

func (SetAllBoneCollisionState) isSubaction() {}

// 0x1C
type SetBoneCollisionState struct {
	ID    uint32 `bit:"8"`
	State uint32 `bit:"18"`
}

func (SetBoneCollisionState) isSubaction() {}

// 0x1D
type AllowJabFollowUp struct {
	Value uint32 `bit:"26"`
}

func (AllowJabFollowUp) isSubaction() {}

// 0x1E
type SetRapidJabFlag struct {
	Value uint32 `bit:"26"`
}

func (SetRapidJabFlag) isSubaction() {}

// 0x1F
type ChangeModelState struct {
	StructID uint32 `bit:"7"`
	_        uint32 `bit:"1"`
	ObjectID int32  `bit:"18"`
}

func (ChangeModelState) isSubaction() {}

// 0x20
type RevertModels struct {
	_ uint32 `bit:"26"`
}

func (RevertModels) isSubaction() {}

// 0x21
type RemoveModels struct {
	_ uint32 `bit:"26"`
}

func (RemoveModels) isSubaction() {}

// 0x22
type Throw struct {
	Type            uint32 `bit:"3"`
	_               uint32 `bit:"14"`
	Damage          uint32 `bit:"9"`
	Angle           uint32 `bit:"9"`
	KnockbackGrowth uint32 `bit:"9"`
	SetKnockback    uint32 `bit:"9"`
	_               uint32 `bit:"5"`
	BaseKnockback   uint32 `bit:"9"`
	Element         uint32 `bit:"4"`
	_               uint32 `bit:"19"`
}

func (Throw) isSubaction() {}

// 0x23
type HeldItemInvisibility struct {
	_    uint32 `bit:"25"`
	Flag uint32 `bit:"1"`
}

func (HeldItemInvisibility) isSubaction() {}

// 0x24
type BodyArticleInvisibility struct {
	_    uint32 `bit:"25"`
	Flag uint32 `bit:"1"`
}

func (BodyArticleInvisibility) isSubaction() {}

// 0x25
type CharacterInvisibility struct {
	_    uint32 `bit:"25"`
	Flag uint32 `bit:"1"`
}

func (CharacterInvisibility) isSubaction() {}

// 0x26
type PseudoRandomSoundEffect struct {
	_ uint32 `bit:"26"`
	_ uint32 `bit:"32"`
	_ uint32 `bit:"32"`
	_ uint32 `bit:"32"`
	_ uint32 `bit:"32"`
	_ uint32 `bit:"32"`
	_ uint32 `bit:"32"`
}

func (PseudoRandomSoundEffect) isSubaction() {}

// 0x27
type Unknown0x27 struct {
	_     uint32 `bit:"26"`
	_     uint32 `bit:"8"`
	_     uint32 `bit:"8"`
	_     uint32 `bit:"8"`
	SFXID uint32 `bit:"32"`
	_     uint32 `bit:"32"`
	_     uint32 `bit:"8"`
	_     uint32 `bit:"8"`
	_     uint32 `bit:"8"`
	_     uint32 `bit:"8"`
}

func (Unknown0x27) isSubaction() {}

// 0x28
type AnimateTexture struct {
	MaterialFlag  uint32 `bit:"1"`
	MaterialIndex uint32 `bit:"7"`
	FrameFlag     uint32 `bit:"7"`
	Frame         uint32 `bit:"11"`
}

func (AnimateTexture) isSubaction() {}

// 0x29
type AnimateModel struct {
	BodyPart uint32 `bit:"7"`
	State    uint32 `bit:"7"`
	_        uint32 `bit:"12"`
}

func (AnimateModel) isSubaction() {}

// 0x2A
type ChangeHeldItemsActionState struct {
	Value1 uint32 `bit:"13"`
	Value2 uint32 `bit:"13"`
}

func (ChangeHeldItemsActionState) isSubaction() {}

// 0x2B
type Rumble struct {
	Flag   bool   `bit:"1"`
	Value1 uint32 `bit:"12"`
	Value2 uint32 `bit:"13"`
}

func (Rumble) isSubaction() {}

// 0x2C
type SetFlag0x221Eto20 struct {
	_    uint32 `bit:"25"`
	Flag bool   `bit:"1"`
}

func (SetFlag0x221Eto20) isSubaction() {}

// 0x2D
type BodyAura struct {
	_        uint32 `bit:"8"`
	Duration uint32 `bit:"18"`
}

func (BodyAura) isSubaction() {}

// 0x2E
type ColorAnimation struct {
	ID       uint32 `bit:"8"`
	Duration uint32 `bit:"18"`
}

func (ColorAnimation) isSubaction() {}

// 0x2F
type RemoveColorAnimation struct {
	ID uint32 `bit:"8"`
	_  uint32 `bit:"18"`
}

func (RemoveColorAnimation) isSubaction() {}

// 0x30
type Unknown0xC0 struct {
	_ uint32 `bit:"26"`
}

func (Unknown0xC0) isSubaction() {}

// 0x31
type SwordTrail struct {
	BeamSword    bool   `bit:"1"`
	_            uint32 `bit:"17"`
	RenderStatus uint32 `bit:"8"`
}

func (SwordTrail) isSubaction() {}

// 0x32
type EnableRagDoll struct {
	BoneID uint32 `bit:"26"`
}

func (EnableRagDoll) isSubaction() {}

// 0x33
type SelfDamage struct {
	_ uint32 `bit:"26"`
}

func (SelfDamage) isSubaction() {}

// 0x34
type FootsnapBehavior struct {
	_ uint32 `bit:"26"`
}

func (FootsnapBehavior) isSubaction() {}

// 0x35
type SetFlag0x2225to10 struct {
	_ uint32 `bit:"26"`
}

func (SetFlag0x2225to10) isSubaction() {}

// 0x36
type FootStepSoundAndGraphicEffect struct {
	_ uint32 `bit:"26"`
	_ uint32 `bit:"32"`
	_ uint32 `bit:"32"`
}

func (FootStepSoundAndGraphicEffect) isSubaction() {}

// 0x37
type SoundAndGraphicEffect struct {
	GraphicID uint32 `bit:"26"`
	SoundID   uint32 `bit:"32"`
	_         uint32 `bit:"32"`
}

func (SoundAndGraphicEffect) isSubaction() {}

// 0x38
type StartSmashCharge struct {
	_              uint32 `bit:"2"`
	ChargeFrames   uint32 `bit:"8"`
	ChargeRate     uint32 `bit:"16"`
	ColorAnimation uint32 `bit:"8"`
	_              uint32 `bit:"24"`
}

func (StartSmashCharge) isSubaction() {}

// 0x39
type Unknown0x39 struct {
	_ uint32 `bit:"26"`
}

func (Unknown0x39) isSubaction() {}

// 0x3A
type WindEffect struct {
	_ uint32 `bit:"26"`
	_ uint32 `bit:"32"`
	_ uint32 `bit:"32"`
	_ uint32 `bit:"32"`
}

func (WindEffect) isSubaction() {}

// Item Specific (include projectiles, etc)

// 0x0B
type CreateHitboxItem struct {
	ID                   uint32 `bit:"3"`
	_                    uint32 `bit:"3"`
	BoneId               uint32 `bit:"7"`
	Damage               uint32 `bit:"13"`
	Size                 uint32 `bit:"16"`
	Z                    int32  `bit:"16"`
	Y                    int32  `bit:"16"`
	X                    int32  `bit:"16"`
	Angle                uint32 `bit:"9"`
	KnockbackGrowth      uint32 `bit:"9"`
	SetKnockback         uint32 `bit:"9"`
	_                    uint32 `bit:"5"`
	BaseKnockback        uint32 `bit:"9"`
	Element              uint32 `bit:"5"`
	Clank                bool   `bit:"1"`
	ShieldDamage         int32  `bit:"8"`
	HitSFXSeverity       uint32 `bit:"3"`
	HitSFXKind           uint32 `bit:"4"`
	HitGroundedFighters  bool   `bit:"1"`
	HitAirBorneFighters  bool   `bit:"1"`
	HitCooldown          uint32 `bit:"8"`
	TimedRehitNonFighter bool   `bit:"1"`
	TimedRehitFighter    bool   `bit:"1"`
	TimedRehitShield     bool   `bit:"1"`
	Reflectable          bool   `bit:"1"`
	Absorbable           bool   `bit:"1"`
	Shieldable           bool   `bit:"1"`
	_                    bool   `bit:"1"`
	Deflectable          bool   `bit:"1"`
	Reflectable2         bool   `bit:"1"`
	_                    uint32 `bit:"15"`
}

func (CreateHitboxItem) isSubaction() {}

type Unknown0x0FItem struct {
	_ uint32 `bit:"26"`
}

func (Unknown0x0FItem) isSubaction() {}

func subactionTypeSwitch(i uint8, isItem bool) (SubAction, error) {
	if isItem {
		return itemSubactionTypeSwitch(i)
	}
	return fighterSubactionTypeSwitch(i)
}

func fighterSubactionTypeSwitch(i uint8) (SubAction, error) {
	switch i {
	case 0x00:
		return &EndOfScript{}, nil
	case 0x01:
		return &SynchronousTimer{}, nil
	case 0x02:
		return &AsynchronousTimer{}, nil
	case 0x03:
		return &SetLoop{}, nil
	case 0x04:
		return &ExecuteLoop{}, nil
	case 0x05:
		return &Subroutine{}, nil
	case 0x06:
		return &SubroutineReturn{}, nil
	case 0x07:
		return &GoTo{}, nil
	case 0x08:
		return &SetTimerAnimation{}, nil
	case 0x09:
		return &Unknown0x09{}, nil
	case 0x0A:
		return &GraphicEffect{}, nil
	case 0x0B:
		return &CreateHitbox{}, nil
	case 0x0C:
		return &UpdateHitboxDamage{}, nil
	case 0x0D:
		return &UpdateHitboxSize{}, nil
	case 0x0E:
		return &SetHitboxHitGroundAir{}, nil
	case 0x0F:
		return &RemoveHitbox{}, nil
	case 0x10:
		return &RemoveAllHitboxes{}, nil
	case 0x11:
		return &SoundEffect{}, nil
	case 0x12:
		return &RandomSmashSFX{}, nil
	case 0x13:
		return &SetFlag{}, nil
	case 0x14:
		return &AllowReverseDirectionAndRapidJabEnd{}, nil
	case 0x15:
		return &SetFlag0x2210to10{}, nil
	case 0x16:
		return &SetFlag0x2210to20{}, nil
	case 0x17:
		return &AllowInterrupt{}, nil
	case 0x18:
		return &ProjectileFlag{}, nil
	case 0x19:
		return &SetJumpState{}, nil
	case 0x1A:
		return &SetBodyCollisionState{}, nil
	case 0x1B:
		return &SetAllBoneCollisionState{}, nil
	case 0x1C:
		return &SetBoneCollisionState{}, nil
	case 0x1D:
		return &AllowJabFollowUp{}, nil
	case 0x1E:
		return &SetRapidJabFlag{}, nil
	case 0x1F:
		return &ChangeModelState{}, nil
	case 0x20:
		return &RevertModels{}, nil
	case 0x21:
		return &RemoveModels{}, nil
	case 0x22:
		return &Throw{}, nil
	case 0x23:
		return &HeldItemInvisibility{}, nil
	case 0x24:
		return &BodyArticleInvisibility{}, nil
	case 0x25:
		return &CharacterInvisibility{}, nil
	case 0x26:
		return &PseudoRandomSoundEffect{}, nil
	case 0x27:
		return &Unknown0x27{}, nil
	case 0x28:
		return &AnimateTexture{}, nil
	case 0x29:
		return &AnimateModel{}, nil
	case 0x2A:
		return &ChangeHeldItemsActionState{}, nil
	case 0x2B:
		return &Rumble{}, nil
	case 0x2C:
		return &SetFlag0x221Eto20{}, nil
	case 0x2D:
		return &BodyAura{}, nil
	case 0x2E:
		return &ColorAnimation{}, nil
	case 0x2F:
		return &RemoveColorAnimation{}, nil
	case 0x30:
		return &Unknown0xC0{}, nil
	case 0x31:
		return &SwordTrail{}, nil
	case 0x32:
		return &EnableRagDoll{}, nil
	case 0x33:
		return &SelfDamage{}, nil
	case 0x34:
		return &FootsnapBehavior{}, nil
	case 0x35:
		return &SetFlag0x2225to10{}, nil
	case 0x36:
		return &FootStepSoundAndGraphicEffect{}, nil
	case 0x37:
		return &SoundAndGraphicEffect{}, nil
	case 0x38:
		return &StartSmashCharge{}, nil
	case 0x39:
		return &Unknown0x39{}, nil
	case 0x3A:
		return &WindEffect{}, nil
	default:
		return &emptySubAction{}, &SubActionNotImplemented{Id: i, IsItem: false}
	}
}

func itemSubactionTypeSwitch(i uint8) (SubAction, error) {
	switch i {
	case 0x00:
		return &EndOfScript{}, nil
	case 0x01:
		return &SynchronousTimer{}, nil
	case 0x02:
		return &AsynchronousTimer{}, nil
	case 0x03:
		return &SetLoop{}, nil
	case 0x04:
		return &ExecuteLoop{}, nil
	case 0x05:
		return &Subroutine{}, nil
	case 0x06:
		return &SubroutineReturn{}, nil
	case 0x07:
		return &GoTo{}, nil
	case 0x08:
		return &SetTimerAnimation{}, nil
	case 0x09:
		return &Unknown0x09{}, nil
	case 0x0B:
		return &CreateHitboxItem{}, nil
	case 0x0C:
		return &UpdateHitboxDamage{}, nil
	case 0x0D:
		return &UpdateHitboxSize{}, nil
	case 0x0E:
		return &RemoveHitbox{}, nil
	case 0x0F:
		return &Unknown0x0FItem{}, nil
	default:
		return &emptySubAction{}, &SubActionNotImplemented{Id: i, IsItem: true}
	}
}

func GetSubActionType(r *binread.Reader, isItem bool) (SubAction, error) {
	enrichError := func(subac SubAction, err error) (SubAction, error) {
		errMsg := fmt.Sprintf("GetSubActionType(Offset:0x%X): ", r.CurrentPosition())
		return subac, fmt.Errorf(errMsg+"%w", err)
	}
	b, err := r.Peek(1)
	if err != nil {
		return enrichError(&emptySubAction{}, err)
	}
	subactionType, err := subactionTypeSwitch(uint8(b[0])>>2, isItem)
	if err != nil {
		return enrichError(subactionType, err)
	}
	return subactionType, nil
}

func DecodeSubAction(r *binread.Reader, s SubAction) error {
	numbits, err := binread.NumBits(s)
	if err != nil {
		return err
	}

	// if numbits%32 != 0 {
	//         return errors.New(fmt.Sprintf("SubAction total bit number +6 should be a multiple of 32. SubAction: %T | Numbits+6=%d", s, numbits))
	// }

	byts := make([]byte, numbits/8)
	err = r.Decode(&byts)
	if err != nil {
		return err
	}

	bits := binread.SplitBytes(byts)
	return binread.BitRead(bits, s, 6)
}

type SubActions []SubAction

func (a *SubActions) BinRead(r *binread.Reader, args ...Args) error {
	var offset Addr
	var err error

	isItem := false

	ok := false
	for _, args := range args {
		offset_, ok_ := args["offset"].(Addr)
		if ok_ {
			offset = offset_
			ok = true
			continue
		}
		// offset_, ok := args["offset"].(Addr)
		// if ok {
		// 	offset = offset_
		// 	continue
		// }
		isItem_, ok_ := args["isItem"].(bool)
		if ok_ {
			isItem = isItem_
			continue
		}
	}

	if !ok {
		err = r.Decode(&offset)
		if err != nil {
			return err
		}
	}

	before := r.CurrentPosition()
	_, err = r.Seek(offset.ToSeek(), io.SeekStart)
	if err != nil {
		return err
	}

	subactions := make([]SubAction, 0)

subactionloop:
	for {
		subac, err := GetSubActionType(r, isItem)
		if err != nil {
			return fmt.Errorf("Subactions(Offset:0x%X): %w", offset.ToSeek(), err)
			// if _, ok := err.(*SubActionNotImplemented); ok {
			// 	var x [4]byte
			// 	r.Decode(&x)
			// 	continue
			// } else {
			// 	return enrichError(err)
			// }
		}

		err = DecodeSubAction(r, subac)
		if err != nil {
			return err
		}

		subactions = append(subactions, subac)
		switch subac.(type) {
		case *EndOfScript, *GoTo, *SubroutineReturn:
			break subactionloop
		default:
		}
	}

	_, err = r.Seek(before, io.SeekStart)
	if err != nil {
		return err
	}

	*a = subactions
	return nil
}
