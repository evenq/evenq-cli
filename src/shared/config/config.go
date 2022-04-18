package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFilePath = "/.evenq"
const configFileName = "config"

type data map[string]string

func GetValue(key string) (string, bool) {
	if file, ok := readFile(); ok {
		fileData := data{}
		if err := json.Unmarshal(file, &fileData); err == nil {
			if val, ok := fileData[key]; ok && val != "" {
				return val, true
			}
		}
	}

	return "", false
}

func SetValue(key string, value string) {
	fileData := data{}

	if file, ok := readFile(); ok {
		_ = json.Unmarshal(file, &fileData)
	}

	fileData[key] = value

	js, err := json.Marshal(fileData)
	if err != nil {
		fmt.Println(err)
		return
	}

	writeFile(js)
}

func readFile() ([]byte, bool) {
	path, ok := getFilePath()
	if !ok {
		return []byte{}, false
	}

	res, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, false
	}

	return res, true
}

func writeFile(data []byte) {
	path, ok := getFilePath()
	if !ok {
		return
	}

	basePath, ok := getBasePath()
	if !ok {
		return
	}

	err := os.MkdirAll(basePath, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}

	// ignore the error since we can go on without the file
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func getFilePath() (string, bool) {
	basePath, ok := getBasePath()
	if !ok {
		return "", false
	}

	return basePath + "/" + configFileName, true
}

func getBasePath() (string, bool) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return "", false
	}

	return dirname + configFilePath, true
}
