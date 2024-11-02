package charm

import (
	"KC-Checker/charm/endPhase"
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
	if !helper.ContainsTypeName("socks4") && !helper.ContainsTypeName("socks5") {
		if !helper.ContainsTypeName("http") {
			common.RemoveHttpJudges()
		} else if !helper.ContainsTypeName("https") {
			common.RemoveHttpsJudges()
		}

		//No more judge left
		if !common.IsAllowedToCheck(helper.GetTypeNames()) {
			errorDisplays.PrintErrorForJudge(helper.GetTypeNames())
		}
	}

	//Check the judges for the fastest
	go common.CheckDomains()

	//Display while waiting for checking the judges
	//Also starts the main checking process in the background if finished
	hosts.Run()

	//Sets the current time & start the checker
	threadPhase.SetTime()
	threadPhase.RunBars()

	endPhase.RunEndScreen()
}
