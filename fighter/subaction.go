package fighter

type EmptySubAction struct{}

func (s EmptySubAction) IsSubAction() bool {
	return true
}

type SubAction interface {
	IsSubAction() bool
}

var _ SubAction = (*EmptySubAction)(nil)

type EndOfScript struct {
	EmptySubAction
	_ uint32 `bit:"26"`
}

type SynchronousTimer struct {
	EmptySubAction
	frame uint32 `bit:"26"`
}

func subactionTypeSwitch(i uint8) SubAction {
	switch i {
	case 0:
		return EndOfScript{}
	default:
		return EmptySubAction{}
	}
}
