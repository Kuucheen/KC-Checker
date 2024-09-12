package helper

import (
	_ "embed"
	"github.com/oschwald/geoip2-golang"
	"net"
)

//go:embed GeoLite2-Country.mmdb
var geoLiteDB []byte

var geoIP2DB *geoip2.Reader

func init() {
	if len(geoLiteDB) > 0 {
		var err error
		geoIP2DB, err = geoip2.FromBytes(geoLiteDB)
		if err != nil {
			geoIP2DB = nil
		}
	} else {
		geoIP2DB = nil
	}
}

func GetCountryCode(proxy *Proxy) string {
	if geoIP2DB == nil {
		return "unknown"
	}

	ip := net.ParseIP(proxy.Ip)
	if ip == nil {
		return "unknown"
	}

	record, err := geoIP2DB.Country(ip)
	if err != nil {
		return "unknown"
	}

	return record.Country.IsoCode
}
