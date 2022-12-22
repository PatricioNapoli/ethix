package utils

import "encoding/json"

func ToJSON(v interface{}) ([]byte, error) {
	j, err := json.Marshal(v)
	return j, err
}

func ToPrettyJSON(v interface{}) ([]byte, error) {
	j, err := json.MarshalIndent(v, "", "    ")
	return j, err
}

func FromJSON(d []byte, v interface{}) error {
	err := json.Unmarshal(d, v)
	return err
}
