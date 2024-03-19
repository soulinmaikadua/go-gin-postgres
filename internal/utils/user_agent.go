package utils

import (
	"regexp"
	"strings"
)

type UserAgentInfo struct {
	Browser        string `json:"browser"`
	BrowserVersion string `json:"browser_version"`
	OS             string `json:"os"`
	OSVersion      string `json:"os_version"`
	Device         string `json:"device"`
	Engine         string `json:"engine"`
}

func ParseUserAgent(userAgent string) UserAgentInfo {
	// Regular expressions to extract various details from the User-Agent string
	browserRe := regexp.MustCompile(`(Chrome|Firefox|Safari|Edge|Opera)/([^\s]+)`)
	osRe := regexp.MustCompile(`\((Windows|Macintosh|Linux|Android|iOS|iPad|iPhone);? ([^;]+)?`)
	deviceRe := regexp.MustCompile(`\((Mobile|Tablet|Desktop)\)`)

	// Extract browser name and version
	browserMatches := browserRe.FindStringSubmatch(userAgent)
	browser, browserVersion := "", ""
	if len(browserMatches) >= 3 {
		browser, browserVersion = browserMatches[1], browserMatches[2]
	}

	// Extract OS name and version
	osMatches := osRe.FindStringSubmatch(userAgent)
	os, osVersion := "", ""
	if len(osMatches) >= 3 {
		os, osVersion = osMatches[1], strings.TrimSpace(osMatches[2])
	}

	// Extract device type
	deviceMatches := deviceRe.FindStringSubmatch(userAgent)
	device := ""
	if len(deviceMatches) >= 2 {
		device = deviceMatches[1]
	}

	// Extract rendering engine
	engine := "Unknown" // Default to "Unknown"
	if strings.Contains(userAgent, "Gecko") {
		engine = "Gecko"
	} else if strings.Contains(userAgent, "WebKit") {
		engine = "WebKit"
	}

	// Create a UserAgentInfo struct with the extracted information
	return UserAgentInfo{
		Browser:        browser,
		BrowserVersion: browserVersion,
		OS:             os,
		OSVersion:      osVersion,
		Device:         device,
		Engine:         engine,
	}
}
