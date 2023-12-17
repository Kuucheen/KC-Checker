package main

import (
	"KC-Checker/charm"
	"KC-Checker/common"
	"github.com/jwalton/go-supportscolor"
)

func main() {
	supportscolor.Stdout()

	common.ReadSettings()

	common.GetLocalIP()

	charm.RunHostsDisplay()
}
