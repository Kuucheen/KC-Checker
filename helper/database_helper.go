package helper

import (
	_ "embed"
	"github.com/oschwald/geoip2-golang"
	"net"
	"strings"
)

//go:embed GeoLite2-ASN.mmdb
var geoLiteASNDB []byte

//go:embed GeoLite2-Country.mmdb
var geoLiteCountryDB []byte

var geoIP2CountryDB *geoip2.Reader
var geoIP2ASNDB *geoip2.Reader

func init() {
	// Load the Country database
	if len(geoLiteCountryDB) > 0 {
		var err error
		geoIP2CountryDB, err = geoip2.FromBytes(geoLiteCountryDB)
		if err != nil {
			geoIP2CountryDB = nil
		}
	} else {
		geoIP2CountryDB = nil
	}

	// Load the ASN database
	if len(geoLiteASNDB) > 0 {
		var err error
		geoIP2ASNDB, err = geoip2.FromBytes(geoLiteASNDB)
		if err != nil {
			geoIP2ASNDB = nil
		}
	} else {
		geoIP2ASNDB = nil
	}
}

// GetCountryCode returns the ISO code of a country based on the IP address
func GetCountryCode(ipAddress string) string {
	if geoIP2CountryDB == nil {
		return "unknown"
	}

	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return "unknown"
	}

	record, err := geoIP2CountryDB.Country(ip)
	if err != nil {
		return "unknown"
	}

	return record.Country.IsoCode
}

// DetermineProxyType classifies the IP as ISP, Datacenter, or Residential
func DetermineProxyType(ipAddress string) string {
	if geoIP2ASNDB == nil {
		return "unknown"
	}

	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return "unknown"
	}

	record, err := geoIP2ASNDB.ASN(ip)
	if err != nil {
		return "unknown"
	}

	// Classify based on the ASN organization name
	org := strings.ToLower(record.AutonomousSystemOrganization)

	switch {
	case strings.Contains(org, "amazon") || strings.Contains(org, "digitalocean") ||
		strings.Contains(org, "google") || strings.Contains(org, "microsoft") ||
		strings.Contains(org, "linode") || strings.Contains(org, "ovh") ||
		strings.Contains(org, "choopa") || strings.Contains(org, "leaseweb"):
		return "Datacenter"
	case strings.Contains(org, "comcast") || strings.Contains(org, "verizon") ||
		strings.Contains(org, "at&t") || strings.Contains(org, "charter") ||
		strings.Contains(org, "spectrum") || strings.Contains(org, "centurylink") ||
		strings.Contains(org, "bt group") || strings.Contains(org, "telecom"):
		return "ISP"
	default:
		return "Residential"
	}
}
