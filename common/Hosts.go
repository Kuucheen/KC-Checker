package common

import (
	"io"
	"net/http"
	"sort"
	"sync"
	"sync/atomic"
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

var (
	CurrentCheckedHosts HostTimes
	wg                  sync.WaitGroup
	mutex               sync.Mutex
	currentThreads      int32
)

func CheckDomains() HostTimes {
	configHosts := GetConfig().Judges
	maxThreads := GetConfig().JudgesThreads

	for i := 0; i < len(configHosts); i++ {
		if int(currentThreads) < maxThreads {
			wg.Add(1)
			go checkTimeAsync(configHosts[i])
			atomic.AddInt32(&currentThreads, 1)
		} else {
			i--
		}
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Create a copy of the unsorted CurrentCheckedHosts
	unsortedHosts := make(HostTimes, len(CurrentCheckedHosts))
	copy(unsortedHosts, CurrentCheckedHosts)

	// Sort the original CurrentCheckedHosts based on response time
	sort.Sort(CurrentCheckedHosts)

	FastestJudge = CurrentCheckedHosts[0].Judge

	// Return the unsorted CurrentCheckedHosts
	return unsortedHosts
}

func checkTimeAsync(host string) {
	defer wg.Done()
	defer atomic.AddInt32(&currentThreads, -1)

	responseTime := checkTime(host)
	hostTime := JudgesTimes{
		Judge:        host,
		ResponseTime: responseTime,
	}

	mutex.Lock()
	CurrentCheckedHosts = append(CurrentCheckedHosts, hostTime)
	mutex.Unlock()
}

func checkTime(host string) time.Duration {
	client := http.Client{
		Timeout: time.Millisecond * time.Duration(config.JudgesTimeOut),
	}

	startTime := time.Now()

	_, err := client.Get(host)

	if err != nil {
		return time.Hour * 999
	}

	return time.Since(startTime)
}

func GetHosts() HostTimes {
	return CurrentCheckedHosts
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
