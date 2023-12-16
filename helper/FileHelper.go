package helper

import (
	"fmt"
	"os"
	"strings"
)

var ProxySum float64

func Write(proxies map[int][]*Proxy, style int) string {
	pType := GetTypeName()

	for _, proxyLevel := range proxies {
		f, err := os.Create(GetFilePath(pType) + GetLevelNameOf(proxyLevel[0].Level-1) + ".txt")
		if err != nil {
			return ""
		}
		for _, proxy := range proxyLevel {
			var proxyString string
			switch style {
			case 0:
				proxyString = proxy.Full
			case 1:
				proxyString = fmt.Sprintf("%s://%s", pType, proxy.Full)
			}

			_, err := fmt.Fprintln(f, proxyString)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return ""
			}
		}
		err = f.Close()
		if err != nil {
			return ""
		}
	}

	return "Wrote to " + GetFilePath(pType)
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

func GetFilePath(name string) string {
	return fmt.Sprintf("output/%s/", name)
}
