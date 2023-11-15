package main

import (
	"KC-Checker/helper"
	"fmt"
)

type test struct {
	val int
}

func main() {
	d := []byte("192.532.213.33")

	helper.Write("output/socks4/elite.txt", d)

	stringProxyArr := helper.GetInput()

	if len(stringProxyArr) == 0 {
		//TODO charm textinput
	}

	proxyArr := helper.ToProxies(stringProxyArr)
	fmt.Print(proxyArr)

}
