package fighter

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestSubAction(t *testing.T) {
	t.Logf("%v\n", spew.Sdump(subactionTypeSwitch(0)))
	t.Logf("%v\n", spew.Sdump(subactionTypeSwitch(1)))

	args := map[string]uint32{"a": 1}

	v, ok := args["a"]
	t.Logf("%v | %v\n", v, ok)
	v, ok = args["b"]
	t.Logf("%v | %v\n", v, ok)
}
