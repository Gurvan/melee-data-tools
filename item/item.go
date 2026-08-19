package item

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/Gurvan/melee-data-tools/binread"
	"github.com/Gurvan/melee-data-tools/descriptor"
	. "github.com/Gurvan/melee-data-tools/lib"
	"github.com/Gurvan/melee-data-tools/model"
	"github.com/Gurvan/melee-data-tools/subaction"
)

type Item struct {
	Attributes Ptr[Attributes]

	// SpecificAttributes points at the article's own attribute block, the article-level counterpart
	// of a fighter's special attributes. Null (0x20 once decoded) for articles that have none; the
	// layout differs per article, so only known (fighter, slot) pairs are decoded.
	SpecificAttributes Addr

	Hurtboxes HurtboxList
	States    States
	Model     Ptr[Model]
	_         Addr // Dynamics

	// Vegetable holds the decoded specific attributes of Peach's vegetable article, nil for every
	// other article. It is not read from this position, so the binary decoder must skip it.
	Vegetable *VegetableAttributes `binread:"-"`
}

// Lifetime reports how many frames the article lasts before it expires. The second result is false
// when the article carries no override, in which case the game's global item default applies.
func (i *Item) Lifetime() (float32, bool) {
	if i.Vegetable != nil {
		return i.Vegetable.Duration, true
	}
	return 0, false
}

// VegetableFace is one row of the table Peach's vegetable draws its face from: Weight is the row's
// share of the draw (a proportion, not a percentage), Damage what that face deals when thrown.
type VegetableFace struct {
	Weight int32
	Damage int32
}

// VegetableAttributes is the specific-attributes block of Peach's vegetable article: how long a
// pulled vegetable lasts before it vanishes, followed by the face table. The row count is stored
// inline ahead of the rows rather than as the usual offset/size pair, so this needs its own decoder.
type VegetableAttributes struct {
	Duration float32
	Faces    []VegetableFace
}

var _ binread.BinReader = (*VegetableAttributes)(nil)

func (a *VegetableAttributes) BinRead(r *binread.Reader, args ...Args) error {
	if err := r.Decode(&a.Duration, args...); err != nil {
		return err
	}

	var rowCount int32
	if err := r.Decode(&rowCount, args...); err != nil {
		return err
	}
	// The row count comes straight out of the file, so validate it before allocating on it.
	if rowCount < 0 || int64(rowCount)*8 > r.Size() {
		return fmt.Errorf("vegetable face table claims %d rows, which does not fit the file", rowCount)
	}

	a.Faces = make([]VegetableFace, rowCount)
	return r.Decode(&a.Faces, args...)
}

// decodeSpecificAttributes decodes the article's specific-attributes block for the articles whose
// layout is known: opt-in per (fighter, slot), and an unlisted article keeps only its raw pointer.
func decodeSpecificAttributes(r *binread.Reader, it *Item, fighterName string, slot int, args ...Args) error {
	if it.SpecificAttributes == Addr(0x20) {
		return nil
	}

	switch {
	case fighterName == "Peach" && slot == 1:
		before := r.CurrentPosition()
		if _, err := r.Seek(it.SpecificAttributes.ToSeek(), io.SeekStart); err != nil {
			return err
		}
		var attributes VegetableAttributes
		if err := r.Decode(&attributes, args...); err != nil {
			return err
		}
		it.Vegetable = &attributes
		if _, err := r.Seek(before, io.SeekStart); err != nil {
			return err
		}
	}

	return nil
}

type ItemFlags struct {
	IsHeavy bool  `bit:"1"`
	_       uint8 `bit:"4"`
	// HoldKind is how the article is meant to be carried, and is what picks the pose the
	// hand carrying it takes: none, open inwards, sword, open downwards or open forwards.
	HoldKind uint8 `bit:"3"`
}

func (f *ItemFlags) BinRead(r *binread.Reader, args ...Args) error {
	var err error

	byts := make([]byte, 4)
	err = r.Decode(&byts)
	if err != nil {
		return err
	}

	bits := binread.SplitBytes(byts)
	err = binread.BitRead(bits, f, 0)
	return err
}

type ECB struct {
	Top    float32
	Bottom float32
	Left   float32
	Right  float32
}

// GrabBox is where the item offers itself to be grabbed: a point off its own middle plus a box
// added to the reach of whoever is trying to take it. All zero on articles nobody can carry.
type GrabBox struct {
	OffsetX    float32
	OffsetY    float32
	HalfWidth  float32
	HalfHeight float32
}

type Attributes struct {
	Flags                ItemFlags
	ThrowSpeedMultiplier float32
	_                    float32
	SpinSpeed            float32
	FallAcceleration     float32
	MaxFallSpeed         float32
	_                    [6]float32 // Unk0x18
	GrabBox              GrabBox
	ECB                  ECB
	_                    [4]float32 // Unk0x50
	ModelScale           float32
	_                    [8]int32 // SFXs
}

type Hurtbox struct {
	BoneIndex uint32
	BaseX     float32
	BaseY     float32
	BaseZ     float32
	TipX      float32
	TipY      float32
	TipZ      float32
	Radius    float32
}

// HurtboxList is how an item's hurtboxes hang off it: a single pointer to a {count, rows} pair.
// Every fighter article in the roster leaves the pointer null; the common items carry real ones.
type HurtboxList []Hurtbox

var _ binread.BinReader = (*HurtboxList)(nil)

func (h *HurtboxList) BinRead(r *binread.Reader, args ...Args) error {
	var offset Addr
	if err := r.Decode(&offset); err != nil {
		return err
	}
	if offset == Addr(0x20) {
		*h = nil
		return nil
	}

	before := r.CurrentPosition()
	if _, err := r.Seek(offset.ToSeek(), io.SeekStart); err != nil {
		return err
	}

	var count uint32
	if err := r.Decode(&count); err != nil {
		return err
	}
	// The count comes straight out of the file, so validate it before allocating on it.
	if int64(count)*0x20 > r.Size() {
		return fmt.Errorf("item hurtbox list claims %d rows, which does not fit the file", count)
	}
	var rows Addr
	if err := r.Decode(&rows); err != nil {
		return err
	}
	if _, err := r.Seek(rows.ToSeek(), io.SeekStart); err != nil {
		return err
	}

	out := make([]Hurtbox, count)
	if err := r.Decode(&out); err != nil {
		return err
	}

	if _, err := r.Seek(before, io.SeekStart); err != nil {
		return err
	}

	*h = out
	return nil
}

const StateSize = 0x10

type State struct {
	_          [3]Addr
	Subactions subaction.SubActions
}

type States []State

func (s *States) BinRead(r *binread.Reader, args ...Args) error {
	var err error
	var offset Addr

	err = r.Decode(&offset)
	if err != nil {
		return err
	}

	var numElem int
	for _, args := range args {
		if reloc, ok := args["relocation"].(*Relocation); ok && reloc != nil {
			numElem = int((*reloc)[offset]) / StateSize
			// fmt.Println("Num elem states:", numElem, elemSize, offset, statesSize, ok)
		}
	}
	before := r.CurrentPosition()
	_, err = r.Seek(offset.ToSeek(), io.SeekStart)
	if err != nil {
		return err
	}

	if args == nil {
		args = make([]Args, 1)
		args[0] = make(Args)
	}
	args[0]["isItem"] = true

	states := make([]State, 0)
	for i := 0; i < numElem; i++ {
		var state State
		err = r.Decode(&state, args...)
		if err != nil {
			return fmt.Errorf("State(Index:%d): %w", i, err)
		}
		states = append(states, state)
	}
	_, err = r.Seek(before, io.SeekStart)
	if err != nil {
		return err
	}

	*s = states
	return nil
}

type Model struct {
	// Optional: some articles carry no model of their own (Peach's vegetable uses a common one).
	Joint        OptionalPtr[model.Joint]
	BoneCount    int32
	BoneAttachID int32
	BitField     int32
}

// CommonItem is one entry of the game's shared item table: an item any fighter can come to hold,
// keyed by the kind replays record as the item's type.
type CommonItem struct {
	Kind int
	Item Item
}

// CommonItems is the shared item table of the common-item file. Kinds decode opt-in per verified
// item, like the per-fighter article slots.
type CommonItems []CommonItem

func (a *CommonItems) BinRead(r *binread.Reader, args ...Args) error {
	var offset Addr
	if err := r.Decode(&offset); err != nil {
		return err
	}

	before := r.CurrentPosition()
	if _, err := r.Seek(offset.ToSeek(), io.SeekStart); err != nil {
		return err
	}
	startPos := r.CurrentPosition()

	items := make([]CommonItem, 0)
	for _, kind := range commonItemSwitch() {
		var t OptionalPtr[Item]
		if _, err := r.Seek(startPos+int64(kind)*0x4, io.SeekStart); err != nil {
			return err
		}
		if err := r.Decode(&t, args...); err != nil {
			return fmt.Errorf("CommonItem(Kind:%d): %w", kind, err)
		}
		if t.ValuePtr == nil {
			continue
		}
		items = append(items, CommonItem{Kind: kind, Item: *t.ValuePtr})
	}

	if _, err := r.Seek(before, io.SeekStart); err != nil {
		return err
	}

	*a = items
	return nil
}

// commonItemSwitch lists the common items that decode: 6 is the bomb.
func commonItemSwitch() []int {
	return []int{6}
}

type Items []Item

func (a *Items) BinRead(r *binread.Reader, args ...Args) error {
	var err error
	var offset Addr

	err = r.Decode(&offset)
	if err != nil {
		return err
	}

	before := r.CurrentPosition()
	_, err = r.Seek(offset.ToSeek(), io.SeekStart)
	if err != nil {
		return err
	}

	var firstRoot string
	for _, args := range args {
		if desc, ok := args["descriptor"].(*descriptor.Descriptor); ok {
			if firstRoot, err = desc.FirstRootName(); err != nil {
				return err
			}
			break
		}
	}

	fighterName, err := fighterName(firstRoot)
	if err != nil {
		return err
	}
	indices := fighterSwitch(fighterName)

	startPos := r.CurrentPosition()

	items := make([]Item, 0)
	for _, i := range indices {
		var t OptionalPtr[Item]
		_, err := r.Seek(startPos+int64(i)*0x4, io.SeekStart)
		err = r.Decode(&t, args...)
		if err != nil {
			return fmt.Errorf("Item(Index:%d): %w", i, err)
		}
		if t.ValuePtr == nil {
			continue
		}
		if err := decodeSpecificAttributes(r, t.ValuePtr, fighterName, i, args...); err != nil {
			return fmt.Errorf("Item(Index:%d): %w", i, err)
		}
		items = append(items, *t.ValuePtr)
	}

	_, err = r.Seek(before, io.SeekStart)
	if err != nil {
		return err
	}

	*a = items
	return nil
}

// fighterName strips the ftData prefix off a fighter data file's first root name, and fails on any
// file that is not a fighter data file.
func fighterName(firstRoot string) (string, error) {
	name := strings.TrimPrefix(firstRoot, "ftData")
	if name == firstRoot {
		return "", errors.New(fmt.Sprintf("File first root %s does not belong to fighter data file.\n", firstRoot))
	}
	return name, nil
}

func fighterSwitch(name string) []int {
	switch name {
	case "Fox":
		return []int{0, 1, 2}
	case "Falco":
		return []int{0, 1, 3}
	case "Peach":
		// 5 slots; 0 and 4 are hitbox-only articles with no model.
		return []int{0, 1, 2, 3, 4}
	default:
		return []int{}
	}
}
