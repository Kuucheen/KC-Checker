package helper

import (
	"KC-Checker/common"
	"crypto/tls"
	"golang.org/x/net/context"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
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
	return RequestCustom(proxy, common.GetFastestJudgeForProtocol(proxy.Protocol), common.GetFastestJudgeNameForProtocol(proxy.Protocol), false)
}

// RequestCustom makes a request to the provided siteUrl with the provided proxy
func RequestCustom(proxyToCheck *Proxy, targetIp string, siteName *url.URL, isBanCheck bool) (string, int, error) {
	// Suppress logging for this operation
	log.SetOutput(io.Discard)

	proxyURL, err := url.Parse(strings.Replace(proxyToCheck.Protocol, "https", "http", 1) + "://" + proxyToCheck.Full)
	if err != nil {
		return "Error parsing proxyToCheck URL", -1, err
	}

	privateTransport := GetSharedTransport()

	dialer := &net.Dialer{
		Timeout: time.Millisecond * time.Duration(common.GetConfig().Timeout),
	}

	privateTransport.Proxy = http.ProxyURL(proxyURL)
	privateTransport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		if strings.Contains(addr, siteName.Hostname()) {
			addr = net.JoinHostPort(targetIp, siteName.Port())
		}
		return dialer.DialContext(ctx, network, addr)
	}
	privateTransport.TLSClientConfig = &tls.Config{
		ServerName:         siteName.Hostname(),
		InsecureSkipVerify: false,
	}

	client := GetClientFromPool()
	client.Transport = privateTransport

	req, err := http.NewRequest("GET", siteName.String(), nil)
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
