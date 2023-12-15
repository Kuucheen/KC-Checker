package main

import (
	"KC-Checker/charm"
	"KC-Checker/common"
)

func main() {
	common.ReadSettings()

	common.GetLocalIP()

	//hostsList := common.CheckDomains()
	//hosts.DisplayHosts(hostsList)

	charm.RunHostsDisplay()
}
