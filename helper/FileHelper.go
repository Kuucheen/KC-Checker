package helper

import (
	"KC-Checker/common"
	"fmt"
	"golang.design/x/clipboard"
	"os"
	"regexp"
	"sort"
)

var ProxySum float64

func Write(proxies map[int][]*Proxy, style int, banCheck bool) string {
	pType := GetTypeName()

	var allFile *os.File
	var allFileErr error

	if common.GetConfig().CopyToClipboard {
		cliperr := clipboard.Init()

		if cliperr != nil {
			return "ClipBoard error"
		}
	}

	clipString := ""

	for _, proxyLevel := range proxies {

		sort.Slice(proxyLevel, func(i, j int) bool {
			return proxyLevel[i].time < proxyLevel[j].time
		})

		filtered := ""

		if banCheck {
			filtered = "BanChecked/"
		}

		f, err := os.Create(GetFilePath(pType) + filtered + GetLevelNameOf(proxyLevel[0].Level-1) + ".txt")

		if allFile == nil {
			allFile, allFileErr = os.Create(GetFilePath(pType) + filtered + "all.txt")
		}

		if err != nil || allFileErr != nil {
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
			_, allFileErr = fmt.Fprintln(allFile, proxyString)

			clipString += proxyString + "\n"

			if err != nil || allFileErr != nil {
				return ""
			}
		}

		err = f.Close()

		if common.GetConfig().CopyToClipboard {
			clipboard.Write(clipboard.FmtText, []byte(clipString))
		}

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
