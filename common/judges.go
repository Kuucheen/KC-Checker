package common

import (
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type JudgesTimes struct {
	Judge        string
	Ip           string
	Regex        string
	ResponseTime time.Duration
}

type HostTimes []JudgesTimes

var (
	UserIP             string
	FastestJudge       string
	FastestJudgeName   *url.URL
	FastestJudgeRegex  string
	FastestJudges      map[string]string
	FastestJudgesName  map[string]*url.URL
	FastestJudgesRegex map[string]string

	standardHeader = []string{"USER-AGENT", "HOST", "ACCEPT", "ACCEPT-ENCODING", "ACCEPT-ENCODING"}
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
	FastestJudges = make(map[string]string)
	FastestJudgesName = make(map[string]*url.URL)
	FastestJudgesRegex = make(map[string]string)

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

	FastestJudge = CurrentCheckedHosts[0].Ip
	FastestJudgeRegex = CurrentCheckedHosts[0].Regex

	u, err := url.Parse(CurrentCheckedHosts[0].Judge)
	if err == nil {
		FastestJudgeName = u
	}

	for _, host := range CurrentCheckedHosts {

		protocol := strings.Split(host.Judge, "://")[0]

		_, ok := FastestJudges[protocol]

		if !ok {
			FastestJudges[protocol] = host.Ip
			FastestJudgesRegex[protocol] = host.Regex
		}

		_, ok = FastestJudgesName[protocol]

		if !ok {
			u, err = url.Parse(host.Judge)
			if err == nil {
				FastestJudgesName[protocol] = u
			}
		}
	}

	// Return the unsorted CurrentCheckedHosts
	return unsortedHosts
}

func checkTimeAsync(host configJudge) {
	defer wg.Done()
	defer atomic.AddInt32(&currentThreads, -1)

	ip, responseTime := checkTime(host.Url, host.Regex)
	hostTime := JudgesTimes{
		Judge:        host.Url,
		Ip:           ip,
		Regex:        host.Regex,
		ResponseTime: responseTime,
	}

	mutex.Lock()
	CurrentCheckedHosts = append(CurrentCheckedHosts, hostTime)
	mutex.Unlock()
}

// Main function to check the time
func checkTime(host string, regex string) (string, time.Duration) {
	// Parse the URL to extract the hostname
	parsedURL, err := url.Parse(host)
	if err != nil {
		return "", time.Hour * 999
	}

	hostname := parsedURL.Hostname()

	tempTransport := &http.Transport{
		DisableKeepAlives: !GetConfig().Transport.KeepAlive,
		MaxIdleConns:      3,
		IdleConnTimeout:   time.Duration(GetConfig().JudgesTimeOut) * time.Millisecond,
	}

	client := &http.Client{
		Transport: tempTransport,
		Timeout:   time.Millisecond * time.Duration(config.JudgesTimeOut),
	}

	ips, err := net.LookupIP(hostname)
	if err != nil || len(ips) == 0 {
		return "", time.Hour * 999
	}

	ip := ips[0].String()

	startTime := time.Now()

	resp, err := client.Get(host)
	if err != nil {
		return ip, time.Hour * 999
	}
	defer resp.Body.Close()

	// Read the response body
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ip, time.Hour * 999
	}

	if !CheckForValidResponse(string(resBody), regex) {
		return ip, time.Hour * 99
	}

	return ip, time.Since(startTime)
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

	panic("Couldn't get the Users IP please provide an other ip sources!")
}

func GetFastestJudgeForProtocol(protocol string) string {
	if strings.HasPrefix(protocol, "socks") {
		return FastestJudge
	}

	return FastestJudges[protocol]
}

func GetFastestJudgeNameForProtocol(protocol string) *url.URL {
	if strings.HasPrefix(protocol, "socks") {
		return FastestJudgeName
	}

	return FastestJudgesName[protocol]
}

func GetFastestJudgeRegexForProtocol(protocol string) string {
	if strings.HasPrefix(protocol, "socks") {
		return FastestJudgeRegex
	}

	return FastestJudgesRegex[protocol]
}

func CheckForValidResponse(html string, regex string) bool {
	if strings.EqualFold(regex, "default") {
		html = strings.ReplaceAll(html, "_", "-")
		html = strings.ToUpper(html)

		for _, header := range standardHeader {
			if !strings.Contains(html, header) {

				return false
			}
		}

		return true
	}

	re, err := regexp.Compile(regex)
	if err != nil {
		return false
	}

	return re.MatchString(html)
}
