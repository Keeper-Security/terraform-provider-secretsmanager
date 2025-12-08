# PAM Field Mapping Reference
**Generated:** 2025-12-08
**Purpose:** Complete field-by-field mapping between MCP schema and Terraform implementation

## How to Read This Document

- **MCP Field**: Official field name from Keeper backend (via `get_record_type_schema`)
- **TF Field**: Terraform resource schema field name (snake_case)
- **Type**: Data type (text, secret, checkbox, etc.)
- **Req**: Required field (Y/N)
- **Status**:
  - ✅ Correctly mapped
  - ❌ Missing in Terraform
  - ⚠️ Extra in Terraform (not in MCP schema)
  - 🐛 Bug in implementation

---

## 1. pamDatabase Field Mapping

| # | MCP Field | TF Field | Type | Req | Status | Notes |
|---|-----------|----------|------|-----|--------|-------|
| 1 | pamHostname.hostName | pam_hostname[0].value[0].hostname | string | N | ✅ | |
| 2 | pamHostname.port | pam_hostname[0].value[0].port | string | N | ✅ | |
| 3 | useSSL | use_ssl[0].value[0] | checkbox | N | ✅ | Label: "useSSL" |
| 4 | pamSettings | pam_settings | pamSettings | N | ✅ | JSON string |
| 5 | trafficEncryptionSeed | - | trafficEncryptionSeed | N | - | System-managed (correct) |
| 6 | rotationScripts.command | rotation_scripts[0].value[0].command | string | N | ✅ | |
| 7 | rotationScripts.fileRef | rotation_scripts[0].value[0].file_ref | string | N | ✅ | |
| 8 | rotationScripts.recordRef | rotation_scripts[0].value[0].record_ref | string | N | ✅ | |
| 9 | databaseId | database_id | text | N | ✅ | |
| 10 | databaseType | database_type | dropdown | N | ✅ | Values: mysql, postgresql, etc. |
| 11 | providerGroup | provider_group | text | N | ✅ | |
| 12 | providerRegion | provider_region | text | N | ✅ | |
| 13 | fileRef | file_ref | fileRef | N | ✅ | |
| 14 | oneTimeCode | totp | otp | N | ✅ | |
| - | - | login | login | N | ⚠️ | Extra (valid for DB credentials) |
| - | - | password | password | N | ⚠️ | Extra (valid for DB credentials) |
| - | - | connect_database | text | N | ⚠️ | Extra (may be legacy field) |

**Summary:**
- **MCP Fields**: 14
- **TF Fields**: 17 (13 mapped + 3 extra + system field excluded)
- **Missing**: 0
- **Extra**: 3 (login, password, connect_database)
- **Bugs**: 0
- **Verdict**: ✅ Complete (extras are acceptable)

---

## 2. pamMachine Field Mapping

| # | MCP Field | TF Field | Type | Req | Status | Notes |
|---|-----------|----------|------|-----|--------|-------|
| 1 | pamHostname.hostName | pam_hostname[0].value[0].hostname | string | **Y** | ✅ | **Required** |
| 2 | pamHostname.port | pam_hostname[0].value[0].port | string | **Y** | ✅ | **Required** |
| 3 | pamSettings | pam_settings | pamSettings | N | ✅ | JSON string |
| 4 | trafficEncryptionSeed | - | trafficEncryptionSeed | N | - | System-managed (correct) |
| 5 | rotationScripts.command | rotation_scripts[0].value[0].command | string | N | ✅ | |
| 6 | rotationScripts.fileRef | rotation_scripts[0].value[0].file_ref | string | N | ✅ | |
| 7 | rotationScripts.recordRef | rotation_scripts[0].value[0].record_ref | string | N | ✅ | |
| 8 | operatingSystem | operating_system | text | N | ✅ | |
| 9 | instanceName | instance_name | text | N | ✅ | |
| 10 | instanceId | instance_id | text | N | ✅ | |
| 11 | providerGroup | provider_group | text | N | ✅ | |
| 12 | providerRegion | provider_region | text | N | ✅ | |
| 13 | fileRef | file_ref | fileRef | N | ✅ | |
| 14 | oneTimeCode | totp | otp | N | ✅ | |
| - | - | login | login | N | ⚠️ | Extra (valid for machine login) |
| - | - | password | password | N | ⚠️ | Extra (valid for machine login) |
| - | **NOT IN MCP** | ssl_verification | checkbox | N | 🐛 | **INVALID - Remove or verify** |

**Summary:**
- **MCP Fields**: 14
- **TF Fields**: 17 (13 mapped + 2 extra + 1 invalid + system field excluded)
- **Missing**: 0
- **Extra**: 2 valid (login, password)
- **Invalid**: 1 (ssl_verification) 🔴 **ACTION REQUIRED**
- **Verdict**: ❌ Has invalid field

**Critical Issue:**
```go
// resource_pam_machine.go:62
"ssl_verification": schemaCheckboxField(),  // ❌ NOT IN MCP SCHEMA
```

**Investigation needed:**
- Check Go SDK for `ssl_verification` field definition
- Search codebase for usage in create/read/update functions
- Verify with backend team or remove

---

## 3. pamUser Field Mapping

| # | MCP Field | TF Field | Type | Req | Status | Notes |
|---|-----------|----------|------|-----|--------|-------|
| 1 | login | login | login | **Y** | ✅ | **Required** |
| 2 | password | password | password | N | ✅ | |
| 3 | rotationScripts.command | rotation_scripts[0].value[0].command | string | N | ✅ | |
| 4 | rotationScripts.fileRef | rotation_scripts[0].value[0].file_ref | string | N | ✅ | |
| 5 | rotationScripts.recordRef | rotation_scripts[0].value[0].record_ref | string | N | ✅ | |
| 6 | **privatePEMKey** | - | secret | N | ❌ | **MISSING - Add to TF** |
| 7 | distinguishedName | distinguished_name | text | N | ✅ | |
| 8 | connectDatabase | connect_database | text | N | ✅ | |
| 9 | managed | managed | checkbox | N | ✅ | |
| 10 | fileRef | file_ref | fileRef | N | ✅ | |
| 11 | oneTimeCode | totp | otp | N | ✅ | |

**Summary:**
- **MCP Fields**: 11
- **TF Fields**: 10 (10 mapped + 1 missing)
- **Missing**: 1 (privatePEMKey) 🟡 **MEDIUM PRIORITY**
- **Extra**: 0
- **Bugs**: 0
- **Verdict**: ❌ Missing SSH key field

**Missing Field:**
```go
// Add to resource_pam_user.go schema (after line 57):
"private_pem_key": schemaSecretField(),  // maps to privatePEMKey
```

**Implementation Notes:**
- Type: `secret` (masked field)
- Required: No
- Use case: SSH private key authentication
- Similar to: password field (sensitive)

---

## 4. pamDirectory Field Mapping

| # | MCP Field | TF Field | Type | Req | Status | Notes |
|---|-----------|----------|------|-----|--------|-------|
| 1 | pamHostname.hostName | pam_hostname[0].value[0].hostname | string | N | ✅ | |
| 2 | pamHostname.port | pam_hostname[0].value[0].port | string | N | ✅ | |
| 3 | **useSSL** | use_ssl | checkbox | N | 🐛 | **Label bug: "Use SSL" should be "useSSL"** |
| 4 | pamSettings | pam_settings | pamSettings | N | ✅ | JSON string |
| 5 | trafficEncryptionSeed | - | trafficEncryptionSeed | N | - | System-managed (correct) |
| 6 | rotationScripts.command | rotation_scripts[0].value[0].command | string | N | ✅ | |
| 7 | rotationScripts.fileRef | rotation_scripts[0].value[0].file_ref | string | N | ✅ | |
| 8 | rotationScripts.recordRef | rotation_scripts[0].value[0].record_ref | string | N | ✅ | |
| 9 | **domainName** | - | text | N | ❌ | **MISSING** |
| 10 | **alternativeIPs** | - | unknown | N | ❌ | **MISSING** (MCP shows error, likely multiline) |
| 11 | **directoryId** | - | text | N | ❌ | **MISSING** |
| 12 | directoryType | directory_type | dropdown | N | ✅ | Values: Active Directory, OpenLDAP |
| 13 | **userMatch** | - | text | N | ❌ | **MISSING** |
| 14 | **providerGroup** | - | text | N | ❌ | **MISSING** |
| 15 | **providerRegion** | - | text | N | ❌ | **MISSING** |
| 16 | fileRef | file_ref | fileRef | N | ✅ | |
| 17 | oneTimeCode | totp | otp | N | ✅ | |
| - | - | login | login | N | ⚠️ | Extra (valid for bind user) |
| - | - | password | password | N | ⚠️ | Extra (valid for bind user) |
| - | - | distinguished_name | text | N | ⚠️ | Extra (not in MCP, but valid for LDAP bind DN) |

**Summary:**
- **MCP Fields**: 19 (including trafficEncryptionSeed)
- **TF Fields**: 13 (10 mapped + 3 extra + 6 missing + system field excluded)
- **Missing**: 6 🟡 **MEDIUM PRIORITY**
- **Extra**: 3 (login, password, distinguished_name - all valid)
- **Bugs**: 1 (useSSL label) 🔴 **CRITICAL**
- **Verdict**: ❌ Critical bug + significant gaps

**Critical Bug (IMMEDIATE FIX):**
```go
// resource_pam_directory.go:170
field.(*core.Checkbox).Label = "Use SSL"  // ❌ WRONG

// Should be:
field.(*core.Checkbox).Label = "useSSL"   // ✅ CORRECT (camelCase)

// resource_pam_directory.go:335
useSSL := getFieldResourceDataWithLabel("checkbox", "fields", secret, "Use SSL")  // ❌ WRONG

// Should be:
useSSL := getFieldResourceDataWithLabel("checkbox", "fields", secret, "useSSL")   // ✅ CORRECT
```

**Missing Fields (Add to Schema):**
```go
// Add to resource_pam_directory.go after line 63:
"domain_name":    schemaTextField(),       // maps to domainName
"alternative_ips": schemaMultilineField(), // maps to alternativeIPs (wait for MCP fix)
"directory_id":   schemaTextField(),       // maps to directoryId
"user_match":     schemaTextField(),       // maps to userMatch
"provider_group": schemaTextField(),       // maps to providerGroup
"provider_region": schemaTextField(),      // maps to providerRegion
```

**Note on alternativeIPs:**
- MCP schema shows: `"type": "unknown"` with error message
- JSON template shows: `{"$ref": "multiline", "label": "alternativeIPs"}`
- Likely backend schema bug - report to Keeper team
- Implementation: Use multiline type (per template) but document uncertainty

---

## Field Type Reference

### Mapping Convention: MCP Type → Terraform Schema Function

| MCP Type | TF Schema Function | Example Field | Notes |
|----------|-------------------|---------------|-------|
| `string` | (inline struct) | pamHostname.hostName | Part of complex field |
| `text` | `schemaTextField()` | databaseId | Plain text |
| `secret` | `schemaSecretField()` | privatePEMKey | Masked/sensitive |
| `password` | `schemaPasswordField()` | password | Password with generation |
| `login` | `schemaLoginField()` | login | Username/login |
| `checkbox` | `schemaCheckboxField()` | useSSL | Boolean checkbox |
| `dropdown` | `schema*TypeField()` | databaseType | Enum (specific function per type) |
| `otp` | `schemaOneTimeCodeField()` | oneTimeCode | TOTP/2FA |
| `fileRef` | `schemaFileRefField()` | fileRef | File reference |
| `multiline` | `schemaMultilineField()` | alternativeIPs | Multi-line text |
| `pamHostname` | `schemaPamHostnameField()` | pamHostname | Complex: hostname + port |
| `pamSettings` | `schemaPamSettingsField()` | pamSettings | Complex JSON |
| `script` | `schemaScriptField()` | rotationScripts | Script execution details |
| `trafficEncryptionSeed` | - | - | System-managed, not exposed |

### Label Conventions

**MCP Field Name → Terraform Field Name:**
```
camelCase     → snake_case
useSSL        → use_ssl
privatePEMKey → private_pem_key
domainName    → domain_name
```

**When setting Label in code:**
```go
// ✅ CORRECT: Use MCP field name (camelCase)
field.(*core.Checkbox).Label = "useSSL"

// ❌ WRONG: Don't use human-readable format
field.(*core.Checkbox).Label = "Use SSL"
```

---

## Implementation Priorities

### 🔴 CRITICAL (Fix Immediately)
1. **pamDirectory: Fix useSSL label bug**
   - File: `resource_pam_directory.go`
   - Lines: 170, 335
   - Change: `"Use SSL"` → `"useSSL"`
   - Impact: Breaks create/read operations

### 🟡 HIGH (Should Fix)
2. **pamMachine: Investigate ssl_verification field**
   - File: `resource_pam_machine.go`
   - Line: 62
   - Action: Check SDK + backend team, remove if invalid

### 🟢 MEDIUM (Feature Additions)
3. **pamUser: Add privatePEMKey field**
   - File: `resource_pam_user.go`
   - Add after line 57: `"private_pem_key": schemaSecretField()`
   - Update Create/Read/Update functions

4. **pamDirectory: Add 6 missing fields**
   - File: `resource_pam_directory.go`
   - Add: domain_name, directory_id, user_match, provider_group, provider_region
   - Hold: alternative_ips (wait for MCP schema fix)

---

## Testing Checklist

After implementing changes, verify each mapping:

```bash
# Set test environment
export TF_ACC=1
export KEEPER_CREDENTIAL=<xxx>

# Test pamDatabase (should pass - no changes)
go test ./secretsmanager -v -run TestAccResourcePamDatabase

# Test pamMachine (after ssl_verification fix)
go test ./secretsmanager -v -run TestAccResourcePamMachine

# Test pamUser (after adding privatePEMKey)
go test ./secretsmanager -v -run TestAccResourcePamUser

# Test pamDirectory (after fixing useSSL label)
go test ./secretsmanager -v -run TestAccResourcePamDirectory

# Full PAM suite
go test ./secretsmanager -v -run TestAccResourcePam
go test ./secretsmanager -v -run TestAccDataSourcePam
```

---

## Change Log

| Date | Change | Type | Files Modified |
|------|--------|------|----------------|
| 2025-12-08 | Initial mapping created | Documentation | pam-field-mapping.md |

---

## See Also

- [PAM Schema Validation Report](./pam-schema-validation-report.md) - Detailed analysis
- [PAM Fields Verification Report](../plans/pam-fields-verification-2025-12-08.md) - Original investigation
- MCP Server: `get_record_type_schema` tool documentation
