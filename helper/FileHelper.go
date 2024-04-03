package helper

import (
	"fmt"
	"os"
	"regexp"
	"sort"
)

var ProxySum float64

func Write(proxies map[int][]*Proxy, style int, banCheck bool) string {
	pType := GetTypeName()

	for _, proxyLevel := range proxies {

		sort.Slice(proxyLevel, func(i, j int) bool {
			return proxyLevel[i].time < proxyLevel[j].time
		})

		filtered := ""

		if banCheck {
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

			case 2:
				proxyString = fmt.Sprintf("%s;%d", proxy.Full, proxy.time)
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

func GetFilePath(name string) string {
	return fmt.Sprintf("output/%s/", name)
}

// GetProxiesFile gets proxies/ips from a file
func GetProxiesFile(file string, full bool) []string {
	dat, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Error while reading proxies: %s", err)
	}

	return GetProxies(string(dat), full)
}

func GetProxies(str string, full bool) []string {
	ipPortPattern := ""

	if full {
		ipPortPattern = `\b(?:\d{1,3}\.){3}\d{1,3}:\d+\b`
	} else {
		ipPortPattern = `\b(?:\d{1,3}\.){3}\d{1,3}\b`
	}

	re := regexp.MustCompile(ipPortPattern)

	matches := re.FindAllString(str, -1)

	ProxySum = float64(len(matches))

	return matches
}
