package helper

import (
	"KC-Checker/common"
	"sync"
	"sync/atomic"
	"time"
)

const (
	Transparent = 1
	Anonymous   = 2
	Elite       = 3
)

var (
	proxyQueue    = ProxyQueue{}
	ProxyMap      = make(map[int][]*Proxy)
	ProxyCountMap = make(map[int]int)
	stop          = false
	Invalid       int32
	threadsActive int32
	mutex         sync.Mutex
	wg            sync.WaitGroup

	HasFinished = false
)

func Dispatcher(proxies []*Proxy) {
	threads := common.GetConfig().Threads

	for len(proxies) > 0 {
		if int(atomic.LoadInt32(&threadsActive)) < threads {
			wg.Add(1)
			go check(proxies[0])
			atomic.AddInt32(&threadsActive, 1)
			proxies = proxies[1:]
		} else {
			time.Sleep(time.Millisecond * 10)
		}
		if stop {
			break
		}
	}

	wg.Wait()
	HasFinished = true
}

func check(proxy *Proxy) {
	for proxy.checks <= common.GetConfig().Retries {
		body, status := Request(proxy)

		if status >= 400 || status == -1 {
			proxy.checks++
			continue
		}

		level := GetProxyLevel(body)

		mutex.Lock()
		switch level {
		case Transparent:
			proxy.Level = Transparent
			ProxyMap[Transparent] = append(ProxyMap[Transparent], proxy)
			ProxyCountMap[Transparent]++
			proxyQueue.Enqueue(proxy)
		case Anonymous:
			proxy.Level = Anonymous
			ProxyMap[Anonymous] = append(ProxyMap[Anonymous], proxy)
			ProxyCountMap[Anonymous]++
			proxyQueue.Enqueue(proxy)
		case Elite:
			proxy.Level = Elite
			ProxyMap[Elite] = append(ProxyMap[Elite], proxy)
			ProxyCountMap[Elite]++
			proxyQueue.Enqueue(proxy)
		default:
			atomic.AddInt32(&Invalid, 1)
		}
		mutex.Unlock()

		atomic.AddInt32(&threadsActive, -1)
		wg.Done()
		return
	}

	defer func() {
		atomic.AddInt32(&threadsActive, -1)
		wg.Done()
	}()
	atomic.AddInt32(&Invalid, 1)
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
