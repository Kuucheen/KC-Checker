package helper

import (
	"fmt"
	"net/http"
	"sort"
	"time"
)

type HostTime struct {
	Host         string
	ResponseTime time.Duration
}

type HostTimes []HostTime

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

func CheckDomains() {
	configHosts := getConfig().Hosts

	for _, value := range configHosts {
		responseTime := checkTime(value)
		hostTime := HostTime{
			Host:         value,
			ResponseTime: responseTime,
		}
		hosts = append(hosts, hostTime)
	}

	// Sort the hosts based on response time
	sort.Sort(hosts)

	// Print the sorted list
	for _, host := range hosts {
		fmt.Printf("%s: %s\n", host.Host, host.ResponseTime)
		//TODO charm implementation
	}
}

func checkTime(host string) time.Duration {
	starttime := time.Now()

	_, err := http.Get(host)
	if err != nil {
		return time.Duration(0)
	}

	return time.Since(starttime)
}
