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

	threadPhase.RunBars()
	//
	//wg.Wait()
	//threadPhase.SetFinished()
}
