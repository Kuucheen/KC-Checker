package helper

import (
	"KC-Checker/charm"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type proxy struct {
	ip   string
	port int
	full string
}

type proxyType struct {
	pType int
}

func ToProxies(arr []string) []proxy {
	newArr := []proxy{}
	for _, value := range arr {
		temp := strings.Split(value, ":")

		if !checkIp(temp[0]) {
			continue
		}

		dat, err := strconv.Atoi(temp[1])

		if err != nil {
			fmt.Print("Not a valid port: ", err)
		}

		newArr = append(newArr, proxy{
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

func (m *proxyType) GetType() int {
	if reflect.TypeOf(m.pType) == nil {
		charm.RunType()
	}
	return m.pType
}

func (m *proxyType) SetType(value int) {
	m.pType = value
}
