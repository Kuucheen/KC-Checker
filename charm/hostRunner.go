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
	helper.SetType(RunType())

	go hosts.Run()
	hostsList := common.CheckDomains()
	hosts.SetFinished()
	DrawLogo()
	hosts.DisplayHosts(hostsList)
	wg.Add(1)

	go func() {
		helper.Dispatcher(helper.ToProxies(helper.GetInput()))
		defer wg.Done()
	}()

	time.Sleep(time.Second * 2)
	threadPhase.RunBars()

	wg.Wait()
	threadPhase.SetFinished()
}
