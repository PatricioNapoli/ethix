package utils

import (
	"io/ioutil"
	"log"
)

func ReadFile(inputFile string) []byte {
	input, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Printf("could not open file. details: %v", err)
	}

	return input
}
