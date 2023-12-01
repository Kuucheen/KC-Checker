package helper

import (
	"KC-Checker/common"
	"golang.org/x/net/context"
	proxy2 "golang.org/x/net/proxy"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func GetProxyLevel(innerhtml string) int {
	if strings.Contains(innerhtml, common.UserIP) {
		return 1
	}

	proxyVars := []string{"HTTP_X_FORWARDED_FOR", "HTTP_FORWARDED", "HTTP_VIA", "HTTP_X_PROXY_ID"}

	for _, value := range proxyVars {
		if strings.Contains(innerhtml, value) {
			return 2
		}
	}

	return 3

}

func Request(proxy *proxy) (string, int) {
	proxyURL, err := url.Parse(GetTypeName() + "://" + proxy.full)
	if err != nil {
		return "Error parsing proxy URL", -1
	}

	var transport *http.Transport

	switch GetTypeName() {
	case "http":
		transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	case "socks4", "socks5":
		dialer, err := proxy2.SOCKS5("tcp", proxy.full, nil, proxy2.Direct)
		if err != nil {
			return "Error creating SOCKS5 dialer", -1
		}
		transport = &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			},
		}
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Millisecond * time.Duration(common.GetConfig().Timeout),
	}

	req, err := http.NewRequest("GET", common.GetHosts()[0].Host, nil)
	if err != nil {
		return "Error creating HTTP request", -1
	}

	req.Header.Set("Connection", "close")

	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "Maximum number of open connections reached.") {
			return "Maximum number of open connections reached", -1
		}

		return "Error making HTTP request", -1
	}
	defer func() {
		if r := recover(); r != nil {
		}

		err := resp.Body.Close()
		if err != nil {
		}
	}()

	status := resp.StatusCode

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Error reading response body", -1
	}

	return string(resBody), status
}
