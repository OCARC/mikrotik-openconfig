package main

import "testing"

func TestTranslateNetconfToMikrotik_SystemNTP(t *testing.T) {
	xml := `<rpc><edit-config><target><running/></target><config><system><ntp><servers><server><address>1.1.1.1</address></server><server><address>2.2.2.2</address></server></servers><enabled>true</enabled></ntp></system></config></edit-config></rpc>`
	cmds, err := TranslateNetconfToMikrotik(xml)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found1, found2 := false, false
	for _, cmd := range cmds {
		if cmd == `/system/ntp/client/set enabled=yes` {
			found1 = true
		}
		if cmd == `/system/ntp/client/add address=1.1.1.1` || cmd == `/system/ntp/client/add address=2.2.2.2` {
			found2 = true
		}
	}
	if !found1 {
		t.Errorf("expected /system/ntp/client/set enabled=yes in commands, got %v", cmds)
	}
	if !found2 {
		t.Errorf("expected /system/ntp/client/add address=1.1.1.1 or 2.2.2.2 in commands, got %v", cmds)
	}
}

func TestTranslateNetconfToMikrotik_SystemDNS(t *testing.T) {
	xml := `<rpc><edit-config><target><running/></target><config><system><dns><servers><server>8.8.8.8</server><server>8.8.4.4</server></servers></dns></system></config></edit-config></rpc>`
	cmds, err := TranslateNetconfToMikrotik(xml)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, cmd := range cmds {
		if cmd == `/ip/dns/set servers=8.8.8.8,8.8.4.4` {
			found = true
		}
	}
	if !found {
		t.Errorf("expected /ip/dns/set servers=8.8.8.8,8.8.4.4 in commands, got %v", cmds)
	}
}
