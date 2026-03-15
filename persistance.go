package main

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type yamlData struct {
	Disas    string   `yaml:"disas,omitempty"`
	Comments []string `yaml:"comments,omitempty"`
}

func loadFile(path string) (*yamlData, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var yamlData yamlData
	err = yaml.Unmarshal(content, &yamlData)
	if err != nil {
		return nil, err
	}

	return &yamlData, nil
}

func saveFile(path string, yamlData *yamlData) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = yaml.NewEncoder(file).Encode(yamlData)
	return err
}
