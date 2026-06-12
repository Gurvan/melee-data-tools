package attributes

import "fmt"

type Hex32 uint32

func (h Hex32) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", h.String())), nil
}

func (h Hex32) String() string {
	return fmt.Sprintf("0x%08x", uint32(h))
}
