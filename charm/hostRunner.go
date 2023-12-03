package charm

import (
	"KC-Checker/helper"
	"sync"
)

var wg sync.WaitGroup

func RunHostsDisplay() {
	go DrawLogo()

	helper.SetType(GetProxyType())

	//go hosts.Run()
	//hostsList := common.CheckDomains()
	//hosts.SetFinished()
	//DrawLogo()
	//hosts.DisplayHosts(hostsList)
	//wg.Add(1)
	//
	//go func() {
	//	helper.Dispatcher(helper.ToProxies(helper.GetInput()))
	//	defer wg.Done()
	//}()
	//
	//time.Sleep(time.Second * 2)
	//threadPhase.RunBars()
	//
	//wg.Wait()
	//threadPhase.SetFinished()
}
