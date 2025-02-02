package helper

import (
	_ "embed"
	"net"
	"regexp"
	"strings"
	"sync"

	"github.com/oschwald/geoip2-golang"
)

//go:embed GeoLite2-ASN.mmdb
var geoLiteASNDB []byte

//go:embed GeoLite2-Country.mmdb
var geoLiteCountryDB []byte

var (
	countryDB   *geoip2.Reader
	asnDB       *geoip2.Reader
	initOnce    sync.Once
	initSuccess bool

	datacenterOrgs = map[string]bool{
		"amazon": true, "google": true, "microsoft": true, "digitalocean": true,
		"linode": true, "hetzner": true, "ovh": true, "vultr": true, "ibm": true,
		"alibaba": true, "tencent": true, "cloudflare": true, "rackspace": true,
		"hostinger": true, "upcloud": true, "azure": true, "gcp": true, "aws": true,
	}

	residentialKeywords = regexp.MustCompile(`(?i)(dyn|pool|dsl|cust|res|ip|adsl|ppp|user|mobile|static|dhcp)`)

	ispKeywords = regexp.MustCompile(`(?i)(isp|broadband|telecom|communications|networks|carrier)`)
)

func init() {
	initOnce.Do(func() {
		var err error
		if len(geoLiteCountryDB) > 0 {
			countryDB, err = geoip2.FromBytes(geoLiteCountryDB)
			if err != nil {
				countryDB = nil
			}
		}

		if len(geoLiteASNDB) > 0 {
			asnDB, err = geoip2.FromBytes(geoLiteASNDB)
			if err != nil {
				asnDB = nil
			}
		}

		initSuccess = (countryDB != nil && asnDB != nil)
	})
}

func GetCountryCode(ipAddress string) string {
	if !initSuccess {
		return "unknown"
	}

	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return "unknown"
	}

	record, err := countryDB.Country(ip)
	if err != nil {
		return "unknown"
	}

	return record.Country.IsoCode
}

func DetermineProxyType(ipAddress string) string {
	if !initSuccess {
		return "unknown"
	}

	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return "unknown"
	}

	// First check reverse DNS for residential indicators
	names, _ := net.LookupAddr(ipAddress)
	for _, name := range names {
		if residentialKeywords.MatchString(name) {
			return "Residential"
		}
	}

	// Check ASN information
	asnRecord, err := asnDB.ASN(ip)
	if err != nil {
		return "unknown"
	}

	org := strings.ToLower(asnRecord.AutonomousSystemOrganization)

	// Check for datacenter organizations
	for keyword := range datacenterOrgs {
		if strings.Contains(org, keyword) {
			return "Datacenter"
		}
	}

	// Check for ISP indicators in ASN organization
	if ispKeywords.MatchString(org) {
		return "ISP"
	}

	// Final check for common residential ASN patterns
	if strings.Contains(org, "customer") || strings.Contains(org, "residential") {
		return "Residential"
	}

	// Default to ISP for unknown organizations
	return "ISP"
}
