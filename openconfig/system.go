package openconfig

import "fmt"

// SystemHostname represents the OpenConfig system/hostname feature
type SystemHostname struct {
	Value *string `xml:"hostname"`
}

func (h *SystemHostname) MikroTikCmd(op string) []string {
	switch op {
	case "get":
		return []string{"/system/identity/print"}
	case "set":
		if h.Value != nil && *h.Value != "" {
			return []string{fmt.Sprintf("/system/identity/set name=\"%s\"", *h.Value)}
		}
	}
	return nil
}

type System struct {
	Hostname *string        `xml:"hostname"`
	Clock    *SystemClock   `xml:"clock"`
	NTP      *SystemNTP     `xml:"ntp"`
	DNS      *SystemDNS     `xml:"dns"`
	AAA      *SystemAAA     `xml:"aaa"`
	Logging  *SystemLogging `xml:"logging"`
	// Add more fields as needed (e.g., ssh, telnet, etc)
}

type SystemClock struct {
	TimezoneName      *string `xml:"timezone-name"`
	TimezoneUTCOffset *string `xml:"timezone-utc-offset"`
	// Add more fields as needed
}

type SystemNTP struct {
	Servers *SystemNTPServers `xml:"servers"`
	Enabled *bool             `xml:"enabled"`
}

type SystemNTPServers struct {
	Server []SystemNTPServer `xml:"server"`
}

type SystemNTPServer struct {
	Address *string `xml:"address"`
	Port    *uint16 `xml:"port"`
	// Add more fields as needed
}

type SystemDNS struct {
	Servers *SystemDNSServers `xml:"servers"`
}

type SystemDNSServers struct {
	Server []string `xml:"server"`
}

type SystemAAA struct {
	Authentication *SystemAuthentication `xml:"authentication"`
}

type SystemAuthentication struct {
	Users *SystemUsers `xml:"users"`
}

type SystemUsers struct {
	User []SystemUser `xml:"user"`
}

type SystemUser struct {
	Username *string `xml:"username"`
	Password *string `xml:"password"`
	Role     *string `xml:"role"`
}

type SystemLogging struct {
	Console       *SystemLoggingConsole       `xml:"console"`
	RemoteServers *SystemLoggingRemoteServers `xml:"remote-servers"`
}

type SystemLoggingConsole struct {
	Severity *string `xml:"severity"`
}

type SystemLoggingRemoteServers struct {
	RemoteServer []SystemLoggingRemoteServer `xml:"remote-server"`
}

type SystemLoggingRemoteServer struct {
	Host     *string `xml:"host"`
	Port     *uint16 `xml:"port"`
	Severity *string `xml:"severity"`
}
