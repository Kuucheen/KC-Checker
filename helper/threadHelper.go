package helper

import (
	"KC-Checker/common"
	"sync"
	"sync/atomic"
)

var (
	elite            []*proxy
	EliteCount       int32
	anonymous        []*proxy
	AnonymousCount   int32
	transparent      []*proxy
	TransparentCount int32
	allProxies       []*proxy
	Invalid          int32
	threadsActive    int32
	mutex            sync.Mutex
	wg               sync.WaitGroup
)

func Dispatcher(proxies []*proxy) {
	threads := common.GetConfig().Threads
	allProxies = proxies

	for len(allProxies) > 0 {
		if int(atomic.LoadInt32(&threadsActive)) <= threads {
			wg.Add(1)
			go check(allProxies[0])
			allProxies = allProxies[1:]
		}
	}

	wg.Wait()
}

func check(proxy *proxy) {
	atomic.AddInt32(&threadsActive, 1)

	for proxy.checks <= common.GetConfig().Retries {
		body, status := Request(proxy)

		if status >= 400 || status == -1 {
			proxy.checks++
			continue
		}

		level := GetProxyLevel(body)

		mutex.Lock()
		switch level {
		case 1:
			transparent = append(transparent, proxy)
			atomic.AddInt32(&TransparentCount, 1)
		case 2:
			anonymous = append(anonymous, proxy)
			atomic.AddInt32(&AnonymousCount, 1)
		case 3:
			elite = append(elite, proxy)
			atomic.AddInt32(&EliteCount, 1)
		default:
			atomic.AddInt32(&Invalid, 1)
		}
		defer mutex.Unlock()

		defer atomic.AddInt32(&threadsActive, -1)
		wg.Done()
		return
	}

	defer atomic.AddInt32(&threadsActive, -1)
	atomic.AddInt32(&Invalid, 1)
	wg.Done()
	return
}

func GetFinishedProxies() map[string][]*proxy {
	return map[string][]*proxy{
		"transparent": transparent,
		"anonymous":   anonymous,
		"elite":       elite,
	}
}

func GetInvalid() int {
	return int(atomic.LoadInt32(&Invalid))
}
