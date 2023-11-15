package helper

import (
	"fmt"
	"os"
)

func Write(file string, data []byte) {
	err := os.WriteFile(file, data, 0644)
	if err != nil {
		fmt.Printf("Error while writing %s: %s", file, err)
	}
}
