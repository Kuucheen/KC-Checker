package common

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Threads     int      `json:"threads"`
	Retries     int      `json:"retries"`
	Timeout     int      `json:"timeout"`
	IpLookup    string   `json:"iplookup"`
	Judges      []string `json:"judges"`
	Blacklisted []string `json:"blacklisted"`
	Bancheck    string   `json:"bancheck"`
	Keywords    []string `json:"keywords"`
}

var config Config

func ReadSettings() {
	data, err := os.ReadFile("settings.json")
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

func GetConfig() Config {
	return config
}

func DoBanCheck() bool {
	return config.Bancheck != ""
}
