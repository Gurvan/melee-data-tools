package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/Gurvan/melee-data-tools/binread"
	"github.com/davecgh/go-spew/spew"
)

func PrettyString(obj interface{}) string {
	bytes, _ := json.MarshalIndent(obj, "  ", "  ")
	return string(bytes)
}

func main() {
	file, err := os.Open("file.dat")
	if err != nil {
		log.Fatal(err)
	}

	reader := binread.NewReader(file)

	// h := Header{}
	h := Descriptor{}

	err = reader.Decode(&h)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("%#+v\n", h)
	spew.Dump(h)
}
