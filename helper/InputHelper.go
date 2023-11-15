package helper

import (
	"fmt"
	"os"
	"strings"
)

func GetInput() []string {
	dat, err := os.ReadFile("proxies.txt")
	if err != nil {
		fmt.Printf("Error while reading proxies: %s", err)
	}

	//Split into lines
	return strings.Split(strings.ReplaceAll(string(dat), "\r\n", "\n"), "\n")
}
