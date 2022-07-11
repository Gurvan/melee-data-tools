package mdt

import (
	"regexp"

	. "github.com/Gurvan/melee-data-tools/common"
	"github.com/Gurvan/melee-data-tools/model"
)

type ModelData struct {
	model.Joint
}

type ModelFile = File[ModelData]

func (m *ModelData) UpdateName(firstRootName string) error {
	re := regexp.MustCompile("Ply.*5K(.*)_Share_joint")
	substrings := re.FindStringSubmatch(firstRootName)
	modelName := "Nr"
	if len(substrings) > 0 && substrings[len(substrings)-1] != "" {
		modelName = substrings[len(substrings)-1]
	}
	// n := len(firstRootName)
	// if n < 14 {
	//         return nil
	// }
	// modelName := firstRootName[n-14 : n-12]
	model := &m.Joint
	name := NullTerminatedString(modelName)
	model.Name.ValuePtr = &name
	return nil
}
