package helper

import (
	"KC-Checker/common"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	clientPool      sync.Pool
	sharedTransport *http.Transport
)

func initClientPool() {
	sharedTransport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(common.GetConfig().Timeout) * time.Millisecond,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          1000, // Adjust based on your needs
		MaxIdleConnsPerHost:   100,  // Adjust based on your needs
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableKeepAlives:     !common.GetConfig().KeepAlive,
	}

	clientPool = sync.Pool{
		New: func() interface{} {
			return &http.Client{
				Transport: sharedTransport.Clone(),
				Timeout:   time.Duration(common.GetConfig().Timeout) * time.Millisecond,
			}
		},
	}
}

// Thread-safe borrowing and returning clients
func GetClientFromPool() *http.Client {
	return clientPool.Get().(*http.Client)
}

func ReturnClientToPool(client *http.Client) {
	clientPool.Put(client)
}

func GetSharedTransport() *http.Transport {
	return sharedTransport.Clone()
}
