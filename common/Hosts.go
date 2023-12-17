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

var CurrentCheckedHosts HostTimes

func CheckDomains() HostTimes {
	configHosts := GetConfig().Hosts

	for _, value := range configHosts {
		responseTime := checkTime(value)
		hostTime := HostTime{
			Host:         value,
			ResponseTime: responseTime,
		}
		CurrentCheckedHosts = append(CurrentCheckedHosts, hostTime)
	}

	// Create a copy of the unsorted CurrentCheckedHosts
	unsortedHosts := make(HostTimes, len(CurrentCheckedHosts))
	copy(unsortedHosts, CurrentCheckedHosts)

	// Sort the original CurrentCheckedHosts based on response time
	sort.Sort(CurrentCheckedHosts)

	// Return the unsorted CurrentCheckedHosts
	return unsortedHosts
}

func GetHosts() HostTimes {
	return CurrentCheckedHosts
}

func checkTime(host string) time.Duration {
	startTime := time.Now()

	_, err := http.Get(host)
	if err != nil {
		return time.Hour * 999
	}

	return time.Since(startTime)
}

func GetLocalIP() string {
	for i := 0; i < 2; i++ {
		resp, err := http.Get(GetConfig().IpLookup)
		if err != nil {
			continue
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

	panic("Couldnt get the Users IP please provide an other ip sources!")
}
