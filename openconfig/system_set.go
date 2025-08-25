package openconfig

type SystemFeatureHandler func(op string, sys *System) []string

var systemFeatureHandlers = map[string]SystemFeatureHandler{
	"hostname": handleSystemHostname,
	"clock":    handleSystemClock,
	"ntp":      handleSystemNTP,
	"dns":      handleSystemDNS,
}

// SystemToMikrotikCmdsRegistry dispatches to feature handlers for get/set
func SystemToMikrotikCmdsRegistry(op string, sys *System) []string {
	var cmds []string
	for feature, handler := range systemFeatureHandlers {
		// Only call handler if the feature is present in the struct
		switch feature {
		case "hostname":
			if sys != nil && sys.Hostname != nil {
				cmds = append(cmds, handler(op, sys)...)
			}
		case "clock":
			if sys != nil && sys.Clock != nil {
				cmds = append(cmds, handler(op, sys)...)
			}
		case "ntp":
			if sys != nil && sys.NTP != nil {
				cmds = append(cmds, handler(op, sys)...)
			}
		case "dns":
			if sys != nil && sys.DNS != nil {
				cmds = append(cmds, handler(op, sys)...)
			}
		}
	}
	return cmds
}

func handleSystemHostname(op string, sys *System) []string {
	switch op {
	case "get":
		return []string{"/system/identity/print"}
	case "set":
		if sys.Hostname != nil && *sys.Hostname != "" {
			return []string{"/system/identity/set", "=name=" + *sys.Hostname}
		}
	}
	return nil
}

func handleSystemClock(op string, sys *System) []string {
	switch op {
	case "set":
		if sys.Clock != nil {
			if sys.Clock.TimezoneName != nil && *sys.Clock.TimezoneName != "" {
				return []string{"/system/clock/set", "=time-zone-name=" + *sys.Clock.TimezoneName}
			}
			if sys.Clock.TimezoneUTCOffset != nil && *sys.Clock.TimezoneUTCOffset != "" {
				// MikroTik does not support setting timezone-utc-offset directly
				// Return a special marker or error command to indicate unsupported
				return []string{"UNSUPPORTED: system/clock/timezone-utc-offset set operation is not supported on MikroTik"}
			}
		}
	}
	return nil
}

func handleSystemNTP(op string, sys *System) []string {
	var cmds []string
	switch op {
	case "set":
		if sys.NTP != nil {
			if sys.NTP.Enabled != nil {
				if *sys.NTP.Enabled {
					cmds = append(cmds, "/system/ntp/client/set", "enabled=yes")
				} else {
					cmds = append(cmds, "/system/ntp/client/set", "enabled=no")
				}
			}
			if sys.NTP.Servers != nil {
				for _, s := range sys.NTP.Servers.Server {
					if s.Address != nil && *s.Address != "" {
						cmds = append(cmds, "/system/ntp/client/servers/add", "address="+*s.Address)
					}
				}
			}
		}
	case "get":
		cmds = append(cmds, "/system/ntp/client/print")
	}
	return cmds
}

func handleSystemDNS(op string, sys *System) []string {
	if op == "set" && sys.DNS != nil && sys.DNS.Servers != nil && len(sys.DNS.Servers.Server) > 0 {
		servers := ""
		for i, s := range sys.DNS.Servers.Server {
			if i > 0 {
				servers += ","
			}
			servers += s
		}
		return []string{"/ip/dns/set servers=" + servers}
	}
	return nil
}
