package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/Gurvan/melee-data-tools/binread"
	"github.com/Gurvan/melee-data-tools/fighter"
	"github.com/davecgh/go-spew/spew"
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
	// spew.Dump(h.Data.ActionTable.Value)
	actions := h.Data.ActionTable.Value
	for _, action := range actions {
		fmt.Println(action.AnimationOffset)
	}

	// fmt.Println(len(h.Data.ActionTable.Value))
	// printSubActions(h.Data.ActionTable.Value)

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
