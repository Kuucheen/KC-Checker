package charm

import (
	"KC-Checker/charm/spinner"
	"KC-Checker/common"
)

func RunHostsDisplay() {
	go spinner.Run()
	hosts := common.CheckDomains()
	spinner.SetFinished()
	DrawLogo()
	DisplayHosts(hosts)
}
