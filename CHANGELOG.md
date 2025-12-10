# Changelog

All notable changes to the Terraform Provider for Keeper Secrets Manager will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.1.8] - 2025-12-08

### Security
- Upgrade Go from 1.24.0 to 1.24.8 to address critical vulnerabilities:
  - **CVE-2025-22871**: net/http chunked encoding request smuggling vulnerability
  - **CVE-2025-58185**: DER payload parsing memory exhaustion vulnerability
- Update GitHub Actions workflows to use Go 1.24.8 for builds and SBOM generation

### Added
- **PAM Record Type Support** (KSM-527):
  - Add `secretsmanager_pam_machine` resource and data source for SSH, RDP, and remote machine credentials
  - Add `secretsmanager_pam_database` resource and data source for PostgreSQL, MySQL, MongoDB, and database credentials
  - Add `secretsmanager_pam_directory` resource and data source for Active Directory and LDAP credentials
  - Enhanced `secretsmanager_pam_user` data source with `private_pem_key` field support
  - Add `pamSettings` field for protocol-specific connection configuration as JSON
  - Add schema functions in `record_fields_pam.go` for PAM-specific fields
  - Add 16 new acceptance tests validating complete CRUD lifecycle for PAM types
  - Add 6 comprehensive example files demonstrating PAM resource and data source usage

- **Regex Pattern Support** (KSM-389):
  - Add `title_patterns` parameter to `secretsmanager_records` data source for filtering with Go regex
  - Support multiple patterns in a single query
  - Combine with existing UIDs and exact title filters
  - Add 4 new acceptance tests for pattern matching functionality
  - Update documentation with regex pattern examples and performance considerations

- Add GitHub Actions workflow for automated testing on pull requests
- Add explicit `contents: read` permissions to test workflow for security compliance

### Fixed
- Fix shortcuts/linked records error (KSM-522) - resolve duplicate UID handling across multiple shared folders
- Fix "changes to folder_uid not allowed" errors during Terraform apply operations
- Use `reflect.DeepEqual` for JSON comparison to handle map ordering correctly instead of string comparison
- Fix PAM field labels to match backend schema (useSSL, connectDatabase)
- Fix field access patterns in PAM data source examples
- Update PAM data source examples to use UIDs instead of fake folder paths
- Fix test helpers in `data_source_records_test.go` (ProviderFactories → Providers)
- Use `t.Skip()` instead of `t.Fatal()` for missing test setup to prevent CI failures

### Changed
- Update version references from 1.1.7 to 1.1.8 across all 44 example files
- Clarify `folder_uid` description to reflect sub-folder support (parent shared folder access sufficient)
- Clarify checkbox field comment to explain Keeper stores values as single-element arrays
- Add Go version compatibility flag (`-compat=1.24.8`) to goreleaser config

### Removed
- Remove `login` and `password` fields from `secretsmanager_pam_database`, `secretsmanager_pam_directory`, and `secretsmanager_pam_machine` (not in official templates; only in `pamUser`)
- Remove invalid `connect_database` field from `pamDatabase`
- Remove invalid `ssl_verification` field from `pamMachine`
- Remove login/password from PAM test cases and resource examples

## [1.1.7] - 2025-11-20

### Security
- Bump golang.org/x/crypto from 0.42.0 to 0.45.0 in the go_modules group

### Fixed
- Fix folder UID validation and empty folder restriction in resource schema descriptions

## [1.1.6] - 2025-10-15

_(Previous releases not documented - will be added retroactively)_

[Unreleased]: https://github.com/Keeper-Security/terraform-provider-secretsmanager/compare/v1.1.8...HEAD
[1.1.8]: https://github.com/Keeper-Security/terraform-provider-secretsmanager/compare/v1.1.7...v1.1.8
[1.1.7]: https://github.com/Keeper-Security/terraform-provider-secretsmanager/compare/v1.1.6...v1.1.7
[1.1.6]: https://github.com/Keeper-Security/terraform-provider-secretsmanager/releases/tag/v1.1.6
