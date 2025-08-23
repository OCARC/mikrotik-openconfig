package main

import (
	"fmt"

	"github.com/go-routeros/routeros"
)

// MikroTikClient wraps the routeros.Client for command execution
// addr: "host:port", e.g. "192.168.88.1:8728"
func NewMikroTikClient(addr, user, pass string) (*routeros.Client, error) {
	return routeros.Dial(addr, user, pass)
}

// CommandRunner abstracts the RunArgs method for mocking/testing
type CommandRunner interface {
	RunArgs([]string) (*routeros.Reply, error)
}

// SendCommands sends a list of MikroTik API commands to the device
func SendCommands(client CommandRunner, cmds []string) error {
	for _, cmd := range cmds {
		_, err := client.RunArgs(parseCommand(cmd))
		if err != nil {
			return fmt.Errorf("failed to run command %q: %w", cmd, err)
		}
	}
	return nil
}

// parseCommand splits a command string into args for RunArgs
func parseCommand(cmd string) []string {
	// Simple split by space, but handle quoted args
	var args []string
	var current string
	inQuotes := false
	for i := 0; i < len(cmd); i++ {
		c := cmd[i]
		if c == '"' {
			inQuotes = !inQuotes
			continue
		}
		if c == ' ' && !inQuotes {
			if current != "" {
				args = append(args, current)
				current = ""
			}
			continue
		}
		current += string(c)
	}
	if current != "" {
		args = append(args, current)
	}
	return args
}
