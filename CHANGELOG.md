# Changelog

All notable changes to the Terraform Provider for Keeper Secrets Manager will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.1.8] - 2025-12-08

### Security
- Upgrade Go from 1.24.0 to 1.24.8 to address critical vulnerabilities
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
  - Add ReDoS protection with 500-character pattern length limit
  - Add 5 new acceptance tests for pattern matching functionality (including length validation)
  - Update documentation with regex pattern examples, performance warnings, and security considerations

- Add GitHub Actions workflow for automated testing on pull requests
- Add explicit `contents: read` permissions to test workflow for security compliance

### Fixed
- Fix shortcuts/linked records error (KSM-522) - resolve duplicate UID handling across multiple shared folders
- Fix "changes to folder_uid not allowed" errors during Terraform apply operations
- Use `reflect.DeepEqual` for JSON comparison to handle map ordering correctly instead of string comparison
- Fix PAM field labels to match backend schema (useSSL, connectDatabase)
- Fix field access patterns in PAM data source examples
- Fix test helpers in `data_source_records_test.go` (ProviderFactories → Providers)
- Use `t.Skip()` instead of `t.Fatal()` for missing test setup to prevent CI failures

### Changed
- Clarify `folder_uid` description to reflect sub-folder support (parent shared folder access sufficient)
- Clarify checkbox field comment to explain Keeper stores values as single-element arrays
- Add Go version compatibility flag (`-compat=1.24.8`) to goreleaser config


## [1.1.7] - 2025-11-20

### Security
- Bump golang.org/x/crypto from 0.42.0 to 0.45.0 in the go_modules group

### Fixed
- Fix folder UID validation and empty folder restriction in resource schema descriptions

[Unreleased]: https://github.com/Keeper-Security/terraform-provider-secretsmanager/compare/v1.1.8...HEAD
[1.1.8]: https://github.com/Keeper-Security/terraform-provider-secretsmanager/compare/v1.1.7...v1.1.8
[1.1.7]: https://github.com/Keeper-Security/terraform-provider-secretsmanager/compare/v1.1.6...v1.1.7
[1.1.6]: https://github.com/Keeper-Security/terraform-provider-secretsmanager/releases/tag/v1.1.6
