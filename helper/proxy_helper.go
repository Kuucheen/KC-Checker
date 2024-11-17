package helper

import (
	"KC-Checker/common"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Proxy struct {
	Ip   string
	Port int
	Full string

	Level    int    // anonymity
	Protocol string // http, socks4...
	Country  string
	Type     string // ISP, Residential or Datacenter

	checks int
	Time   int //in ms

	// For authentication if existing
	username string
	password string
}

var (
	proxyType     []int
	Blacklisted   []string
	AllProxies    []*Proxy
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
			value = strings.Replace(value, "@", ":", 1)
			splitProxyAuth := strings.Split(value, ":")

			tempLength := len(splitProxyAuth)

			ip := splitProxyAuth[0]

			if (tempLength != 4 && tempLength != 2) || !checkIp(ip) {
				return
			}

			dat, err := strconv.Atoi(splitProxyAuth[1])

			if err != nil {
				return
			}

			newProxy := &Proxy{
				Ip:   ip,
				Port: dat,
				Full: splitProxyAuth[0] + ":" + splitProxyAuth[1],
			}

			if tempLength == 4 {
				newProxy.username = splitProxyAuth[2]
				newProxy.password = splitProxyAuth[3]
			}

			proxyChan <- newProxy
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

	if len(typeNames) == 0 {
		return arr
	}

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
	if AllProxiesSum > 0 {
		withAllProtocols := AddAllProtocols(AllProxies)

		AllProxiesSum = float64(len(withAllProtocols))

		return withAllProtocols
	}

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

	AllProxies = AddAllProtocols(cleaned)

	AllProxiesSum = float64(len(cleaned))

	if AllProxiesSum == 0 {
		AllProxiesSum = -1
	}

	return AllProxies
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
