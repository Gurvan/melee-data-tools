package fighter

// ItemPickupBox is how far a fighter reaches for an item they are trying to pick up: a box centred
// a little ahead of them and a little above their feet, in unscaled character units.
type ItemPickupBox struct {
	OffsetX    float32
	OffsetY    float32
	HalfWidth  float32
	HalfHeight float32
}

// ItemPickupBoxes are the three reaches a fighter has, in the order the data stores them: standing,
// two-handed (crates and the like), and airborne. The game picks the third while a fighter is off
// the ground and the first while they are on it - read out of a memory dump mid-jump, where the
// airborne reach in use was the third box (and it is the one that reaches below the feet, which
// only an airborne fighter can use).
type ItemPickupBoxes struct {
	Ground ItemPickupBox
	Crate  ItemPickupBox
	Air    ItemPickupBox
}
