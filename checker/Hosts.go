package checker

import (
	"fmt"
	"log"
	"net"
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

// GetLocalIP gets the outgoing ip address when a packet is sent
func GetLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	UserIP = conn.LocalAddr().(*net.UDPAddr).IP.String()

	return UserIP
}
