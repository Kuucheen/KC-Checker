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

	debug.SetMaxThreads(999999999)

	common.ReadSettings()

	common.GetLocalIP()

	helper.GetBlacklisted()

	charm.RunHostsDisplay()
}
