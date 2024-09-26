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

func init() {
	configTransport := common.GetConfig().Transport

	sharedTransport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(common.GetConfig().Timeout) * time.Millisecond,
			KeepAlive: time.Duration(configTransport.KeepAliveSeconds) * time.Second,
		}).DialContext,
		DisableKeepAlives:     !configTransport.KeepAlive,
		MaxIdleConns:          configTransport.MaxIdleConns,
		MaxIdleConnsPerHost:   configTransport.MaxIdleConnsPerHost,
		IdleConnTimeout:       time.Duration(configTransport.IdleConnTimeout) * time.Second,
		TLSHandshakeTimeout:   time.Duration(configTransport.TLSHandshakeTimeout) * time.Second,
		ExpectContinueTimeout: time.Duration(configTransport.ExpectContinueTimeout) * time.Second,
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

// GetClientFromPool Thread-safe borrowing and returning clients
// This also creates a new http Client if non are available
func GetClientFromPool() *http.Client {
	return clientPool.Get().(*http.Client)
}

func ReturnClientToPool(client *http.Client) {
	clientPool.Put(client)
}

func GetSharedTransport() *http.Transport {
	return sharedTransport.Clone()
}
