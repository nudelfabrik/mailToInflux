package settings

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type Settings struct {
	URL    string `json:"url"`
	Token  string `json:"token"`
	Bucket string `json:"bucket"`
	Org    string `json:"org"`
}

func LoadSettings(altPath string) (*Settings, error) {
	var file []byte

	var err error

	paths := []string{"/usr/local/etc/mti.json", "./mti.json"}

	if altPath != "" {
		// Use the explicitly specified path
		file, err = ioutil.ReadFile(altPath)
		if err != nil {
			return nil, err
		}
	} else {
		// Try all default paths
		for _, path := range paths {
			file, err = ioutil.ReadFile(path)
			if err == nil {
				break
			}
		}
	}

	if file == nil {
		return nil, errors.New("No File found")
	}

	var setting Settings
	err = json.Unmarshal(file, &setting)

	return &setting, err
}
