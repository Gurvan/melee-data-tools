package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Gurvan/melee-data-tools/binread"
	"github.com/Gurvan/melee-data-tools/fighter"
)

func PrettyString(obj interface{}) string {
	bytes, _ := json.MarshalIndent(obj, "  ", "  ")
	return string(bytes)
}

func printSubActions(actionTable fighter.ActionTable) {
	for _, action := range actionTable {
		switch action.Name.Value {
		// case "PlyCaptain5K_Share_ACTION_AttackAirHi_figatree", "PlyCaptain5K_Share_ACTION_RunBrake_figatree":
		default:
			if len(action.Name.Value) > 36 {
				fmt.Printf("%s:\n", action.Name.Value[26:len(action.Name.Value)-9])
			} else {
				fmt.Println("NoName:")
			}
			for _, subac := range action.Subactions.Value {
				// fmt.Printf("\t%T\n", subac)
				// fmt.Printf("\t%#v\n", subac)
				fmt.Printf("\t%s\n", fighter.SubActionString(subac))
			}
			// fmt.Println()
		}
	}
}

func main() {
	file, err := os.Open("file.dat")
	// file, err := os.Open("file2.dat")
	if err != nil {
		log.Fatal(err)
	}

	reader := binread.NewReader(file)

	// h := Header{}
	// h := Descriptor{}
	h := FighterFile{}

	err = reader.Decode(&h)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("%#+v\n", h)
	// spew.Dump(h.Data.ActionTable.Value)

	// fmt.Println(len(h.Data.ActionTable.Value))
	printSubActions(h.Data.ActionTable.Value)
}
