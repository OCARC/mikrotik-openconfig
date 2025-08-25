# DEVELOPMENT.md

## File Naming Conventions

- Use lowercase letters and underscores (`_`) to separate words in filenames.
- Group related features or modules in subdirectories (e.g., `openconfig/`).
- For OpenConfig feature handlers, use the pattern: `system_get.go`, `system_set.go`, `system.go`.
- Test files should use the `_test.go` suffix and match the feature/module they test (e.g., `translator_system_test.go`).
- Avoid redundant or legacy files; remove obsolete files as the codebase evolves.

## Method Naming Conventions

- Use `CamelCase` for exported (public) functions and methods (e.g., `TranslateNetconfToMikrotik`).
- Use `camelCase` for unexported (private) functions and methods (e.g., `handleSystemHostname`).
- Prefix handler functions with the feature or operation they handle (e.g., `handleSystemHostname`, `SystemGetToMikroTikCmds`).
- For registry or dispatcher functions, use descriptive names indicating their role (e.g., `SystemToMikroTikCmdsRegistry`).
- Test functions should start with `Test` followed by the function or feature under test (e.g., `TestTranslateNetconfToMikrotik_SystemHostname`).

## General Guidelines

- Keep file and method names descriptive and concise.
- Follow Go community conventions for naming and organization.
- Update this document as conventions evolve with the project.
