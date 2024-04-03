package helper

import (
	"KC-Checker/common"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Proxy struct {
	Ip     string
	Port   int
	Full   string
	Level  int
	checks int
	time   int //in ms
}

var (
	proxyType   int
	Blacklisted []string
)

// ToProxies converts a String array of proxies to proxy types
func ToProxies(arr []string) []*Proxy {
	var newArr []*Proxy
	for _, value := range arr {
		temp := strings.Split(value, ":")

		if !checkIp(temp[0]) {
			continue
		}

		dat, err := strconv.Atoi(temp[1])

		if err != nil {
			fmt.Print("Not a valid Port: ", err)
		}

		newArr = append(newArr, &Proxy{
			Ip:   temp[0],
			Port: dat,
			Full: temp[0] + ":" + temp[1],
		})
	}

	return newArr
}

func checkIp(ip string) bool {
	temp := strings.Split(ip, ".")

	for _, value := range temp {
		value, err := strconv.Atoi(value)

		if err != nil || (value > 255 || value < 0) {
			return false
		}
	}

	return true
}

func GetTypeName() string {
	names := []string{"http", "socks4", "socks5"}
	return names[proxyType]
}

func GetLevelNameOf(typ int) string {
	names := []string{"transparent", "anonymous", "elite"}
	return names[typ]
}

func SetType(typ int) {
	proxyType = typ
}

func GetCleanedProxies() []*Proxy {
	forbidden := GetProxiesFile("blacklisted.txt", false)
	for _, val := range Blacklisted {
		forbidden = append(forbidden, val)
	}

	normal := ToProxies(GetProxiesFile("proxies.txt", true))

	var cleaned []*Proxy

	for _, value := range normal {
		if !hasString(forbidden, value.Ip) {
			cleaned = append(cleaned, value)
		}
	}

	return cleaned
}

func hasString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
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
		resp, _ := http.Get(site)

		respBody, err := io.ReadAll(resp.Body)
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
