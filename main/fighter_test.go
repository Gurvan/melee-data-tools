package main

import (
	"fmt"
	"log"
	"testing"

	mdt "github.com/Gurvan/melee-data-tools"
	"github.com/davecgh/go-spew/spew"
)

func TestFighter(t *testing.T) {
	// fighterFile, err := os.Open("fighter.dat")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer fighterFile.Close()
	// reader := binread.NewReader(fighterFile)
	//
	// fighter := mdt.FighterFile{}
	//
	// err = reader.Decode(&fighter)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	var fighterFile mdt.FighterFile
	// fighterData, desc, err := fighterFile.ReadFromFile("fighter.dat")
	// fighterData, desc, err := fighterFile.ReadFromFile("../../../re/files/PlCa.dat")
	// fighterData, desc, err := fighterFile.ReadFromFile("../../../re/files/PlMs.dat")
	fighterData, desc, err := fighterFile.ReadFromFile("../../../re/files/PlFc.dat")
	if err != nil {
		fmt.Println(err)
		log.Fatalf("Could not parse %s as FighterFile.\n", "fighter.dat")
	}
	spew.Config.DisableCapacities = true
	spew.Config.DisablePointerAddresses = true

	// fmt.Printf("%#+v\n", h)
	spew.Dump(desc.Header)
	// spew.Dump(fighter.Desc.Roots)
	// spew.Dump(fighter.Desc.Relocation)
	// spew.Dump(fighterData.AttributesCommon)
	// spew.Dump(fighterData.AttributesSpecial)
	// spew.Dump(fighterData.ModelParams)
	spew.Dump(fighterData.Items)

	// spew.Dump(fighterData.Items[0].ValuePtr.States)
	// for _, i := range fighterData.Items {
	// 	if i.ValuePtr == nil {
	// 		continue
	// 	}
	// 	spew.Dump(i.ValuePtr.States)
	// 	// for _, s := range i.ValuePtr.States {
	// 	// 	spew.Dump(s.Subactions)
	// 	// }
	// }
	// spew.Dump(fighterData.ActionTable)
	// for _, a := range fighterData.ActionTable.GetValue() {
	// 	spew.Dump(a.Subactions)
	// }
	t.FailNow()
}
