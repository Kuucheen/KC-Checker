package helper

import (
	"fmt"
	"os"
	"strings"
)

var ProxySum float64

func Write(file string, data []byte) {
	err := os.WriteFile(file, data, 0644)
	if err != nil {
		fmt.Printf("Error while writing %s: %s", file, err)
	}
}

func GetInput(file string) []string {
	dat, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Error while reading proxies: %s", err)
	}
	split := strings.Split(strings.ReplaceAll(string(dat), "\r\n", "\n"), "\n")

	ProxySum = float64(len(split))

	return split
}
