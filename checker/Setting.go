package checker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Retries int      `json:"retries"`
	Timeout int      `json:"timeout"`
	Hosts   []string `json:"hosts"`
}

var config Config

func ReadSettings() {
	data, err := ioutil.ReadFile("settings.json")
	if err != nil {
		fmt.Println("Error reading settings file:", err)
		return
	}

	// Unmarshal the JSON data into the Config struct
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

}

func getConfig() Config {
	return config
}
