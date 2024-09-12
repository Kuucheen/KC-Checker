package helper

import (
	"KC-Checker/common"
	"crypto/tls"
	"golang.org/x/net/context"
	"golang.org/x/net/proxy"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
)

var (
	proxyHeader = []string{"HTTP_X_FORWARDED_FOR", "HTTP_FORWARDED", "HTTP_VIA", "HTTP_X_PROXY_ID"}
)

func GetProxyLevel(html string) int {
	//When the headers contain UserIp proxy is transparent
	if strings.Contains(html, common.UserIP) {
		return 1
	}

	//When containing one of these headers proxy is anonymous
	for _, header := range proxyHeader {
		if strings.Contains(html, header) {
			return 2
		}
	}

	//Proxy is elite
	return 3

}

func Request(proxy *Proxy) (string, int, error) {
	return RequestCustom(proxy, common.GetFastestJudgeForProtocol(proxy.Protocol), false)
}

// RequestCustom makes a request to the provided siteUrl with the provided proxy
func RequestCustom(proxyToCheck *Proxy, siteUrl string, isBanCheck bool) (string, int, error) {
	// Suppress logging for this operation
	log.SetOutput(io.Discard)

	proxyURL, err := url.Parse(strings.Replace(proxyToCheck.Protocol, "https", "http", 1) + "://" + proxyToCheck.Full)
	if err != nil {
		return "Error parsing proxyToCheck URL", -1, err
	}

	privateTransport := GetSharedTransport()

	switch proxyToCheck.Protocol {
	case "http", "https":
		privateTransport.Proxy = http.ProxyURL(proxyURL)

		if proxyToCheck.Protocol == "https" {
			privateTransport.TLSClientConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}

	case "socks4", "socks5":
		dialer, err := proxy.SOCKS5("tcp", proxyToCheck.Full, nil, proxy.Direct)
		if err != nil {
			return "Error creating SOCKS dialer", -1, err
		}
		privateTransport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		}
		privateTransport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	client := GetClientFromPool()
	client.Transport = privateTransport

	req, err := http.NewRequest("GET", siteUrl, nil)
	if err != nil {
		ReturnClientToPool(client)
		return "Error creating HTTP request", -1, err
	}

	resp, err := client.Do(req)
	ReturnClientToPool(client)
	if err != nil {
		return "Error making HTTP request", -1, err
	}

	defer resp.Body.Close()

	status := resp.StatusCode
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Error reading response body", -1, err
	}

	html := string(resBody)
	if !isBanCheck && !common.CheckForValidResponse(html) {
		return "Invalid response", -1, nil
	}

	return html, status, nil
}
