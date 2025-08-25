package openconfig

import (
	"encoding/xml"
)

// SystemGetToMikrotikCmds parses the filter XML and dispatches to the unified handler for get operations.
func SystemGetToMikrotikCmds(filterXML string) []string {
	type systemFilter struct {
		Hostname *struct{} `xml:"hostname"`
		Clock    *struct {
			TimezoneName      *struct{} `xml:"timezone-name"`
			TimezoneUTCOffset *struct{} `xml:"timezone-utc-offset"`
		} `xml:"clock"`
		NTP *struct {
			Enabled *struct{} `xml:"enabled"`
			Servers *struct{} `xml:"servers"`
		} `xml:"ntp"`
	}
	var sys systemFilter
	_ = xml.Unmarshal([]byte(filterXML), &sys)

	var cmds []string
	if sys.Hostname != nil {
		h := &SystemHostname{}
		cmds = append(cmds, h.MikroTikCmd("get")...)
	}
	if sys.Clock != nil && (sys.Clock.TimezoneUTCOffset != nil || sys.Clock.TimezoneName != nil) {
		// MikroTik uses gmt-offset and time-zone-name in /system/clock/print
		cmds = append(cmds, "/system/clock/print")
	}
	if sys.NTP != nil {
		// Use the existing NTP handler for GET operations
		ntpSys := &System{NTP: &SystemNTP{}}
		cmds = append(cmds, handleSystemNTP("get", ntpSys)...)
	}
	// Add more features as needed, using their unified handler
	return cmds
}
