package mdt

import "github.com/Gurvan/melee-data-tools/model"

type ModelData struct {
	model.Joint
}

type ModelFile = File[ModelData]
