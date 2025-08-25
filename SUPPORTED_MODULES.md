# Supported OpenConfig Modules and Features

This document tracks which OpenConfig modules and features are currently supported by the MikroTik OpenConfig Translator.

## Supported Features

| IETF Standard | OpenConfig Path | Get | Set | Merge | Test File |
|---------------|-----------------|-----|-----|-------|-----------|
| [RFC 7317](https://datatracker.ietf.org/doc/html/rfc7317) | `openconfig-system:system/hostname` | ✅ | ✅ | ✅ | [set_hostname.xml](netconf-tests/set_hostname.xml) |
| [RFC 7317](https://datatracker.ietf.org/doc/html/rfc7317) | `openconfig-system:system/clock/timezone-name` | ✅ | ✅ | ✅ | [set_timezone.xml](netconf-tests/set_timezone.xml) |
| [RFC 7317](https://datatracker.ietf.org/doc/html/rfc7317) | `openconfig-system:system/clock/timezone-utc-offset` | ✅ | ❌ | ❌ | [get_clock_info.xml](netconf-tests/get_clock_info.xml) |
| [RFC 5905](https://datatracker.ietf.org/doc/html/rfc5905) | `openconfig-system:system/ntp/enabled` | ✅ | ✅ | ✅ | [enable_ntp.xml](netconf-tests/enable_ntp.xml) |
| [RFC 5905](https://datatracker.ietf.org/doc/html/rfc5905) | `openconfig-system:system/ntp/servers/server/address` | ✅ | ✅ | ✅ | [add_ntp_servers.xml](netconf-tests/add_ntp_servers.xml) |
| [RFC 1035](https://datatracker.ietf.org/doc/html/rfc1035) | `openconfig-system:system/dns/servers/server` | ❌ | ✅ | ✅ | *No test yet* |

## Conditional Operations

| Feature | Description | Test File |
|---------|-------------|-----------|
| Hostname Conditional Set | Set hostname only if current hostname is not "testing" | [netconf_set_hostname_if_not_testing.xml](netconf-tests/netconf_set_hostname_if_not_testing.xml) |

## Planned/Not Yet Supported

| IETF Standard | OpenConfig Module | Priority | Notes |
|---------------|-------------------|----------|-------|
| [RFC 7223](https://datatracker.ietf.org/doc/html/rfc7223) | `openconfig-interfaces` | High | Interface configuration and state |
| [RFC 791](https://datatracker.ietf.org/doc/html/rfc791) | `openconfig-if-ip` | High | IPv4/IPv6 address configuration |
| [RFC 3411](https://datatracker.ietf.org/doc/html/rfc3411) | `openconfig-aaa` | Medium | Authentication, Authorization, Accounting |
| [RFC 3164](https://datatracker.ietf.org/doc/html/rfc3164) | `openconfig-logging` | Medium | System logging configuration |
| [RFC 4251](https://datatracker.ietf.org/doc/html/rfc4251) | `openconfig-system:ssh` | Low | SSH server configuration |
| [RFC 854](https://datatracker.ietf.org/doc/html/rfc854) | `openconfig-system:telnet` | Low | Telnet server configuration |

## Legend

- ✅ **Supported**: Feature is implemented and tested
- ❌ **Not Supported**: Feature is recognized but not implemented
- **Get**: Retrieve current configuration/state
- **Set**: Modify configuration 
- **Merge**: Merge configuration changes (NETCONF merge operation)

## Implementation Status

### Fully Supported (Get + Set + Merge)
- System hostname
- System timezone (timezone-name)
- NTP client enabled/disabled
- NTP server addresses

### Partially Supported (Get only)
- System timezone (timezone-utc-offset) - read-only due to MikroTik limitations

### Set Only
- DNS servers - verification not yet implemented

## Testing Coverage

All supported features have corresponding E2E tests in the `netconf-tests/` directory that:
- ✅ Execute the NETCONF operation against a live MikroTik device
- ✅ Verify the configuration change took effect
- ✅ Support both single operations and conditional operations

## Contributing

When adding new OpenConfig features:

1. **Update this table** with the new feature details
2. **Add IETF RFC references** for the underlying protocol/standard
3. **Create test files** in `netconf-tests/` directory
4. **Implement translator logic** in the `openconfig/` package
5. **Verify E2E tests pass** against a live MikroTik device

---

**Last Updated:** $(date)  
**Total Supported Features:** 6  
**Test Coverage:** 83% (5/6 features have tests)
