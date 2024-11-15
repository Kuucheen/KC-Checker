package helper

import (
	"KC-Checker/common"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

var (
	proxyQueue            = ProxyQueue{}
	ProxyMap              = make(map[int][]*Proxy)
	ProxyMapFiltered      = make(map[int][]*Proxy) //ProxyMapFiltered is for Banchecked proxies
	ProxyCountMap         = make(map[int]int)
	ProxyProtocolCountMap = make(map[string]map[int]int)
	stop                  = false
	Invalid               int32
	threadsActive         int32
	mutex                 sync.Mutex
	wg                    sync.WaitGroup

	retries int

	HasFinished = false
)

type ProxyListing struct {
	mu      sync.Mutex
	proxies []*Proxy
	length  int64
	index   atomic.Int64
}

var proxyList = ProxyListing{}

func Init() {
	InitializeCPM()

	ProxyProtocolCountMap["http"] = make(map[int]int)
	ProxyProtocolCountMap["https"] = make(map[int]int)
	ProxyProtocolCountMap["socks4"] = make(map[int]int)
	ProxyProtocolCountMap["socks5"] = make(map[int]int)
}

func Dispatcher(proxies []*Proxy) {
	Init()
	threads := common.GetConfig().Threads
	retries = common.GetConfig().Retries
	proxyList.proxies = proxies
	proxyList.length = int64(len(proxies))
	proxyList.index.Store(0)

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go threadHandling()
		atomic.AddInt32(&threadsActive, 1)
		time.Sleep(time.Millisecond * 1) //This is to prevent adding unnecessary httpClients to the pool
	}

	wg.Wait()
	EndCPMCounter()
	HasFinished = true
}

func threadHandling() {
	for proxyList.index.Load() < proxyList.length {
		proxyList.mu.Lock()
		if proxyList.index.Load() < proxyList.length {
			proxy := proxyList.proxies[proxyList.index.Load()]
			proxyList.index.Add(1)
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

	for proxy.checks <= retries {
		timeStart := time.Now()
		body, status, err := Request(proxy)
		timeEnd := time.Now()
		IncrementCheckCount()

		proxy.Time = int(timeEnd.Sub(timeStart).Milliseconds())

		if err != nil {
			status = -1
		}

		if status >= 400 || status == -1 {
			proxy.checks++
			continue
		}

		level = GetProxyLevel(body)
		proxy.Level = level

		mutex.Lock()
		ProxyMap[level] = append(ProxyMap[level], proxy)
		ProxyCountMap[level]++
		ProxyProtocolCountMap[proxy.Protocol][level]++
		proxyQueue.Enqueue(proxy)
		mutex.Unlock()

		AddToWriteQueue(proxy)
		proxy.Country = GetCountryCode(proxy.Ip)
		proxy.Type = DetermineProxyType(proxy.Ip)

		responded = true
		break
	}

	u, err := url.Parse(common.GetConfig().Bancheck)
	//Ban check for websites
	if responded && common.DoBanCheck() && err == nil {

		for i := 0; i < retries; i++ {
			body, status, err := RequestCustom(proxy, common.GetConfig().Bancheck, u, "", true)
			IncrementCheckCount()

			if err != nil {
				status = -1
			}

			if !(status >= 400) && status != -1 {
				keywords := common.GetConfig().Keywords

				if len(keywords) == 0 || len(keywords[0]) == 0 || ContainsSlice(keywords, body) {
					mutex.Lock()
					ProxyMapFiltered[level] = append(ProxyMapFiltered[level], proxy)
					ProxyCountMap[-1]++
					ProxyProtocolCountMap[proxy.Protocol][-1]++
					AddToBancheckWriteQueue(proxy)
					mutex.Unlock()
					break
				}
			}
		}
	}
	if !responded {
		atomic.AddInt32(&Invalid, 1)
	}
}

func GetInvalid() int {
	return int(atomic.LoadInt32(&Invalid))
}

func GetThreadsActive() int {
	return int(atomic.LoadInt32(&threadsActive))
}

func GetQueue() ProxyQueue {
	return proxyQueue
}

func StopThreads() {
	stop = true
}

func GetProxyProtocolCountMap() map[string]map[int]int {
	var countMap map[string]map[int]int
	mutex.Lock()
	countMap = ProxyProtocolCountMap
	mutex.Unlock()
	return countMap
}
