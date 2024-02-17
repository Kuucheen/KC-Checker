package helper

import (
	"KC-Checker/common"
	"sync"
	"sync/atomic"
	"time"
)

const (
	DelayBetweenChecks = time.Millisecond * 10
)

var (
	proxyQueue       = ProxyQueue{}
	ProxyMap         = make(map[int][]*Proxy)
	ProxyMapFiltered = make(map[int][]*Proxy)
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

var cpmCounter = CPMCounter{}

func Dispatcher(proxies []*Proxy) {
	threads := common.GetConfig().Threads
	retries = common.GetConfig().Retries

	for len(proxies) > 0 {
		if int(atomic.LoadInt32(&threadsActive)) < threads {
			wg.Add(1)
			go check(proxies[0])
			atomic.AddInt32(&threadsActive, 1)
			proxies = proxies[1:]
		} else {
			time.Sleep(DelayBetweenChecks)
		}
		if stop {
			break
		}
	}

	wg.Wait()
	cpmCounter.mu.Lock()
	cpmCounter.checks = 0
	cpmCounter.mu.Unlock()
	HasFinished = true
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
		body, status := Request(proxy)

		if status >= 400 || status == -1 {
			proxy.checks++
			continue
		}

		level = GetProxyLevel(body)
		levels := []int{1, 2, 3}

		if isInList(level, levels) {
			mutex.Lock()

			proxy.Level = level
			ProxyMap[level] = append(ProxyMap[level], proxy)
			ProxyCountMap[level]++
			proxyQueue.Enqueue(proxy)

			mutex.Unlock()
		} else {
			atomic.AddInt32(&Invalid, 1)
		}

		responded = true
		break
	}

	//Ban check for websites
	if responded && common.DoBanCheck() {
		for i := 0; i < retries; i++ {
			body, status := RequestCustom(proxy, common.GetConfig().Bancheck)

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
	} else {
		atomic.AddInt32(&Invalid, 1)
	}

	defer func() {
		atomic.AddInt32(&threadsActive, -1)
		wg.Done()
	}()
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

func isInList(target int, list []int) bool {
	for _, value := range list {
		if value == target {
			return true
		}
	}
	return false
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
