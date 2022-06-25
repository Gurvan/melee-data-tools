package main

import (
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
	file, err := os.Open("file.dat")
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
	spew.Dump(h.Desc.Header)
	// spew.Dump(h.Desc.Relocation)
	// spew.Dump(h.Data.AttributesSpecial)
	// spew.Dump(h.Data.ECB)
	// spew.Dump(h.Data.JostleBox)
	// spew.Dump(h.Data.ActionTable.Value)
	// actions := h.Data.ActionTable.Value
	// for _, action := range actions {
	//         fmt.Println(action.AnimationOffset, action.Name, action.Subactions.Value)
	// }

	// fmt.Println(len(h.Data.ActionTable.Value))
	printActions(h.Data.ActionTable.Value, true)

	// ANIMATION
	// fileAnim, err := os.Open("file_anim.dat")
	// if err != nil {
	//         log.Fatal(err)
	// }

	// reader = binread.NewReader(fileAnim)
	// reader.Seek(actions[0].AnimationOffset.ToSeek(), io.SeekStart)

	// filesData, err := fighter.SplitAnimationFile(reader)
	// if err != nil {
	//         log.Fatal(err)
	// }
	// fmt.Println(len(filesData))
	// for _, data := range filesData {
	//         r := binread.NewReader(bytes.NewReader(data))
	//         d := descriptor.Descriptor{}
	//         err := r.Decode(&d)
	//         if err != nil {
	//                 log.Fatal(err)
	//         }
	//         // fmt.Println(d.Footer.Roots[0].Name, animNameRegex.FindStringSubmatch(string(d.Footer.Roots[0].Name))[1]) // .MatchString(string(d.Footer.Roots[0].Name)))
	//         fmt.Println(animNameRegex.FindStringSubmatch(string(d.Footer.Roots[0].Name))[1])
	// }
}
