package helper

import (
	"KC-Checker/common"
	"golang.org/x/net/context"
	proxy2 "golang.org/x/net/proxy"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func GetProxyLevel(html string) int {
	//When the headers contain UserIp proxy is transparent
	if strings.Contains(html, common.UserIP) {
		return 1
	}

	//When containing one of these headers proxy is anonymous
	proxyVars := []string{"HTTP_X_FORWARDED_FOR", "HTTP_FORWARDED", "HTTP_VIA", "HTTP_X_PROXY_ID"}

	for _, value := range proxyVars {
		if strings.Contains(html, value) {
			return 2
		}
	}

	//Proxy is elite
	return 3

}

func Request(proxy *Proxy) (string, int) {
	return RequestCustom(proxy, common.FastestJudge)
}

// RequestCustom makes a request to the provided siteUrl with the provided proxy
func RequestCustom(proxy *Proxy, siteUrl string) (string, int) {
	//Errors would destroy the whole display while checking
	log.SetOutput(io.Discard)

	proxyURL, err := url.Parse(GetTypeName() + "://" + proxy.Full)
	if err != nil {
		return "Error parsing proxy URL", -1
	}

	var transport *http.Transport

	switch GetTypeName() {
	case "http":
		transport = &http.Transport{Proxy: http.ProxyURL(proxyURL), DisableKeepAlives: true}
	case "socks4", "socks5":
		//udp doesn't work for some reason
		dialer, err := proxy2.SOCKS5("tcp", proxy.Full, nil, proxy2.Direct)
		if err != nil {
			return "Error creating SOCKS dialer", -1
		}
		transport = &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			}, DisableKeepAlives: true,
		}
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Millisecond * time.Duration(common.GetConfig().Timeout),
	}

	req, err := http.NewRequest("GET", siteUrl, nil)
	if err != nil {
		return "Error creating HTTP request", -1
	}

	req.Header.Set("Connection", "close")

	resp, err := client.Do(req)
	if err != nil {
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
