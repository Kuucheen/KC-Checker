package common

import (
	"io"
	"net/http"
	"sort"
	"time"
)

type JudgesTimes struct {
	Judge        string
	ResponseTime time.Duration
}

type HostTimes []JudgesTimes

var (
	UserIP       string
	FastestJudge string
)

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
	configHosts := GetConfig().Judges

	for _, value := range configHosts {
		responseTime := checkTime(value)
		hostTime := JudgesTimes{
			Judge:        value,
			ResponseTime: responseTime,
		}
		CurrentCheckedHosts = append(CurrentCheckedHosts, hostTime)
	}

	// Create a copy of the unsorted CurrentCheckedHosts
	unsortedHosts := make(HostTimes, len(CurrentCheckedHosts))
	copy(unsortedHosts, CurrentCheckedHosts)

	// Sort the original CurrentCheckedHosts based on response time
	sort.Sort(CurrentCheckedHosts)

	FastestJudge = CurrentCheckedHosts[0].Judge

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

		UserIP = string(respBody)

		return UserIP
	}

	panic("Couldnt get the Users IP please provide an other ip sources!")
}
