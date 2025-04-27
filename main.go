package main

import (
	"log"

	"github.com/harshjoeyit/word-dict/dict"
)

func main() {
	// err := dict.BuildNewDict()
	// if err != nil {
	// 	log.Fatalf("Error building new dictionary: %v", err)
	// }

	d1, err := dict.NewDict()
	if err != nil {
		log.Fatalf("Error creating new dictionary: %v", err)
	}
	defer d1.Close()

	// Query the dictionary for a word
	def, ok := d1.QueryWord("abandon")
	if ok {
		log.Printf("Definition: %s", def)
	}

	err = dict.UpdateDict()
	if err != nil {
		log.Fatalf("Error updating dictionary: %v", err)
	}

	d2, err := dict.NewDict()
	if err != nil {
		panic(err)
	}
	defer d2.Close()

	// Query the dictionary for a word
	def, ok = d2.QueryWord("abandon")
	if ok {
		log.Printf("Definition: %s", def)
	}
}
