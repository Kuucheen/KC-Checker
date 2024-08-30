package helper

import (
	"KC-Checker/common"
	"sync"
	"sync/atomic"
	"time"
)

var (
	proxyQueue       = ProxyQueue{}
	ProxyMap         = make(map[int][]*Proxy)
	ProxyMapFiltered = make(map[int][]*Proxy) //ProxyMapFiltered is for Banchecked proxies
	ProxyCountMap    = make(map[int]int)
	stop             = false
	Invalid          int32
	threadsActive    int32
	mutex            sync.Mutex
	wg               sync.WaitGroup

	retries int

	HasFinished = false
)

type CPMCounter struct {
	mu          sync.Mutex
	checks      int
	lastUpdated time.Time
}

type ProxyListing struct {
	mu      sync.Mutex
	proxies []*Proxy
}

var cpmCounter = CPMCounter{}
var proxyList = ProxyListing{}

func Dispatcher(proxies []*Proxy) {
	threads := common.GetConfig().Threads
	retries = common.GetConfig().Retries
	proxyList.proxies = proxies

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go threadHandling()
		atomic.AddInt32(&threadsActive, 1)
	}

	wg.Wait()
	cpmCounter.mu.Lock()
	cpmCounter.checks = 0
	cpmCounter.mu.Unlock()
	HasFinished = true
}

func threadHandling() {
	for len(proxyList.proxies) > 0 {
		proxyList.mu.Lock()
		if len(proxyList.proxies) > 0 {
			proxy := proxyList.proxies[0]
			proxyList.proxies = proxyList.proxies[1:]
			proxyList.mu.Unlock()
			check(proxy)
		} else {
			proxyList.mu.Unlock()
		}
		if stop {
			break
		}
	}

	defer func() {
		atomic.AddInt32(&threadsActive, -1)
		wg.Done()
	}()
}

func check(proxy *Proxy) {
	responded := false
	level := 0

	cpmCounter.mu.Lock()
	cpmCounter.checks++
	now := time.Now()

	if now.Sub(cpmCounter.lastUpdated) >= time.Minute {
		cpmCounter.checks = 1
		cpmCounter.lastUpdated = now
	} else {
		cpmCounter.checks++
	}
	cpmCounter.mu.Unlock()

	for proxy.checks <= retries {

		timeStart := time.Now()
		body, status, err := Request(proxy)
		timeEnd := time.Now()

		proxy.time = int(timeEnd.Sub(timeStart).Milliseconds())

		if err != nil {
			status = -1
		}

		if status >= 400 || status == -1 {
			proxy.checks++
			continue
		}

		level = GetProxyLevel(body)
		mutex.Lock()

		proxy.Level = level
		ProxyMap[level] = append(ProxyMap[level], proxy)
		ProxyCountMap[level]++
		proxyQueue.Enqueue(proxy)

		mutex.Unlock()

		responded = true
		break
	}

	//Ban check for websites
	if responded {
		//Extra if because of performance (else)
		if common.DoBanCheck() {
			for i := 0; i < retries; i++ {
				body, status, err := RequestCustom(proxy, common.GetConfig().Bancheck)

				if err != nil {
					status = -1
				}

				if !(status >= 400) && status != -1 {
					keywords := common.GetConfig().Keywords

					if len(keywords) == 0 || len(keywords[0]) == 0 || ContainsSlice(keywords, body) {
						mutex.Lock()
						ProxyMapFiltered[level] = append(ProxyMapFiltered[level], proxy)
						ProxyCountMap[-1]++
						mutex.Unlock()
						break
					}
				}
			}
		}
	} else {
		atomic.AddInt32(&Invalid, 1)
	}
}

func GetInvalid() int {
	return int(atomic.LoadInt32(&Invalid))
}

func GetActive() int {
	return int(atomic.LoadInt32(&threadsActive))
}

func GetQueue() ProxyQueue {
	return proxyQueue
}

func StopThreads() {
	stop = true
}

func GetCPM() float64 {
	cpmCounter.mu.Lock()
	defer cpmCounter.mu.Unlock()

	now := time.Now()
	if now.Sub(cpmCounter.lastUpdated) >= time.Minute {
		defer func() { cpmCounter.checks = 0 }()
		cpmCounter.lastUpdated = now
	}

	return float64(cpmCounter.checks) / now.Sub(cpmCounter.lastUpdated).Minutes()
}
