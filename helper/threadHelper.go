package helper

import (
	"KC-Checker/common"
	"fmt"
	"sync"
	"sync/atomic"
)

var (
	elite         []*proxy
	anonymous     []*proxy
	transparent   []*proxy
	allProxies    []*proxy
	invalid       int32
	threadsActive int32
	mutex         sync.Mutex
	wg            sync.WaitGroup
)

func Dispatcher(proxies []*proxy) {
	threads := common.GetConfig().Threads
	allProxies = proxies

	fmt.Println("starting dispatcher")

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
			fmt.Println("Failed ", proxy.checks)
			proxy.checks++
			continue
		}

		level := GetProxyLevel(body)

		mutex.Lock()
		switch level {
		case 1:
			transparent = append(transparent, proxy)
		case 2:
			anonymous = append(anonymous, proxy)
		case 3:
			elite = append(elite, proxy)
		default:
			atomic.AddInt32(&invalid, 1)
		}
		defer mutex.Unlock()

		fmt.Println("Proxy level: ", level, proxy.full)

		defer atomic.AddInt32(&threadsActive, -1)
		wg.Done()
		return
	}

	defer atomic.AddInt32(&threadsActive, -1)
	atomic.AddInt32(&invalid, 1)
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
	return int(atomic.LoadInt32(&invalid))
}
