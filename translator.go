package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"

	"github.com/OCARC/mikrotik-openconfig/openconfig"
)

// --- Schema Validation ---
// ValidateOpenConfigSchema validates NETCONF XML against IETF OpenConfig YANG/XSD schema.
// In production, use a YANG-to-XSD tool to generate XSDs, then validate XML using a Go XSD library.
// For now, this is a stub for demonstration and extension.
func ValidateOpenConfigSchema(xmlInput string) error {
	// TODO: Integrate with a Go XML Schema validator (e.g., github.com/lestrrat-go/libxml2 or similar)
	// and use XSDs generated from OpenConfig YANG models.
	// Example:
	//   schema, _ := xsd.Parse(xsdData)
	//   doc, _ := libxml2.ParseString(xmlInput)
	//   err := schema.Validate(doc)
	//   return err
	// For now, just return nil (no validation)
	return nil
}

// --- NETCONF and OpenConfig Structures ---
type NetconfRPC struct {
	XMLName      xml.Name      `xml:"rpc"`
	Get          *Get          `xml:"get"`
	EditConfig   *EditConfig   `xml:"edit-config"`
	DeleteConfig *DeleteConfig `xml:"delete-config"`
}

type Get struct {
	Filter Filter `xml:"filter"`
}

type Filter struct {
	Type  string `xml:"type,attr"`
	Value string `xml:",innerxml"`
}

type EditConfig struct {
	Target Target `xml:"target"`
	Config Config `xml:"config"`
}

type DeleteConfig struct {
	Target Target `xml:"target"`
}

type Target struct {
	Running string `xml:"running"`
}

type Config struct {
	System *openconfig.System `xml:"system"`
	// Extend for more OpenConfig modules
}

// --- Translation Logic ---

// TranslateNetconfToMikrotik takes NETCONF XML and returns MikroTik API commands
func TranslateNetconfToMikrotik(xmlInput string) ([]string, error) {
	// Step 1: Validate against OpenConfig schema
	if err := ValidateOpenConfigSchema(xmlInput); err != nil {
		return nil, fmt.Errorf("schema validation failed: %w", err)
	}

	// Step 2: Parse and translate
	var rpc NetconfRPC
	err := xml.Unmarshal([]byte(xmlInput), &rpc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse NETCONF XML: %w", err)
	}

	var cmds []string

	// Handle <get>
	if rpc.Get != nil {
		cmds = append(cmds, handleGet(rpc.Get)...)
	}

	// Handle <edit-config>
	if rpc.EditConfig != nil {
		editCmds, err := handleEditConfig(rpc.EditConfig)
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, editCmds...)
	}

	// Handle <delete-config>
	if rpc.DeleteConfig != nil {
		cmds = append(cmds, handleDeleteConfig(rpc.DeleteConfig)...)
	}

	if len(cmds) == 0 {
		return nil, errors.New("no supported NETCONF operations found")
	}
	return cmds, nil
}

// --- Handlers for NETCONF operations ---

func handleGet(get *Get) []string {
	// Delegate to openconfig system get handler
	return openconfig.SystemGetToMikrotikCmds(get.Filter.Value)
}

func handleEditConfig(edit *EditConfig) ([]string, error) {
	var cmds []string
	// Delegate to openconfig system set handler (registry-based)
	cmds = append(cmds, openconfig.SystemToMikrotikCmdsRegistry("set", edit.Config.System)...)
	// Extend for more OpenConfig modules
	if len(cmds) == 0 {
		return nil, errors.New("no supported edit-config elements found")
	}
	return cmds, nil
}

func handleDeleteConfig(del *DeleteConfig) []string {
	// Example: delete all addresses (very basic)
	return []string{"/ip/address/remove [find]"}
}

// --- Utility functions ---
func parseBool(val string) bool {
	return val == "true" || val == "yes" || val == "1"
}

func mapBoolToMikrotik(b bool) string {
	if b {
		return "no" // MikroTik uses 'disabled=no' to enable
	}
	return "yes" // 'disabled=yes' to disable
}

// --- Main CLI ---
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: translator '<netconf-xml>'")
		os.Exit(1)
	}
	input := os.Args[1]
	cmds, err := TranslateNetconfToMikrotik(input)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	for _, cmd := range cmds {
		fmt.Println(cmd)
	}
}
