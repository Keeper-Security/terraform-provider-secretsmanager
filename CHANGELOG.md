# Changelog

All notable changes to the Terraform Provider for Keeper Secrets Manager will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Custom Fields** (KSM-388):
  - Add `custom` block to all 22 resource types (`login`, `bank_account`, `bank_card`, `birth_certificate`, `contact`, `database_credentials`, `driver_license`, `encrypted_notes`, `file`, `health_insurance`, `membership`, `passport`, `photo`, `server_credentials`, `software_license`, `ssh_keys`, `ssn_card`, `address`, `pam_database`, `pam_directory`, `pam_machine`, `pam_remote_browser`, `pam_user`)
  - Supports 43+ Keeper field types including `text`, `secret`, `url`, `email`, `phone`, `date`, `birthDate`, `expirationDate`, `name`, `address`, `paymentCard`, `bankAccount`, `host`, `keyPair`, `securityQuestion`, `checkbox`, `multiline`, and more
  - Simple types use a plain string `value`; complex types use `value = jsonencode({...})` for a single entry or `value = jsonencode([{...},{...}])` for multiple entries in one field
  - `pam_machine` and `pam_user` use merge-aware logic to preserve the vault-managed "Private Key Passphrase" custom field across create/update operations
  - `required` and `privacy_screen` attributes round-trip correctly from vault state (no perpetual diff on import)

### Fixed
- **Custom fields — `paymentCard` perpetual diff** (KSM-888): `jsonencode()` values must use camelCase keys — `cardNumber`, `cardExpirationDate`, `cardSecurityCode` — matching Keeper's API format. Snake_case keys (`card_number`, etc.) were previously silently ignored, causing the field to be written empty and producing a perpetual plan diff.
- **Custom fields — non-canonical `checkbox` values** (KSM-889): only `"true"` or `"false"` are accepted. Other strings like `"yes"` or `"1"` now return a clear error instead of being silently coerced to `false`.
- **Custom fields — non-canonical date values** (KSM-889): `date`, `birthDate`, and `expirationDate` only accept YYYY-MM-DD format. RFC3339 input (e.g. `"2026-03-20T14:30:00Z"`) now returns a clear error instead of causing a perpetual plan diff (config kept RFC3339; state returned YYYY-MM-DD).

## [1.3.0]

### Security
- Bump `cloudflare/circl` to v1.6.3 and `grpc` to v1.79.3 to address known vulnerabilities

### Added
- Document Alpine Linux and musl-based container compatibility — all Linux binaries are statically compiled (`CGO_ENABLED=0`) with no C library dependencies and run on Alpine and other musl-based systems without modification (KSM-885)

- **Ephemeral Resources** (KSM-871):
  - Add ephemeral resource support for Terraform 1.10+, ensuring secrets are never stored in `terraform.tfstate`
  - Ephemeral resources available for all 25 record types: `login`, `field`, `record`, `database_credentials`, `server_credentials`, `ssh_keys`, `encrypted_notes`, `address`, `bank_account`, `bank_card`, `birth_certificate`, `contact`, `driver_license`, `health_insurance`, `membership`, `passport`, `photo`, `software_license`, `ssn_card`, `file`, `pam_user`, `pam_machine`, `pam_database`, `pam_directory`, `pam_remote_browser`
  - Use `ephemeral "secretsmanager_<type>" "name" { ... }` instead of `data` blocks to keep secrets out of state
  - Introduce Terraform Plugin Framework alongside existing SDKv2 via protocol v6 mux server
  - Add `terraform-plugin-framework` v1.18.0 and `terraform-plugin-mux` v0.22.0 dependencies
  - All existing resources and data sources remain fully backward compatible
  - Add documentation and examples for all ephemeral resources

- **PAM Remote Browser** (KSM-871):
  - Add `secretsmanager_pam_remote_browser` resource, data source, and ephemeral resource
  - Support for Remote Browser Isolation (RBI) URL, browser settings (JSON), traffic encryption seed, file references, and TOTP
  - Full CRUD lifecycle with import support

### Fixed
- Fix wrong JSON key in `pamHostnameToListValue` causing incorrect PAM hostname field mapping (KSM-884)
- Remove invalid `DiffSuppressFunc` and `ValidateFunc` from computed-only `pam_remote_browser_settings` field in the `pam_remote_browser` data source
- Add nil-check guard in all ephemeral resource `Open()` methods to prevent panics if provider configuration is missing
- Surface warning diagnostics when referenced `addressRef` or `cardRef` records cannot be fetched, instead of silently returning empty fields
- Mark `credential` provider attribute as sensitive to prevent credentials appearing in plan output
- Mark sensitive fields across all record types to prevent secrets appearing in plan output: payment card numbers and security codes, bank account and routing numbers, PIN codes, TOTP seeds, license numbers, and secret field values

## [1.2.0]

### Security
- Upgrade Go from 1.24.0 to 1.24.13 to address critical vulnerabilities

### Added
- **SSH Key Generation** (KSM-788):
  - Add automatic SSH key pair generation to `secretsmanager_ssh_keys` resource
  - Support ED25519, RSA (2048/3072/4096), and ECDSA (P-256/P-384/P-521) key types
  - Generate SSH keys via `generate = "yes"` on `key_pair` block
  - Automatic private key encryption with passphrase using OpenSSH bcrypt+aes256-ctr format
  - Add unit tests for all key types and acceptance tests for ED25519 and RSA+passphrase
  - Update documentation with generation examples and key type options

- **PAM SSH Key Generation** (KSM-789):
  - Add SSH key generation support to PAM User and Machine resources
  - Add `private_pem_key` field (standard secret) to `secretsmanager_pam_user` and `secretsmanager_pam_machine`
  - Add `private_key_passphrase` field (custom secret) for encrypted key storage
  - Generate SSH keys via `generate = "yes"` on `private_pem_key` block
  - Support same key types as ssh_keys resource (ED25519, RSA, ECDSA)
  - Passphrase stored as custom field for kdnrm PAM interoperability
  - Add comprehensive documentation for pam_user and pam_machine SSH key generation

- **PAM Record Type Support** (KSM-527):
  - Add `secretsmanager_pam_machine` resource and data source for SSH, RDP, and remote machine credentials
  - Add `secretsmanager_pam_database` resource and data source for PostgreSQL, MySQL, MongoDB, and database credentials
  - Add `secretsmanager_pam_directory` resource and data source for Active Directory and LDAP credentials
  - Enhanced `secretsmanager_pam_user` data source with `private_pem_key` field support
  - Add `pamSettings` field for protocol-specific connection configuration as JSON
  - PAM-specific fields use flat value syntax consistent with standard fields: `database_type = "postgresql"`, `directory_type = "Active Directory"`, `use_ssl { value = true }`
  - Add schema functions in `record_fields_pam.go` for PAM-specific fields
  - Add 21 acceptance tests covering full CRUD lifecycle, data source field readback, and auto-generated UID for all 4 PAM types
  - Add data source documentation for all 4 PAM types (`pam_machine`, `pam_database`, `pam_directory`, `pam_user`)
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
- Fix PAM `pam_settings` field readback for PAM data sources and PAM Machine resource lifecycle handling (KSM-796)
- Fix PAM Machine data source missing `login`, `password`, `private_pem_key`, `private_key_passphrase`, and `ssl_verification` fields (KSM-797)
- Fix PAM data source coverage gaps for `folder_uid`, `totp`, and PAM User `private_key_passphrase` readback (KSM-798)
- Fix PAM User data source returning empty values for `connect_database` and `private_pem_key` fields (KSM-794)
- Fix shortcuts/linked records error (KSM-522) - resolve duplicate UID handling across multiple shared folders
- Fix "changes to folder_uid not allowed" errors during Terraform apply operations
- Use `reflect.DeepEqual` for JSON comparison to handle map ordering correctly instead of string comparison
- Fix test helpers in `data_source_records_test.go` (ProviderFactories → Providers)
- Use `t.Skip()` instead of `t.Fatal()` for missing test setup to prevent CI failures

### Changed
- **Resource Documentation Improvements** (KSM-790):
  - Update 18 resource documentation files to clarify resources "create and manage" secrets (not just "access")
  - Add "Example Usage" sections with working Terraform code samples from examples directory
  - Improve clarity around resource CRUD lifecycle capabilities
  - Resources updated: address, bank_account, bank_card, birth_certificate, contact, database_credentials, driver_license, encrypted_notes, file, health_insurance, login, membership, passport, photo, server_credentials, software_license, ssh_keys, ssn_card
- Clarify `folder_uid` description to reflect sub-folder support (parent shared folder access sufficient)
- Clarify checkbox field comment to explain Keeper stores values as single-element arrays
- Add Go version compatibility flag (`-compat=1.24.13`) to goreleaser config


## [1.1.7]

### Security
- Bump golang.org/x/crypto from 0.42.0 to 0.45.0 in the go_modules group

### Fixed
- Fix folder UID validation and empty folder restriction in resource schema descriptions

[Unreleased]: https://github.com/Keeper-Security/terraform-provider-secretsmanager/compare/v1.3.0...HEAD
[1.3.0]: https://github.com/Keeper-Security/terraform-provider-secretsmanager/compare/v1.2.0...v1.3.0
[1.2.0]: https://github.com/Keeper-Security/terraform-provider-secretsmanager/compare/v1.1.7...v1.2.0
[1.1.7]: https://github.com/Keeper-Security/terraform-provider-secretsmanager/compare/v1.1.6...v1.1.7
[1.1.6]: https://github.com/Keeper-Security/terraform-provider-secretsmanager/releases/tag/v1.1.6
