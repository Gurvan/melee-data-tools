package fighter

// ItemPickupBox is how far a fighter reaches for an item they are trying to pick up: a box centred
// a little ahead of them and a little above their feet, in unscaled character units.
type ItemPickupBox struct {
	OffsetX    float32
	OffsetY    float32
	HalfWidth  float32
	HalfHeight float32
}

// ItemPickupBoxes are the three reaches a fighter has, in the order the data stores them: airborne,
// two-handed (crates and the like), and standing.
type ItemPickupBoxes struct {
	Air    ItemPickupBox
	Crate  ItemPickupBox
	Ground ItemPickupBox
}
