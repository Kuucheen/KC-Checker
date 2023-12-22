package charm

import (
	"KC-Checker/charm/hosts"
	"KC-Checker/charm/threadPhase"
	"KC-Checker/common"
	"KC-Checker/helper"
)

func RunHostsDisplay() {
	DrawLogo()

	helper.SetType(GetProxyType())

	go common.CheckDomains()
	hosts.Run()

	if helper.ProxySum < 1 {
		hosts.WaitForProxies()
		go helper.Dispatcher(helper.GetCleanedProxies())
	}

	threadPhase.RunBars()
}
