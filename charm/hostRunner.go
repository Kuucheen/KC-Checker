package charm

import (
	"KC-Checker/charm/hosts"
	"KC-Checker/charm/threadPhase"
	"KC-Checker/common"
	"KC-Checker/helper"
	"sync"
	"time"
)

var wg sync.WaitGroup

func RunHostsDisplay() {
	DrawLogo()

	helper.SetType(GetProxyType())

	go common.CheckDomains()
	hosts.Run()
	wg.Add(1)

	go func() {
		helper.Dispatcher(helper.ToProxies(helper.GetInput("proxies.txt")))
		defer wg.Done()
	}()

	time.Sleep(time.Second * 2)
	threadPhase.RunBars()
	//
	//wg.Wait()
	//threadPhase.SetFinished()
}
