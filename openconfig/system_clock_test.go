package openconfig

import (
	"reflect"
	"testing"
)

func TestSystemClock_MikroTikCmd_Set(t *testing.T) {
	tz := "America/New_York"
	sys := &System{
		Clock: &SystemClock{
			TimezoneName: &tz,
		},
	}
	cmds := handleSystemClock("set", sys)
	expected := []string{"/system/clock/set", "=time-zone-name=America/New_York"}
	if !reflect.DeepEqual(cmds, expected) {
		t.Errorf("expected %v, got %v", expected, cmds)
	}
}

func TestSystemClock_MikroTikCmd_Set_Empty(t *testing.T) {
	sys := &System{
		Clock: &SystemClock{},
	}
	cmds := handleSystemClock("set", sys)
	if cmds != nil && len(cmds) > 0 {
		t.Errorf("expected no commands, got %v", cmds)
	}
}
