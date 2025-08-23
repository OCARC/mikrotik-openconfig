package openconfig

// SystemToMikrotikCmds translates OpenConfig system tree to MikroTik API commands.
func SystemToMikrotikCmds(sys *System) []string {
	var cmds []string
	if sys == nil {
		return cmds
	}
	if sys.Hostname != nil && *sys.Hostname != "" {
		cmds = append(cmds, "/system/identity/set name=\""+*sys.Hostname+"\"")
	}
	if sys.Clock != nil && sys.Clock.TimezoneName != nil && *sys.Clock.TimezoneName != "" {
		cmds = append(cmds, "/system/clock/set time-zone-name=\""+*sys.Clock.TimezoneName+"\"")
	}
	// NTP
	if sys.NTP != nil {
		if sys.NTP.Enabled != nil {
			if *sys.NTP.Enabled {
				cmds = append(cmds, "/system/ntp/client/set enabled=yes")
			} else {
				cmds = append(cmds, "/system/ntp/client/set enabled=no")
			}
		}
		if sys.NTP.Servers != nil {
			for _, s := range sys.NTP.Servers.Server {
				if s.Address != nil && *s.Address != "" {
					cmds = append(cmds, "/system/ntp/client/add address="+*s.Address)
				}
			}
		}
	}
	// DNS
	if sys.DNS != nil && sys.DNS.Servers != nil && len(sys.DNS.Servers.Server) > 0 {
		servers := ""
		for i, s := range sys.DNS.Servers.Server {
			if i > 0 {
				servers += ","
			}
			servers += s
		}
		cmds = append(cmds, "/ip/dns/set servers="+servers)
	}
	return cmds
}
