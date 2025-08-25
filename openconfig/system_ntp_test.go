package openconfig

import (
	"reflect"
	"testing"
)

func TestSystemNTP_MikroTikCmd_Set_Enabled(t *testing.T) {
	enabled := true
	sys := &System{
		NTP: &SystemNTP{Enabled: &enabled},
	}
	cmds := handleSystemNTP("set", sys)
	expected := []string{"/system/ntp/client/set", "enabled=yes"}
	if !reflect.DeepEqual(cmds, expected) {
		t.Errorf("expected %v, got %v", expected, cmds)
	}
}

func TestSystemNTP_MikroTikCmd_Set_Disabled(t *testing.T) {
	enabled := false
	sys := &System{
		NTP: &SystemNTP{Enabled: &enabled},
	}
	cmds := handleSystemNTP("set", sys)
	expected := []string{"/system/ntp/client/set", "enabled=no"}
	if !reflect.DeepEqual(cmds, expected) {
		t.Errorf("expected %v, got %v", expected, cmds)
	}
}

func TestSystemNTP_MikroTikCmd_Set_Servers(t *testing.T) {
	addr1 := "1.2.3.4"
	addr2 := "5.6.7.8"
	sys := &System{
		NTP: &SystemNTP{
			Servers: &SystemNTPServers{
				Server: []SystemNTPServer{
					{Address: &addr1},
					{Address: &addr2},
				},
			},
		},
	}
	cmds := handleSystemNTP("set", sys)
	expected := []string{"/system/ntp/client/servers/add", "address=1.2.3.4", "/system/ntp/client/servers/add", "address=5.6.7.8"}
	if !reflect.DeepEqual(cmds, expected) {
		t.Errorf("expected %v, got %v", expected, cmds)
	}
}

func TestSystemNTP_MikroTikCmd_Get(t *testing.T) {
	sys := &System{NTP: &SystemNTP{}}
	cmds := handleSystemNTP("get", sys)
	expected := []string{"/system/ntp/client/print"}
	if !reflect.DeepEqual(cmds, expected) {
		t.Errorf("expected %v, got %v", expected, cmds)
	}
}
