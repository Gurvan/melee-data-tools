package fighter

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/Gurvan/melee-data-tools/binread"
)

type emptySubAction struct{}

type SubActionNotImplemented struct {
	Id uint8
}

func (e *SubActionNotImplemented) Error() string {
	return fmt.Sprintf("SubAction 0x%X is not implemented", e.Id)
}

// func (s emptySubAction) IsSubAction() bool {
//         return true
// }

type SubAction interface {
	// IsSubAction() bool
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

// 0x01
type SynchronousTimer struct {
	// emptySubAction
	Frame uint32 `bit:"26"`
}

// 0x02
type AsynchronousTimer struct {
	// emptySubAction
	Frame uint32 `bit:"26"`
}

// 0x03
type SetLoop struct {
	Count uint32 `bit:"26"`
}

// 0x04
type ExecuteLoop struct {
	_ uint32 `bit:"26"`
}

// 0x05
type Subroutine struct {
	Target  uint32 `bit:"26"`
	Pointer uint32 `bit:"32"`
}

// 0x06
type SubroutineReturn struct {
	_ uint32 `bit:"26"`
}

// 0x07
type GoTo struct {
	Target  uint32 `bit:"26"`
	Pointer uint32 `bit:"32"`
}

// 0x08
type SetTimerAnimation struct {
	_ uint32 `bit:"26"`
}

// 0x09
type Unknown0x09 struct {
	_ uint32 `bit:"26"`
}

// 0x0A
type GraphicEffect struct {
	BoneId               uint32 `bit:"32"`
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
	_                    uint32 `bit:"8"`
}

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

// 0x0C
type AdjusteHitboxDamage struct {
	ID     uint32 `bit:"3"`
	Damage uint32 `bit:"23"`
}

// 0x0D
type AdjusteHitboxSize struct {
	ID   uint32 `bit:"3"`
	Size uint32 `bit:"23"`
}

// 0x0E
type SetHitboxInteraction struct {
	ID    uint32 `bit:"24"`
	Type  bool   `bit:"1"`
	Value bool   `bit:"1"`
}

// 0x0F
type RemoveHitbox struct {
	ID uint32 `bit:"26"`
}

// 0x10
type RemoveAllHitboxes struct {
	_ uint32 `bit:"26"`
}

// 0x11
type SoundEffect struct {
	Behavior uint32 `bit:"8"`
	_        uint32 `bit:"18"`
	SFXID    uint32 `bit:"32"`
	_        uint32 `bit:"16"`
	Volume   uint32 `bit:"8"`
	_        uint32 `bit:"8"`
}

// 0x12
type RandomSmashSFX struct {
	Value uint32 `bit:"26"`
}

// 0x13
type SetFlag struct {
	Flag  uint32 `bit:"2"`
	Value uint32 `bit:"24"`
}

// 0x14
type ReverseDirection struct {
	Enabled uint32 `bit:"26"`
}

// 0x15
type SetFlag0x2210to10 struct {
	_ uint32 `bit:"26"`
}

// 0x16
type SetFlag0x2210to20 struct {
	_ uint32 `bit:"26"`
}

// 0x17
type AllowInterrupt struct {
	_ uint32 `bit:"26"`
}

// 0x18
type ProjectileFlag struct {
	_ uint32 `bit:"26"`
}

// 0x19
type SetJumpState struct {
	_ uint32 `bit:"26"`
}

// 0x1A
type SetBodyCollisionState struct {
	State uint32 `bit:"26"`
}

// 0x1B
type SetAllBoneCollisionState struct {
	State uint32 `bit:"26"`
}

// 0x1C
type SetBoneCollisionState struct {
	ID    uint32 `bit:"8"`
	State uint32 `bit:"18"`
}

// 0x1D
type EnableJabFollowUp struct {
	_ uint32 `bit:"26"`
}

// 0x1E
type SetRapidJabFlag struct {
	Value uint32 `bit:"26"`
}

// 0x1F
type ChangeModelState struct {
	StructID uint32 `bit:"7"`
	_        uint32 `bit:"1"`
	ObjectID int32  `bit:"18"`
}

// 0x20
type RevertModels struct {
	_ uint32 `bit:"26"`
}

// 0x21
type RemoveModels struct {
	_ uint32 `bit:"26"`
}

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

// 0x23
type HeldItemInvisibility struct {
	_    uint32 `bit:"25"`
	Flag uint32 `bit:"1"`
}

// 0x24
type BodyArticleInvisibility struct {
	_    uint32 `bit:"25"`
	Flag uint32 `bit:"1"`
}

// 0x25
type CharacterInvisibility struct {
	_    uint32 `bit:"25"`
	Flag uint32 `bit:"1"`
}

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

// 0x27
type Unknown0x27 struct {
	_     uint32 `bit:"26"`
	_     uint32 `bit:"32"`
	_     uint32 `bit:"32"`
	_     uint32 `bit:"32"`
	SFXID uint32 `bit:"32"`
	_     uint32 `bit:"32"`
	_     uint32 `bit:"32"`
	_     uint32 `bit:"32"`
	_     uint32 `bit:"32"`
	_     uint32 `bit:"32"`
}

// 0x28
type AnimateTexture struct {
	MaterialFlag  uint32 `bit:"1"`
	MaterialIndex uint32 `bit:"7"`
	FrameFlag     uint32 `bit:"7"`
	Frame         uint32 `bit:"11"`
}

// 0x29
type AnimateModel struct {
	BodyPart uint32 `bit:"7"`
	State    uint32 `bit:"7"`
	_        uint32 `bit:"12"`
}

// 0x2A
type ChangeHeldItemsActionState struct {
	Value1 uint32 `bit:"13"`
	Value2 uint32 `bit:"13"`
}

// 0x2B
type Rumble struct {
	Flag   bool   `bit:"1"`
	Value1 uint32 `bit:"12"`
	Value2 uint32 `bit:"13"`
}

// 0x2C
type SetFlag0x221Eto20 struct {
	_    uint32 `bit:"25"`
	Flag bool   `bit:"1"`
}

// 0x2D
type BodyAura struct {
	_        uint32 `bit:"8"`
	Duration uint32 `bit:"18"`
}

// 0x2E
type ColorAnimation struct {
	ID       uint32 `bit:"8"`
	Duration uint32 `bit:"18"`
}

// 0x2F
type RemoveColorAnimation struct {
	ID uint32 `bit:"8"`
	_  uint32 `bit:"18"`
}

// 0x30
type Unknown0xC0 struct {
	_ uint32 `bit:"26"`
}

// 0x31
type SwordTrail struct {
	BeamSword    bool   `bit:"1"`
	_            uint32 `bit:"17"`
	RenderStatus uint32 `bit:"8"`
}

// 0x32
type EnableRagDoll struct {
	BoneID uint32 `bit:"26"`
}

// 0x33
type SelfDamage struct {
	_ uint32 `bit:"26"`
}

// 0x34
type FootsnapBehavior struct {
	_ uint32 `bit:"26"`
}

// 0x35
type SetFlag0x2225to10 struct {
	_ uint32 `bit:"26"`
}

// 0x36
type FootStepSoundAndGraphicEffect struct {
	_ uint32 `bit:"26"`
	_ uint32 `bit:"32"`
	_ uint32 `bit:"32"`
}

// 0x37
type SoundAndGraphicEffect struct {
	GraphicID uint32 `bit:"26"`
	SoundID   uint32 `bit:"32"`
	_         uint32 `bit:"32"`
}

// 0x38
type StartSmashCharge struct {
	_              uint32 `bit:"26"`
	ChargeFrames   uint32 `bit:"8"`
	CargeRate      uint32 `bit:"16"`
	ColorAnimation uint32 `bit:"8"`
	_              uint32 `bit:"32"`
}

// 0x39
type Unknown0x39 struct {
	_ uint32 `bit:"26"`
}

// 0x3A
type WindEffect struct {
	_ uint32 `bit:"26"`
	_ uint32 `bit:"32"`
	_ uint32 `bit:"32"`
	_ uint32 `bit:"32"`
}

func subactionTypeSwitch(i uint8) (SubAction, error) {
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
		return &AdjusteHitboxDamage{}, nil
	case 0x0D:
		return &AdjusteHitboxSize{}, nil
	case 0x0E:
		return &SetHitboxInteraction{}, nil
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
		return &ReverseDirection{}, nil
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
		return &EnableJabFollowUp{}, nil
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
		return &emptySubAction{}, &SubActionNotImplemented{Id: i}
	}
}

func GetSubActionType(r *binread.Reader) (SubAction, error) {
	b, err := r.Peek(1)
	if err != nil {
		return &emptySubAction{}, err
	}
	return subactionTypeSwitch(uint8(b[0]) >> 2)
}

type Bit = byte

func SplitBytes(byts []byte) []Bit {
	bits := make([]Bit, 0)
	for _, b := range byts {
		for i := 0; i < 8; i++ {
			v := 0
			if b&(1<<(7-i)) > 0 {
				v = 1
			}
			bits = append(bits, byte(v))
		}
	}
	return bits
}

func JoinBits(bits []Bit, padto int) []byte {
	byts := make([]byte, 0)
	if padto > len(bits) {
		padding := make([]Bit, padto-len(bits))
		bits = append(padding, bits...)
	} else {
		numbytes := (len(bits)-1)/8 + 1
		padding := make([]Bit, 8*numbytes-len(bits))
		bits = append(padding, bits...)
	}

	var by int
	for i, b := range bits {
		i = i % 8
		if i == 0 {
			by = 0
		}
		by += int(b) * (1 << (7 - i))
		if i == 7 {
			byts = append(byts, byte(by))
		}
	}
	return byts
}

func NumBits(s SubAction) (int, error) {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr {
		return 0, errors.New("BitRead should be called with a pointer to a SubAction.")
	}

	v = v.Elem()
	t := reflect.TypeOf(s).Elem()
	var totalbits int = 6
	var numbits int
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous {
			continue
		}
		if bit, ok := field.Tag.Lookup("bit"); ok {
			fmt.Sscanf(bit, "%d", &numbits)
			totalbits += numbits
			// fmt.Println(field.Name, bit, numbits, totalbits)
			// if v.Field(i).CanSet() {
			//         v.Field(i).Set(reflect.ValueOf(uint32(numbits)))
			// }
		}
	}
	return totalbits, nil
}

func BitRead(bits []Bit, s SubAction) error {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr {
		return errors.New("BitRead should be called with a pointer to a SubAction.")
	}

	v = v.Elem()
	t := reflect.TypeOf(s).Elem()

	var p int = 6
	var numbits int
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous {
			continue
		}
		if bit, ok := field.Tag.Lookup("bit"); ok {
			fmt.Sscanf(bit, "%d", &numbits)
			switch t.Field(i).Type.Kind() {
			case reflect.Uint32:
				byts := JoinBits(bits[p:p+numbits], 32)
				r := bytes.NewReader(byts)
				var value uint32
				err := binary.Read(r, binary.BigEndian, &value)
				if err != nil {
					return err
				}
				if v.Field(i).CanSet() {
					v.Field(i).Set(reflect.ValueOf(value))
				}
			case reflect.Int32:
				byts := JoinBits(bits[p:p+numbits], 32)
				r := bytes.NewReader(byts)
				var value int32
				err := binary.Read(r, binary.BigEndian, &value)
				if err != nil {
					return err
				}
				if v.Field(i).CanSet() {
					v.Field(i).Set(reflect.ValueOf(value))
				}
			case reflect.Bool:
				byts := JoinBits(bits[p:p+numbits], 8)
				r := bytes.NewReader(byts)
				var value bool
				err := binary.Read(r, binary.BigEndian, &value)
				if err != nil {
					return err
				}
				if v.Field(i).CanSet() {
					v.Field(i).Set(reflect.ValueOf(value))
				}
			default:
			}
			p += numbits
		} else {
			return errors.New(fmt.Sprintf("All SubAction fields should have a `bit` number tag. SubAction: %v | Field: %v", t.Name(), field.Name))
		}
	}
	return nil
}

func DecodeSubAction(r *binread.Reader, s SubAction) error {
	// fmt.Printf("%#+v\n", s)
	numbits, err := NumBits(s)
	if err != nil {
		return err
	}

	if numbits%32 != 0 {
		return errors.New(fmt.Sprintf("SubAction total bit number +6 should be a multiple of 32. SubAction: %T | Numbits+6=%d", s, numbits))
	}

	byts := make([]byte, numbits/8)
	err = r.Decode(&byts)
	if err != nil {
		return err
	}

	bits := SplitBytes(byts)
	return BitRead(bits, s)
}
