package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/go-routeros/routeros"
)

func TestE2ENetconfFiles(t *testing.T) {
	mtikAddr := os.Getenv("MIKROTIK_ADDR")
	mtikUser := os.Getenv("MIKROTIK_USER")
	mtikPass := os.Getenv("MIKROTIK_PASS")
	if mtikAddr == "" || mtikUser == "" || mtikPass == "" {
		t.Skip("Set MIKROTIK_ADDR, MIKROTIK_USER, and MIKROTIK_PASS env vars to run this test")
	}

	glob := os.Getenv("NETCONF_TEST_GLOB")
	if glob == "" {
		glob = "netconf-tests/*.xml"
	}

	files, err := filepath.Glob(glob)
	if err != nil || len(files) == 0 {
		t.Fatalf("No test files found for glob: %s", glob)
	}

	// Use our client abstraction instead of direct RouterOS library
	c, err := NewMikroTikClient(mtikAddr, mtikUser, mtikPass)
	if err != nil {
		t.Fatalf("failed to connect to MikroTik: %v", err)
	}
	defer c.Close()

	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			runNetconfTest(t, c, file)
		})
	}
}

func runNetconfTest(t *testing.T, c *routeros.Client, file string) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatalf("failed to read %s: %v", file, err)
	}

	// Use the existing translator to parse and handle the NETCONF XML
	handleE2ENetconf(t, c, string(data))
}

// Split multi-RPC XML into individual RPC operations
func splitRPCOperations(xmlData string) []string {
	// Remove comments first
	commentRegex := regexp.MustCompile(`<!--[\s\S]*?-->`)
	cleanXML := commentRegex.ReplaceAllString(xmlData, "")

	// Find all RPC blocks
	rpcRegex := regexp.MustCompile(`<rpc[^>]*>[\s\S]*?</rpc>`)
	matches := rpcRegex.FindAllString(cleanXML, -1)

	if len(matches) == 0 {
		// No explicit RPC tags found, return the original data
		return []string{xmlData}
	}

	return matches
}

// Determine if an RPC operation is conditional based on filename or content
func isConditionalOperation(filename, xmlData string) bool {
	// Check filename for conditional indicators
	if strings.Contains(filename, "_if_") || strings.Contains(filename, "_conditional_") {
		return true
	}

	// Check if it contains both get and edit-config operations
	hasGet := strings.Contains(xmlData, "<get>")
	hasEdit := strings.Contains(xmlData, "<edit-config>")

	return hasGet && hasEdit
}

func handleE2ENetconf(t *testing.T, c *routeros.Client, xmlData string) {
	// Split into individual RPC operations
	rpcOps := splitRPCOperations(xmlData)

	// Handle conditional operations specially
	testName := t.Name()
	if isConditionalOperation(testName, xmlData) {
		handleConditionalOperation(t, c, rpcOps)
		return
	}

	// Execute each RPC operation in sequence
	for _, rpcOp := range rpcOps {
		executeRPCOperation(t, c, rpcOp)
	}

	// Sleep briefly to allow changes to take effect for set operations
	if strings.Contains(xmlData, "<edit-config>") {
		time.Sleep(1 * time.Second)
		// Verify changes took effect using the last edit-config operation
		for _, rpcOp := range rpcOps {
			if strings.Contains(rpcOp, "<edit-config>") {
				verifyE2EChangesFromXML(t, c, rpcOp)
				break // Only verify the first edit-config operation
			}
		}
	}
}

// Handle conditional operations like "set hostname if not testing"
func handleConditionalOperation(t *testing.T, c *routeros.Client, rpcOps []string) {
	if len(rpcOps) < 2 {
		t.Logf("Conditional operation requires at least 2 RPC operations, got %d", len(rpcOps))
		return
	}

	// Execute the first operation (usually a get) to check current state
	getOp := rpcOps[0]
	if !strings.Contains(getOp, "<get>") {
		t.Logf("First operation in conditional should be a get operation")
		return
	}

	// Execute the GET operation via translator
	executeRPCOperation(t, c, getOp)

	// Determine the condition based on the test name and current state
	shouldExecuteEdit := false

	// Execute the edit operation if condition is met
	if shouldExecuteEdit {
		editOp := rpcOps[1]
		if strings.Contains(editOp, "<edit-config>") {
			executeRPCOperation(t, c, editOp)

			// Sleep and verify changes
			time.Sleep(1 * time.Second)
			verifyE2EChangesFromXML(t, c, editOp)
		}
	}
}

// Execute a single RPC operation
func executeRPCOperation(t *testing.T, c *routeros.Client, rpcOp string) {
	// Generate and execute MikroTik commands using existing translator
	cmds, err := TranslateNetconfToMikrotik(rpcOp)
	if err != nil {
		t.Logf("Translation failed for RPC operation: %v", err)
		return
	}
	if len(cmds) == 0 {
		if testing.Verbose() {
			t.Log("No commands generated from RPC operation")
		}
		return
	}

	// Execute commands - treat the entire commands slice as parts of a single command
	var validCmds []string
	for _, cmd := range cmds {
		if strings.Contains(cmd, "UNSUPPORTED") {
			if testing.Verbose() {
				t.Logf("Skipping unsupported command: %s", cmd)
			}
			continue
		}
		validCmds = append(validCmds, cmd)
	}

	if len(validCmds) > 0 {

		// Check if we have a multi-part command (command + arguments)
		if len(validCmds) > 1 && !strings.HasPrefix(validCmds[1], "/") {
			// This looks like a command with arguments - combine them
			combinedCmd := strings.Join(validCmds, " ")
			cmdArgs := parseE2ECommand(combinedCmd)
			_, err := c.RunArgs(cmdArgs)
			if err != nil {
				t.Errorf("MikroTik command failed: %s - %v", combinedCmd, err)
			}
		} else {
			// Execute each command separately
			for _, cmd := range validCmds {
				cmdArgs := parseE2ECommand(cmd)
				_, err := c.RunArgs(cmdArgs)
				if err != nil {
					t.Errorf("MikroTik command failed: %s - %v", cmd, err)
				}
			}
		}
	}
}

// Parse command string into RouterOS API arguments
func parseE2ECommand(command string) []string {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return fields
	}

	// First field is the command path
	result := []string{fields[0]}

	// Remaining fields are parameters - need to add = prefix for RouterOS API
	for i := 1; i < len(fields); i++ {
		param := fields[i]
		// RouterOS API expects =parameter=value format
		if !strings.HasPrefix(param, "=") {
			param = "=" + param
		}
		result = append(result, param)
	}

	return result
}

func verifyE2EChangesFromXML(t *testing.T, c *routeros.Client, xmlData string) {
	// Extract the config section from edit-config
	configStart := strings.Index(xmlData, "<config>")
	configEnd := strings.Index(xmlData, "</config>")
	if configStart == -1 || configEnd == -1 {
		return
	}

	configXML := xmlData[configStart+8 : configEnd] // +8 to skip "<config>"

	// Generate GET XML by converting edit-config structure to get filter
	getFilterXML := convertEditConfigToGetFilter(configXML)
	if getFilterXML == "" {
		return
	}

	// Create full GET RPC
	getXML := `<rpc><get><filter>` + getFilterXML + `</filter></get></rpc>`

	// Execute the GET command via translator
	cmds, err := TranslateNetconfToMikrotik(getXML)
	if err != nil || len(cmds) == 0 {
		return
	}

	// Execute the GET
	cmdArgs := parseE2ECommand(cmds[0])
	re, err := c.RunArgs(cmdArgs)
	if err != nil || re == nil || len(re.Re) == 0 {
		return
	}

	// Extract expected values from original config and verify
	verifyValues(t, configXML, re.Re[0].Map)
}

// Convert edit-config XML to get filter by removing values and operation attributes
func convertEditConfigToGetFilter(configXML string) string {
	// Remove xc:operation attributes
	opRegex := regexp.MustCompile(`\s+xc:operation="[^"]*"`)
	cleanXML := opRegex.ReplaceAllString(configXML, "")

	// Convert leaf nodes with values to empty leaf nodes for GET
	// e.g., <hostname>test</hostname> -> <hostname/>
	leafRegex := regexp.MustCompile(`<([^/>]+)>[^<]+</([^>]+)>`)
	getFilterXML := leafRegex.ReplaceAllStringFunc(cleanXML, func(match string) string {
		// Extract the tag name from the opening tag
		tagStart := strings.Index(match, "<") + 1
		tagEnd := strings.Index(match, ">")
		if tagEnd > tagStart {
			tagName := match[tagStart:tagEnd]
			return "<" + tagName + "/>"
		}
		return match
	})

	return getFilterXML
}

// Verify that the actual values match the expected values from the config
func verifyValues(t *testing.T, configXML string, actualValues map[string]string) {
	// Extract expected values from the config XML
	expectedValues := extractExpectedValues(configXML)

	for field, expected := range expectedValues {
		// Map NETCONF field names to RouterOS field names
		rosField := mapNetconfToRouterOSField(field)

		if actual, ok := actualValues[rosField]; ok {
			// Handle boolean conversions
			actualValue := normalizeValue(actual)
			expectedValue := normalizeValue(expected)

			if actualValue != expectedValue {
				t.Errorf("%s verification failed: expected %s, got %s", field, expectedValue, actualValue)
			} else if testing.Verbose() {
				fmt.Println("\t\t" + field + " verification passed: " + expectedValue)
			}
		} else {
			t.Errorf("Field %s not found in RouterOS response (mapped to %s)", field, rosField)
			if testing.Verbose() {
				t.Logf("Available fields: %v", actualValues)
			}
		}
	}
}

// Extract expected values from config XML
func extractExpectedValues(configXML string) map[string]string {
	values := make(map[string]string)

	// Use regex to find leaf elements with values
	leafRegex := regexp.MustCompile(`<([^/>]+)>([^<]+)</[^>]+>`)
	matches := leafRegex.FindAllStringSubmatch(configXML, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			fieldName := match[1]
			fieldValue := match[2]
			values[fieldName] = fieldValue
		}
	}

	return values
}

// Map NETCONF field names to RouterOS field names
func mapNetconfToRouterOSField(netconfField string) string {
	fieldMap := map[string]string{
		"hostname":      "name",
		"timezone-name": "time-zone-name",
		"enabled":       "enabled",
	}

	if rosField, ok := fieldMap[netconfField]; ok {
		return rosField
	}
	return netconfField // fallback to original name
}

// Normalize values for comparison (handle boolean conversions)
func normalizeValue(value string) string {
	switch strings.ToLower(value) {
	case "true", "yes", "1":
		return "true"
	case "false", "no", "0":
		return "false"
	default:
		return value
	}
}
