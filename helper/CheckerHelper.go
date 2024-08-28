package helper

import (
	"KC-Checker/common"
	"golang.org/x/net/context"
	"golang.org/x/net/proxy"
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

func Request(proxy *Proxy) (string, int, error) {
	return RequestCustom(proxy, common.GetFastestJudgeForProtocol(proxy.Protocol))
}

// RequestCustom makes a request to the provided siteUrl with the provided proxy
func RequestCustom(proxyToCheck *Proxy, siteUrl string) (string, int, error) {
	//Errors would destroy the whole display while checking
	log.SetOutput(io.Discard)

	proxyURL, err := url.Parse(strings.Replace(proxyToCheck.Protocol, "https", "http", 1) +
		"://" + proxyToCheck.Full)
	if err != nil {
		return "Error parsing proxyToCheck URL", -1, err
	}

	var transport *http.Transport

	switch proxyToCheck.Protocol {
	case "http", "https":
		transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	case "socks4", "socks5":
		//udp doesn't work for some reason
		dialer, err := proxy.SOCKS5("tcp", proxyToCheck.Full, nil, proxy.Direct)
		if err != nil {
			return "Error creating SOCKS dialer", -1, err
		}
		transport = &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			},
			DisableKeepAlives: !common.GetConfig().KeepAlive,
		}
	}

	transport.DisableKeepAlives = !common.GetConfig().KeepAlive
	transport.MaxIdleConns = 3
	transport.IdleConnTimeout = time.Duration(common.GetConfig().Timeout) * time.Millisecond

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Millisecond * time.Duration(common.GetConfig().Timeout),
	}

	req, err := http.NewRequest("GET", siteUrl, nil)
	if err != nil {
		return "Error creating HTTP request", -1, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "Error making HTTP request", -1, err
	}

	status := resp.StatusCode

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Error reading response body", -1, err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			return
		}
	}()

	return string(resBody), status, nil
}
