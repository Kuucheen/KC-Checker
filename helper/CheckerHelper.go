package helper

import (
	"KC-Checker/checker"
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func GetProxyLevel(innerhtml string) int {

	fmt.Println(innerhtml)

	if strings.Contains(innerhtml, checker.UserIP) {
		return 1
	}

	var data map[string]string
	data = make(map[string]string)

	if !isJSON(innerhtml) {

		if !strings.Contains(innerhtml, "=") && strings.Contains(innerhtml, checker.UserIP) {
			return 1
		}

		if strings.Contains(innerhtml, "<pre>") {
			doc, _ := html.Parse(strings.NewReader(innerhtml))
			innerhtml = ExtractElementContent(doc, "pre")
		}

		//Splits values by "=" and puts it into the map
		lines := strings.Split(innerhtml, "\n")
		for _, line := range lines {
			parts := strings.Split(line, "=")
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				data[key] = value
			}
		}
	} else {

		// Unmarshal the JSON string into the map
		err := json.Unmarshal([]byte(innerhtml), &data)
		if err != nil {
			return 0
		}
	}

	return getLevel(data)
}

func getLevel(data map[string]string) int {

	//proxyVars[1] contains all header vars that are used by transparent proxies
	//proxyVars[2] all used by anonymous
	proxyVars := [][]string{{"HTTP_X_FORWARDED_FOR", "HTTP_FORWARDED"}, {"HTTP_VIA", "HTTP_X_PROXY_ID"}}

	for level := 0; level < len(proxyVars); level++ {

		//if server knows request is by a proxy
		isProxy := false

		for _, value := range proxyVars[level] {

			val, ok := data[value]

			if ok {
				//if header value is ip its transparent
				if strings.Contains(val, checker.UserIP) {
					fmt.Println(val)
					return 1
				}

				isProxy = true
			}

			if isProxy {
				return 2
			}

		}
	}

	//No information gathered by the server so proxy is elite
	return 3
}

func isJSON(str string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(str), &js) == nil
}

func Request(proxy *proxy) (string, int) {
	proxyURL, err := url.Parse(GetTypeName() + "://" + proxy.full)
	if err != nil {
		fmt.Println("Error parsing proxy URL:", err)
		return "Error parsing proxy URL", -1
	}

	client := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)},
		Timeout:   time.Second * time.Duration(checker.GetConfig().Timeout),
	}

	req, err := http.NewRequest("GET", checker.GetHosts()[0].Host, nil)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return "Error creating HTTP request", -1
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return "Error making HTTP request", -1
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}

		err := resp.Body.Close()
		if err != nil {
			fmt.Println("Error closing response body:", err)
		}
	}()

	status := resp.StatusCode

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "Error reading response body", -1
	}

	return string(resBody), status
}
