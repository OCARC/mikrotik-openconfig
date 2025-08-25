package main

import (
	"errors"
	"testing"

	"github.com/go-routeros/routeros"
)

// mockClient implements the minimal routeros.Client interface for testing
// Only RunArgs is used in SendCommands
type mockClient struct {
	calls [][]string
	fail  bool
}

func (m *mockClient) RunArgs(args []string) (*routeros.Reply, error) {
	m.calls = append(m.calls, args)
	if m.fail {
		return nil, errors.New("mock failure")
	}
	return &routeros.Reply{}, nil
}

func TestSendCommands_Success(t *testing.T) {
	mc := &mockClient{}
	cmds := []string{"/system/identity/set name=\"router1\"", "/system/clock/set time-zone-name=\"Europe/London\""}
	err := SendCommands(mc, cmds)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(mc.calls) != 2 {
		t.Errorf("expected 2 calls, got %d", len(mc.calls))
	}
}

func TestSendCommands_Failure(t *testing.T) {
	mc := &mockClient{fail: true}
	cmds := []string{"/system/identity/set name=\"router1\""}
	err := SendCommands(mc, cmds)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseCommand(t *testing.T) {
	cmd := "/system/identity/set name=\"router1\""
	args := parseCommand(cmd)
	expected := []string{"/system/identity/set", "name=router1"}
	if len(args) != len(expected) {
		t.Fatalf("expected %d args, got %d", len(expected), len(args))
	}
	for i := range args {
		if args[i] != expected[i] {
			t.Errorf("arg %d: expected %q, got %q", i, expected[i], args[i])
		}
	}
}
