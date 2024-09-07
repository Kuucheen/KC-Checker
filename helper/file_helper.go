package helper

import (
	"KC-Checker/common"
	"bufio"
	"fmt"
	"golang.design/x/clipboard"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"sync"
	"time"
)

var (
	ProxySum int

	muProxies      sync.Mutex
	proxiesToWrite = make(map[int][]*Proxy)

	muBancheck             sync.Mutex
	bancheckProxiesToWrite = make(map[int][]*Proxy)
)

func Write(proxies map[int][]*Proxy, style int, banCheck bool, appendToFile bool) string {
	pTypes := GetTypeNames()

	if common.GetConfig().CopyToClipboard {
		clipErr := clipboard.Init()
		if clipErr != nil {
			return "ClipBoard error"
		}
	}

	clipString := ""

	for _, pType := range pTypes {

		var allFile *os.File
		var allFileErr error

		for _, proxyLevel := range proxies {

			sort.Slice(proxyLevel, func(i, j int) bool {
				return proxyLevel[i].time < proxyLevel[j].time
			})

			filtered := ""
			if banCheck {
				filtered = "BanChecked/"
			}

			fileMode := os.O_CREATE | os.O_WRONLY
			if appendToFile {
				fileMode |= os.O_APPEND
			} else {
				fileMode |= os.O_TRUNC
			}

			f, err := os.OpenFile(GetFilePath(pType)+filtered+GetLevelNameOf(proxyLevel[0].Level-1)+".txt", fileMode, 0644)
			if allFile == nil {
				allFile, allFileErr = os.OpenFile(GetFilePath(pType)+filtered+"all.txt", fileMode, 0644)
			}

			if err != nil || allFileErr != nil {
				return ""
			}

			for _, proxy := range proxyLevel {
				if proxy.Protocol != pType {
					continue
				}

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
			if err != nil {
				return ""
			}

			if common.GetConfig().CopyToClipboard {
				clipboard.Write(clipboard.FmtText, []byte(clipString))
			}
		}
	}

	return "Wrote to output folder"
}

func GetFilePath(name string) string {
	return fmt.Sprintf("output/%s/", name)
}

func GetFullProxies(file string) []string {
	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("Error while reading proxies: %s", err)
		return nil
	}
	defer f.Close()

	fileInfo, err := f.Stat()
	if err != nil {
		fmt.Printf("Error while getting file info: %s", err)
		return nil
	}

	// Preallocate slice capacity based on an estimate
	var proxies []string
	if size := fileInfo.Size(); size > 0 {
		proxies = make([]string, 0, size/20) // Assuming average line length is 20 bytes
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		proxies = append(proxies, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error while scanning file: %s", err)
	}

	return proxies
}

// GetProxiesFile gets proxies/ips from a file with an option to filter full proxies.
func GetProxiesFile(file string, full bool) []string {
	if full {
		return GetFullProxies(file)
	}

	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("Error while reading proxies: %s", err)
		return nil
	}
	defer f.Close()

	var proxies []string
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		proxies = append(proxies, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error while scanning file: %s", err)
	}

	ProxySum = len(proxies)

	return proxies
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

	ProxySum = len(matches)

	return matches
}

func clearOutputFolder(path string) error {
	return filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing %s: %v", filePath, err)
		}

		// Skip the root folder itself
		if filePath == path {
			return nil
		}

		// If it's a directory, just continue
		if info.IsDir() {
			return nil
		}

		// Remove the file
		err = os.Remove(filePath)
		if err != nil {
			return fmt.Errorf("failed to remove file %s: %v", filePath, err)
		}

		return nil
	})
}

func StartAutoOutputManager() {
	_ = clearOutputFolder("./output")

	timeBetween := time.Duration(common.GetConfig().AutoOutput.TimeBetweenSafes) * time.Second

	for {
		if len(proxiesToWrite) > 0 {
			muProxies.Lock()
			Write(proxiesToWrite, common.GetAutoOutput(), false, true)
			proxiesToWrite = make(map[int][]*Proxy)
			muProxies.Unlock()
		}

		if len(bancheckProxiesToWrite) > 0 {
			muBancheck.Lock()
			Write(bancheckProxiesToWrite, common.GetAutoOutput(), true, true)
			bancheckProxiesToWrite = make(map[int][]*Proxy)
			muBancheck.Unlock()
		}

		time.Sleep(timeBetween)
	}
}

func AddToWriteQueue(proxy *Proxy) {
	muProxies.Lock()
	defer muProxies.Unlock()

	proxiesToWrite[proxy.Level] = append(proxiesToWrite[proxy.Level], proxy)
}

func AddToBancheckWriteQueue(proxy *Proxy) {
	muBancheck.Lock()
	defer muBancheck.Unlock()

	bancheckProxiesToWrite[proxy.Level] = append(bancheckProxiesToWrite[proxy.Level], proxy)
}
