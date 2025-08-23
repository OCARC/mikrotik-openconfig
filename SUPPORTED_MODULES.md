# Supported OpenConfig Modules and Features

This document tracks which OpenConfig modules and features are currently supported by the MikroTik OpenConfig Translator.

## Supported Modules

### openconfig-system
- system/hostname
- system/clock/timezone-name
- system/ntp/enabled
- system/ntp/servers/server/address
- system/dns/servers/server

## Not Yet Supported / Planned
- openconfig-interfaces
- openconfig-if-ip
- openconfig-aaa
- openconfig-logging
- (and others)

---

**Process:**
- Update this file whenever a new OpenConfig module or feature is added or removed from the codebase.
- Keep this file in sync with the README and code comments.
- Use this as a checklist for future development and coverage.
