package utils

import (
	"testing"
)

func AssertExpectedJSON(t *testing.T, obj interface{}, filename string) {
	expected := string(ReadTestFile(t, filename))

	res, err := ToJSON(obj)
	if err != nil {
		t.Errorf("%v", err)
	}

	AssertString(t, string(res), expected)
}

func AssertString(t *testing.T, res string, expected string) {
	if res != expected {
		t.Errorf("got: %s - expected: %s", res, expected)
	}
}
