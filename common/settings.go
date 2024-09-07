package common

import (
	"KC-Checker/charm/errorDisplays"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Threads         int        `json:"threads"`
	Retries         int        `json:"retries"`
	Timeout         int        `json:"timeout"`
	IpLookup        string     `json:"iplookup"`
	JudgesThreads   int        `json:"judges_threads"`
	JudgesTimeOut   int        `json:"judges_timeout"`
	Judges          []string   `json:"judges"`
	Blacklisted     []string   `json:"blacklisted"`
	Bancheck        string     `json:"bancheck"`
	Keywords        []string   `json:"keywords"`
	PrivacyMode     bool       `json:"privacy_mode"`
	CopyToClipboard bool       `json:"copyToClipboard"`
	AutoSelect      autoSelect `json:"autoSelect"`
	Transport       transport  `json:"transport"`
}

type autoSelect struct {
	Http   bool
	Https  bool
	Socks4 bool
	Socks5 bool
}

type transport struct {
	KeepAlive             bool `json:"KeepAlive"`
	KeepAliveSeconds      int  `json:"KeepAliveSeconds"`
	MaxIdleConns          int  `json:"MaxIdleConns"`
	MaxIdleConnsPerHost   int  `json:"MaxIdleConnsPerHost"`
	IdleConnTimeout       int  `json:"IdleConnTimeout"`
	TLSHandshakeTimeout   int  `json:"TLSHandshakeTimeout"`
	ExpectContinueTimeout int  `json:"ExpectContinueTimeout"`
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
		errorDisplays.PrintErrorForSettings(err)
		return
	}
}

func RemoveHttpJudges() {
	removeJudge("https://")
}

func RemoveHttpsJudges() {
	removeJudge("http://")
}

func removeJudge(str string) {
	var httpsJudges []string

	for _, i2 := range config.Judges {
		if strings.HasPrefix(i2, str) {
			httpsJudges = append(httpsJudges, i2)
		}
	}

	config.Judges = httpsJudges
}

func GetConfig() Config {
	return config
}

func GetPrivacyMode() bool { return config.PrivacyMode }

func DoBanCheck() bool {
	return config.Bancheck != ""
}

func IsAllowedToCheck(typeNames []string) bool {

	for _, name := range typeNames {
		hasBeenFound := false

		for _, judge := range config.Judges {
			if strings.HasPrefix(judge, name) {
				hasBeenFound = true
				break
			}
		}

		if !hasBeenFound {
			return false
		}
	}

	return true
}