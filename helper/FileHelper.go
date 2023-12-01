package helper

import (
	"fmt"
	"os"
	"strings"
)

func Write(file string, data []byte) {
	err := os.WriteFile(file, data, 0644)
	if err != nil {
		fmt.Printf("Error while writing %s: %s", file, err)
	}
}

func GetInput() []string {
	dat, err := os.ReadFile("proxies.txt")
	if err != nil {
		fmt.Printf("Error while reading proxies: %s", err)
	}

	//Split into lines
	return strings.Split(strings.ReplaceAll(string(dat), "\r\n", "\n"), "\n")
}
