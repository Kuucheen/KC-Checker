package common

import (
	"io"
	"net/http"
	"sort"
	"time"
)

type HostTime struct {
	Host         string
	ResponseTime time.Duration
}

type HostTimes []HostTime

var UserIP string

func (ht HostTimes) Len() int {
	return len(ht)
}

func (ht HostTimes) Less(i, j int) bool {
	return ht[i].ResponseTime < ht[j].ResponseTime
}

func (ht HostTimes) Swap(i, j int) {
	ht[i], ht[j] = ht[j], ht[i]
}

var hosts HostTimes

func CheckDomains() HostTimes {
	configHosts := GetConfig().Hosts

	for _, value := range configHosts {
		responseTime := checkTime(value)
		hostTime := HostTime{
			Host:         value,
			ResponseTime: responseTime,
		}
		hosts = append(hosts, hostTime)
	}

	// Create a copy of the unsorted hosts
	unsortedHosts := make(HostTimes, len(hosts))
	copy(unsortedHosts, hosts)

	// Sort the original hosts based on response time
	sort.Sort(hosts)

	// Return the unsorted hosts
	return unsortedHosts
}

func GetHosts() HostTimes {
	return hosts
}

func checkTime(host string) time.Duration {
	startTime := time.Now()

	_, err := http.Get(host)
	if err != nil {
		return time.Nanosecond * 999999999
	}

	return time.Since(startTime)
}

// GetLocalIP gets the outgoing ip address when a packet is sent
func GetLocalIP() string {
	resp, err := http.Get(GetConfig().IpLookup)
	if err != nil {
		panic("Couldnt get the users ip!")
	}

	respBody, err := io.ReadAll(resp.Body)
	err = resp.Body.Close()
	if err != nil {
		return ""
	}
	if err != nil {
		panic(err)
	}

	UserIP = string(respBody)

	return UserIP
}
