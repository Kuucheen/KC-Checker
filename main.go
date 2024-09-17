package main

import (
	"KC-Checker/charm"
	"KC-Checker/common"
	"KC-Checker/helper"
	"github.com/jwalton/go-supportscolor"
	"runtime/debug"
)

func main() {
	//Lets the terminal on Windows 10 support true color
	supportscolor.Stdout()

	common.ReadSettings()

	common.GetLocalIP()

	helper.GetBlacklisted()

	if common.GetConfig().DebugMaxThreads > 0 {
		debug.SetMaxThreads(common.GetConfig().DebugMaxThreads)
	}

	charm.RunHostsDisplay()
}
