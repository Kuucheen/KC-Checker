package helper

import (
	"fmt"
	"os"
	"regexp"
)

var ProxySum float64

func Write(proxies map[int][]*Proxy, style int, bancheck bool) string {
	pType := GetTypeName()

	for _, proxyLevel := range proxies {
		filtered := ""

		if bancheck {
			filtered = "BanChecked/"
		}

		f, err := os.Create(GetFilePath(pType) + filtered + GetLevelNameOf(proxyLevel[0].Level-1) + ".txt")
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

	ipPortPattern := `\b(?:\d{1,3}\.){3}\d{1,3}:\d+\b`
	re := regexp.MustCompile(ipPortPattern)

	matches := re.FindAllString(string(dat), -1)

	ProxySum = float64(len(matches))

	return matches
}

func GetFilePath(name string) string {
	return fmt.Sprintf("output/%s/", name)
}
