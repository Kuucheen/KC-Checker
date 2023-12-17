package main

import (
	"KC-Checker/charm"
	"KC-Checker/common"
)

func main() {
	common.ReadSettings()

	common.GetLocalIP()

	charm.RunHostsDisplay()
}
