package main

import (
	"log"

	"github.com/harshjoeyit/word-dict/dict"
	"github.com/harshjoeyit/word-dict/s3dict"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Uncomment the following lines to build a new dictionary
	//
	// err := dict.BuildNewDict()
	// if err != nil {
	// 	log.Fatalf("Error building new dictionary: %v", err)
	// }

	d1, err := dict.New()
	if err != nil {
		log.Fatalf("Error creating new dictionary: %v", err)
	}
	defer d1.Close()

	// Query the dictionary for a word
	def, ok := d1.QueryWord("abandon")
	if ok {
		log.Printf("Definition: %s", def)
	}

	// Uncomment the following lines to update the dictionary
	//
	// err = dict.UpdateDict()
	// if err != nil {
	// 	log.Fatalf("Error updating dictionary: %v", err)
	// }

	// d2, err := dict.New()
	// if err != nil {
	// 	panic(err)
	// }
	// defer d2.Close()

	// // Query the dictionary for a word
	// def, ok = d2.QueryWord("abandon")
	// if ok {
	// 	log.Printf("Definition: %s", def)
	// }

	s3d, err := s3dict.New()
	if err != nil {
		log.Fatalf("Error creating new S3 dictionary: %v", err)
	}

	// Query the dictionary for a word
	def, ok = s3d.QueryWord("tiger")
	if ok {
		log.Printf("Definition from S3: %s", def)
	}
}
