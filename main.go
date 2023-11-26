package main

import (
	"KC-Checker/common"
	"KC-Checker/helper"
	"fmt"
)

func main() {
	common.ReadSettings()

	stringProxyArr := helper.GetInput()

	if len(stringProxyArr) == 0 {
		fmt.Print("Looks like you forgot to add some proxies to proxies.txt!")
	}

	proxyArr := helper.ToProxies(stringProxyArr)

	helper.GetType()

	fmt.Println("got type: ", helper.GetTypeName())

	common.CheckDomains()

	common.GetLocalIP()

	helper.Dispatcher(proxyArr)

	fmt.Println("FINISHED!")

	fmt.Println("Elite: ", len(helper.GetFinishedProxies()["elite"]))
	fmt.Println("Anonymous: ", len(helper.GetFinishedProxies()["anonymous"]))
	fmt.Println("Transparent: ", len(helper.GetFinishedProxies()["transparent"]))
	fmt.Println("Invalid: ", helper.GetInvalid())
}
