package main

import "testing"

func TestTranslateNetconfToMikrotik_SystemHostname(t *testing.T) {
	xml := `<rpc><edit-config><target><running/></target><config><system><hostname>router1</hostname></system></config></edit-config></rpc>`
	cmds, err := TranslateNetconfToMikrotik(xml)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, cmd := range cmds {
		if cmd == `/system/identity/set name="router1"` {
			found = true
		}
	}
	if !found {
		t.Errorf("expected /system/identity/set name=\"router1\" in commands, got %v", cmds)
	}
}

func TestTranslateNetconfToMikrotik_SystemClock(t *testing.T) {
	xml := `<rpc><edit-config><target><running/></target><config><system><clock><timezone-name>Europe/London</timezone-name></clock></system></config></edit-config></rpc>`
	cmds, err := TranslateNetconfToMikrotik(xml)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, cmd := range cmds {
		if cmd == `/system/clock/set time-zone-name="Europe/London"` {
			found = true
		}
	}
	if !found {
		t.Errorf("expected /system/clock/set time-zone-name=\"Europe/London\" in commands, got %v", cmds)
	}
}
