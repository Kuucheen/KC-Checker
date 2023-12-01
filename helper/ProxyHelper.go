package helper

import (
	"fmt"
	"strconv"
	"strings"
)

type proxy struct {
	ip     string
	port   int
	full   string
	checks int
}

var proxyType int

func ToProxies(arr []string) []*proxy {
	var newArr []*proxy
	for _, value := range arr {
		temp := strings.Split(value, ":")

		if !checkIp(temp[0]) {
			continue
		}

		dat, err := strconv.Atoi(temp[1])

		if err != nil {
			fmt.Print("Not a valid port: ", err)
		}

		newArr = append(newArr, &proxy{
			ip:   temp[0],
			port: dat,
			full: temp[0] + ":" + temp[1],
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

func GetType() int {
	return proxyType
}

func GetTypeName() string {
	names := []string{"http", "socks4", "socks5"}
	return names[proxyType]
}

func SetType(typ int) {
	proxyType = typ
}
