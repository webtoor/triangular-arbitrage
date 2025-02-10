package file

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v3"
)

// ReadFromYAML reads the YAML file and pass to the object
// args:
//
//	path: file path location
//	target: object which will hold the value
//
// returns:
//
//	error: operation state error
func ReadFromYAML(path string, target interface{}) error {
	yf, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(yf, target)
}

func ReadFromJSON(path string, target any) error {
	yf, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(yf, target)
}

func WriteToJson(path string, in any) error {
	c, err := os.Create(path)
	if err != nil {
		return err
	}

	defer c.Close()

	v, _ := json.Marshal(in)

	_, err = c.Write(v)

	if err != nil {
		return err
	}

	return nil
}
