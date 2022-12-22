package utils

import (
	"io/ioutil"
	"testing"
)

func ReadTestFile(t *testing.T, inputFile string) []byte {
	input, err := ioutil.ReadFile(inputFile)
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	return input
}
