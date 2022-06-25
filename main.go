package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/Gurvan/melee-data-tools/binread"
	"github.com/Gurvan/melee-data-tools/fighter"
	"github.com/davecgh/go-spew/spew"
)

func PrettyString(obj interface{}) string {
	bytes, _ := json.MarshalIndent(obj, "  ", "  ")
	return string(bytes)
}

func printActions(actionTable fighter.ActionTable, printSubactions bool) {
	fmt.Println(len(actionTable))
	for index, action := range actionTable {
		switch action.Name.Value {
		// case "PlyCaptain5K_Share_ACTION_AttackAirHi_figatree", "PlyCaptain5K_Share_ACTION_RunBrake_figatree":
		// case "PlyCaptain5K_Share_ACTION_Attack100Loop_figatree":
		// case "":
		default:
			fmt.Printf("%d: ", index)
			if strings.Contains(string(action.Name.Value), "PlyCaptain5K") {
				fmt.Printf("%s\n", action.Name.Value[26:len(action.Name.Value)-9])
			} else if strings.Contains(string(action.Name.Value), "PlyTaro") {
				fmt.Printf("%s\n", action.Name.Value[21:len(action.Name.Value)-9])
			} else if len(action.Name.Value) == 0 {
				fmt.Println("NoName")
			} else {
				fmt.Printf("%s\n", action.Name.Value)
			}
			if !printSubactions {
				continue
			}
			for _, subac := range action.Subactions.Value {
				fmt.Printf("\t%s\n", fighter.SubActionString(subac))
			}
		}
	}
}

var animNameRegex = regexp.MustCompile(`.*_([a-zA-Z0-9]+)_figatree`)

func main() {
	fighterFile, err := os.Open("fighter.dat")
	if err != nil {
		log.Fatal(err)
	}
	defer fighterFile.Close()
	reader := binread.NewReader(fighterFile)

	fighter := FighterFile{}

	err = reader.Decode(&fighter)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("%#+v\n", h)
	// spew.Dump(fighter.Desc.Header)
	// spew.Dump(fighter.Desc.Relocation)
	// spew.Dump(fighter.Data.AttributesSpecial)
	// spew.Dump(fighter.Data.ECB)
	// spew.Dump(fighter.Data.JostleBox)
	// spew.Dump(fighter.Data.Hurtboxes)
	// printActions(fighter.Data.ActionTable.Value, true)

	// ANIMATION
	animationFile, err := os.Open("animation.dat")
	if err != nil {
		log.Fatal(err)
	}
	defer animationFile.Close()
	reader = binread.NewReader(animationFile)

	animationFilesAsBytes, err := SplitAnimationFile(reader)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range animationFilesAsBytes {
		reader = binread.NewReader(bytes.NewReader(file))
		animation := AnimationFile{}

		err = reader.Decode(&animation)
		if err != nil {
			log.Fatal(err)
		}

		// spew.Dump(animation.Desc.Header)
		spew.Dump(animation.Desc.Footer.Roots[0].Name)
		spew.Dump(animation.Data)
		break
	}

	fmt.Printf("Num anim files: %d\n", len(animationFilesAsBytes))
}
