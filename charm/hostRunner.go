package charm

import (
	"KC-Checker/charm/errorDisplays"
	"KC-Checker/charm/hosts"
	"KC-Checker/charm/threadPhase"
	"KC-Checker/common"
	"KC-Checker/helper"
)

// RunHostsDisplay This has the complete process of the program
func RunHostsDisplay() {
	//Draw logo on top
	DrawLogo()

	//Sets the proxy type (http, socks4 or socks5) from user input (GetProxyType())
	helper.SetType(GetProxyType())

	//Select only the judges for the selected type
	if helper.GetTypeName() == "https" {
		common.RemoveHttpJudges()
	} else if helper.GetTypeName() == "http" {
		common.RemoveHttpsJudges()
	}

	//No more judge left
	if len(common.GetConfig().Judges) < 1 {
		errorDisplays.PrintErrorForJudge(helper.GetTypeName())
	}

	//Check the judges for the fastest
	go common.CheckDomains()

	//Display while waiting for checking the judges
	//Also starts the main checking process in the background if finished
	hosts.Run()

	//If there are no proxies in proxies.txt then wait for the user
	//to put some in the file
	if helper.ProxySum < 1 {
		hosts.WaitForProxies()
		go helper.Dispatcher(helper.GetCleanedProxies())
	}

	//Sets the current time & start the checker
	threadPhase.SetTime()
	threadPhase.RunBars()
}
