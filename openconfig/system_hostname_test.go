package openconfig

import (
	"reflect"
	"testing"
)

func TestSystemHostname_MikroTikCmd_Get(t *testing.T) {
	h := &SystemHostname{}
	cmds := h.MikroTikCmd("get")
	expected := []string{"/system/identity/print"}
	if !reflect.DeepEqual(cmds, expected) {
		t.Errorf("expected %v, got %v", expected, cmds)
	}
}

func TestSystemHostname_MikroTikCmd_Set(t *testing.T) {
	name := "router1"
	h := &SystemHostname{Value: &name}
	cmds := h.MikroTikCmd("set")
	expected := []string{"/system/identity/set name=\"router1\""}
	if !reflect.DeepEqual(cmds, expected) {
		t.Errorf("expected %v, got %v", expected, cmds)
	}
}
