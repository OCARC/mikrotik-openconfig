package openconfig

import (
	"reflect"
	"testing"
)

func TestSystemGetToMikrotikCmds_TimezoneUTCOffset(t *testing.T) {
	xml := `<system><clock><timezone-utc-offset/></clock></system>`
	cmds := SystemGetToMikrotikCmds(xml)
	expected := []string{"/system/clock/print"}
	if !reflect.DeepEqual(cmds, expected) {
		t.Errorf("expected %v, got %v", expected, cmds)
	}
}

func TestSystemGetToMikrotikCmds_TimezoneName(t *testing.T) {
	xml := `<system><clock><timezone-name/></clock></system>`
	cmds := SystemGetToMikrotikCmds(xml)
	expected := []string{"/system/clock/print"}
	if !reflect.DeepEqual(cmds, expected) {
		t.Errorf("expected %v, got %v", expected, cmds)
	}
}
