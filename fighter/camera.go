package fighter

// CameraBox is the room the battle camera keeps around a fighter, before the model scale: a
// centre height above the fighter's position, then signed offsets from that centre. A sixth
// value follows in the file; the battle camera never reads it.
type CameraBox struct {
	CenterHeight float32
	Forward      float32
	Rear         float32 // negative: behind the fighter
	Top          float32
	Bottom       float32 // negative: below the centre
}
