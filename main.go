package main

import (
	"log"

	"github.com/harshjoeyit/word-dict/dict"
)

func main() {
	d, err := dict.NewDict()
	if err != nil {
		panic(err)
	}
	defer d.Close()

	// Query the dictionary for a word
	def, ok := d.QueryWord("abandon")
	if ok {
		log.Printf("Definition: %s", def)
	}
}
