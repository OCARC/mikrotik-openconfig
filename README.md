# MikroTik OpenConfig Translator

This project translates OpenConfig/NETCONF XML into MikroTik RouterOS API commands and can apply them to real hardware using the go-routeros library.

## Features
- Modular translation of OpenConfig trees (system, interfaces, ip, etc.)
- Schema validation stub for future YANG/XSD enforcement
- Unit tests for all translation logic
- MikroTik API integration via go-routeros

## Supported OpenConfig System Features
- Hostname
- Clock (timezone)
- NTP (servers, enable/disable)
- DNS (servers)
- (Extendable: AAA, Logging, etc.)

## Directory Structure
- `openconfig/` — OpenConfig tree models and translation logic
- `translator.go` — Main translation entry point
- `client.go` — MikroTik API client wrapper
- `*_test.go` — Unit tests

## Keeping Documentation Up-to-Date
1. **Every new feature or syntax mapping must be documented in this README under "Supported Features".**
2. **Each translation function should have a Go doc comment describing what it maps.**
3. **When adding a new OpenConfig subtree or MikroTik mapping:**
   - Update the README's Supported Features section.
   - Add/expand doc comments in the relevant Go files.
   - Add or update unit tests to cover the new feature.
4. **Pull requests or commits should be blocked if documentation or tests are missing for new features.**

## Example Usage
```
go run translator.go '<netconf-xml>'
```

## Contributing
- Add new OpenConfig features by creating a new file in `openconfig/` and updating the main translation logic.
- Ensure all new features are covered by unit tests and documented in this README.

---

For more details, see the code comments and test cases.
