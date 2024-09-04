package helper

import (
	"KC-Checker/common"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Proxy struct {
	Ip       string
	Port     int
	Full     string
	Level    int
	Protocol string
	checks   int
	time     int //in ms
}

var (
	proxyType     []int
	Blacklisted   []string
	AllProxiesSum float64
)

// ToProxies converts a String array of proxies to proxy types
func ToProxies(arr []string) []*Proxy {
	var wg sync.WaitGroup
	proxyChan := make(chan *Proxy, len(arr))

	// Process each proxy line concurrently
	for _, value := range arr {
		wg.Add(1)
		go func(value string) {
			defer wg.Done()
			temp := strings.Split(value, ":")

			if len(temp) != 2 || !checkIp(temp[0]) {
				return
			}

			dat, err := strconv.Atoi(temp[1])

			if err != nil {
				fmt.Printf("Not a valid Port: %v\n", err)
				return
			}

			proxyChan <- &Proxy{
				Ip:   temp[0],
				Port: dat,
				Full: value,
			}
		}(value)
	}

	// Close the channel when all processing is done
	go func() {
		wg.Wait()
		close(proxyChan)
	}()

	// Collect results from the channel
	var newArr []*Proxy
	for proxy := range proxyChan {
		if proxy != nil {
			newArr = append(newArr, proxy)
		}
	}

	return newArr
}

// AddAllProtocols adds for every Protocol selected a proxy with the Protocol
func AddAllProtocols(arr []*Proxy) []*Proxy {
	typeNames := GetTypeNames()

	var newArr []*Proxy

	for _, protocol := range typeNames {
		for _, proxy := range arr {
			newProxy := *proxy
			newProxy.Protocol = protocol
			newArr = append(newArr, &newProxy)
		}
	}

	return newArr
}

func checkIp(ip string) bool {
	return net.ParseIP(ip) != nil && net.ParseIP(ip).To4() != nil
}

func GetTypeNames() []string {
	names := []string{"http", "https", "socks4", "socks5"}

	var selected = []string{}

	for _, i := range proxyType {
		selected = append(selected, names[i])
	}

	return selected
}

func ContainsTypeName(str string) bool {
	for _, s := range GetTypeNames() {
		if s == str {
			return true
		}
	}

	return false
}

func GetLevelNameOf(typ int) string {
	names := []string{"transparent", "anonymous", "elite"}
	return names[typ]
}

func SetType(typ []int) {
	proxyType = typ
}

func GetCleanedProxies() []*Proxy {
	forbidden := GetProxiesFile("blacklisted.txt", false)
	for _, val := range Blacklisted {
		forbidden = append(forbidden, val)
	}

	forbiddenSet := make(map[string]struct{}, len(forbidden))
	for _, ip := range forbidden {
		forbiddenSet[ip] = struct{}{}
	}

	normal := ToProxies(GetProxiesFile("proxies.txt", true))

	var cleaned []*Proxy

	for _, value := range normal {
		if _, found := forbiddenSet[value.Ip]; !found {
			cleaned = append(cleaned, value)
		}
	}

	cleaned = AddAllProtocols(cleaned)

	AllProxiesSum = float64(len(cleaned))

	return cleaned
}

func ContainsSlice(slice []string, str string) bool {
	for _, s := range slice {
		if strings.Contains(str, s) {
			return true
		}
	}
	return false
}

func GetBlacklisted() []string {
	var blist []string

	for _, site := range common.GetConfig().Blacklisted {
		resp, err := http.Get(site)
		if err != nil {
			continue
		}

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}
		err = resp.Body.Close()

		if err != nil {
			continue
		}

		for _, value := range GetProxies(string(respBody), false) {
			blist = append(blist, value)
		}
	}

	Blacklisted = blist

	return Blacklisted
}
