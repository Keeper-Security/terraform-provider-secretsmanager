# PAM Schema Validation Report
**Generated:** 2025-12-08
**Method:** MCP `get_record_type_schema` vs Terraform Resource Schemas
**Source of Truth:** Keeper Backend via MCP Server

## Executive Summary

This report compares the official Keeper backend schemas (retrieved via MCP server) against the Terraform provider implementation for all PAM record types.

**Status:** ⚠️ **DISCREPANCIES FOUND**
- **pamDatabase**: ✅ Complete (all MCP fields implemented)
- **pamMachine**: ❌ Has extra field `ssl_verification` not in MCP schema
- **pamUser**: ❌ Missing `privatePEMKey` field
- **pamDirectory**: ❌ Missing 6 fields + has `use_ssl` label bug

---

## Validation Methodology

```
┌─────────────────────┐
│   Keeper Backend    │
│   (Source of Truth) │
└──────────┬──────────┘
           │
           │ MCP Protocol
           ▼
┌─────────────────────┐       ┌──────────────────────┐
│   MCP Server Tool   │◄─────►│  Terraform Provider  │
│ get_record_type_    │       │  resource_pam_*.go   │
│     schema()        │       └──────────────────────┘
└─────────────────────┘
           │
           │ Comparison
           ▼
┌─────────────────────┐
│  Validation Report  │
│   (This Document)   │
└─────────────────────┘
```

**MCP Schema Fields Structure:**
```json
{
  "record_type": "pamDatabase",
  "fields": [
    {
      "name": "pamHostname.hostName",
      "description": "...",
      "type": "string",
      "required": false
    }
  ]
}
```

---

## 1. pamDatabase Schema Validation

### MCP Schema (Official - 14 fields)
```
✓ pamHostname.hostName      (string, optional)
✓ pamHostname.port          (string, optional)
✓ useSSL                    (checkbox, optional)
✓ pamSettings               (pamSettings, optional)
✓ trafficEncryptionSeed     (trafficEncryptionSeed, optional)
✓ rotationScripts.command   (string, optional)
✓ rotationScripts.fileRef   (string, optional)
✓ rotationScripts.recordRef (string, optional)
✓ databaseId                (text, optional)
✓ databaseType              (dropdown, optional)
✓ providerGroup             (text, optional)
✓ providerRegion            (text, optional)
✓ fileRef                   (fileRef, optional)
✓ oneTimeCode               (otp, optional)
```

### Terraform Implementation (Lines 24-69)
```go
Schema: map[string]*schema.Schema{
    "pam_hostname":     schemaPamHostnameField(),      // ✅ maps to pamHostname.*
    "use_ssl":          schemaCheckboxField(),         // ✅ maps to useSSL
    "pam_settings":     schemaPamSettingsField(),      // ✅ maps to pamSettings
    "login":            schemaLoginField(),            // ⚠️ NOT in MCP (but valid for DB users)
    "password":         schemaPasswordField(""),       // ⚠️ NOT in MCP (but valid for DB users)
    "rotation_scripts": schemaScriptField(),           // ✅ maps to rotationScripts.*
    "connect_database": schemaTextField(),             // ⚠️ NOT in MCP
    "database_id":      schemaTextField(),             // ✅ maps to databaseId
    "database_type":    schemaDatabaseTypeField(),     // ✅ maps to databaseType
    "provider_group":   schemaTextField(),             // ✅ maps to providerGroup
    "provider_region":  schemaTextField(),             // ✅ maps to providerRegion
    "file_ref":         schemaFileRefField(),          // ✅ maps to fileRef
    "totp":             schemaOneTimeCodeField(),      // ✅ maps to oneTimeCode
}
```

### Analysis
**Status:** ✅ **ACCEPTABLE** (all MCP fields present, extras are valid)

**Extra fields in Terraform:**
- `login`, `password` - Valid for database user credentials
- `connect_database` - May be legacy/undocumented field
- Note: `trafficEncryptionSeed` NOT exposed in Terraform (system-managed, correct)

**Verdict:** No action needed. Schema is complete.

---

## 2. pamMachine Schema Validation

### MCP Schema (Official - 14 fields)
```
✓ pamHostname.hostName      (string, required)
✓ pamHostname.port          (string, required)
✓ pamSettings               (pamSettings, optional)
✓ trafficEncryptionSeed     (trafficEncryptionSeed, optional)
✓ rotationScripts.command   (string, optional)
✓ rotationScripts.fileRef   (string, optional)
✓ rotationScripts.recordRef (string, optional)
✓ operatingSystem           (text, optional)
✓ instanceName              (text, optional)
✓ instanceId                (text, optional)
✓ providerGroup             (text, optional)
✓ providerRegion            (text, optional)
✓ fileRef                   (fileRef, optional)
✓ oneTimeCode               (otp, optional)
✗ ssl_verification          ❌ NOT IN MCP SCHEMA
```

### Terraform Implementation (Lines 24-69)
```go
Schema: map[string]*schema.Schema{
    "pam_hostname":     schemaPamHostnameField(),      // ✅ maps to pamHostname.*
    "pam_settings":     schemaPamSettingsField(),      // ✅ maps to pamSettings
    "login":            schemaLoginField(),            // ⚠️ NOT in MCP (but valid for machine login)
    "password":         schemaPasswordField(""),       // ⚠️ NOT in MCP (but valid for machine login)
    "rotation_scripts": schemaScriptField(),           // ✅ maps to rotationScripts.*
    "operating_system": schemaTextField(),             // ✅ maps to operatingSystem
    "ssl_verification": schemaCheckboxField(),         // ❌ NOT IN MCP SCHEMA
    "instance_name":    schemaTextField(),             // ✅ maps to instanceName
    "instance_id":      schemaTextField(),             // ✅ maps to instanceId
    "provider_group":   schemaTextField(),             // ✅ maps to providerGroup
    "provider_region":  schemaTextField(),             // ✅ maps to providerRegion
    "file_ref":         schemaFileRefField(),          // ✅ maps to fileRef
    "totp":             schemaOneTimeCodeField(),      // ✅ maps to oneTimeCode
}
```

### Analysis
**Status:** ❌ **INVALID FIELD DETECTED**

**Problem:** `ssl_verification` field exists in Terraform but NOT in MCP schema

**Action Required:**
1. Check Go SDK source code for `ssl_verification` field
2. Verify with backend team if this is:
   - Legacy field (should be removed)
   - Undocumented field (MCP schema needs update)
   - Terraform-specific field (document why)
3. If invalid: Remove from resource_pam_machine.go line 62

**Lines to investigate:**
- `resource_pam_machine.go:62` - Schema definition
- `resource_pam_machine.go:168-177` - Create function (if used)
- Search for read function usage

---

## 3. pamUser Schema Validation

### MCP Schema (Official - 11 fields)
```
✓ login                     (login, required)
✓ password                  (password, optional)
✓ rotationScripts.command   (string, optional)
✓ rotationScripts.fileRef   (string, optional)
✓ rotationScripts.recordRef (string, optional)
✓ privatePEMKey             (secret, optional) ❌ MISSING IN TERRAFORM
✓ distinguishedName         (text, optional)
✓ connectDatabase           (text, optional)
✓ managed                   (checkbox, optional)
✓ fileRef                   (fileRef, optional)
✓ oneTimeCode               (otp, optional)
```

### Terraform Implementation (Lines 23-63)
```go
Schema: map[string]*schema.Schema{
    "login":              schemaLoginField(),          // ✅ maps to login
    "password":           schemaPasswordField(""),     // ✅ maps to password
    "rotation_scripts":   schemaScriptField(),         // ✅ maps to rotationScripts.*
    // ❌ MISSING: "private_pem_key"
    "distinguished_name": schemaTextField(),           // ✅ maps to distinguishedName
    "connect_database":   schemaTextField(),           // ✅ maps to connectDatabase
    "managed":            schemaCheckboxField(),       // ✅ maps to managed
    "file_ref":           schemaFileRefField(),        // ✅ maps to fileRef
    "totp":               schemaOneTimeCodeField(),    // ✅ maps to oneTimeCode
}
```

### Analysis
**Status:** ❌ **MISSING FIELD**

**Problem:** `privatePEMKey` field in MCP schema but NOT in Terraform

**Impact:** Users cannot manage SSH private keys for PAM users via Terraform

**Action Required:**
Add to resource_pam_user.go:
```go
"private_pem_key": schemaSecretField(),  // maps to privatePEMKey
```

**Implementation Notes:**
- Field type: `secret` (sensitive, masked)
- Required: No (optional)
- Use case: SSH key authentication for PAM users
- Priority: MEDIUM (legitimate use case)

---

## 4. pamDirectory Schema Validation

### MCP Schema (Official - 19 fields)
```
✓ pamHostname.hostName      (string, optional)
✓ pamHostname.port          (string, optional)
✓ useSSL                    (checkbox, optional) ⚠️ LABEL BUG IN TERRAFORM
✓ pamSettings               (pamSettings, optional)
✓ trafficEncryptionSeed     (trafficEncryptionSeed, optional)
✓ rotationScripts.command   (string, optional)
✓ rotationScripts.fileRef   (string, optional)
✓ rotationScripts.recordRef (string, optional)
✓ domainName                (text, optional) ❌ MISSING IN TERRAFORM
✓ alternativeIPs            (unknown, optional) ❌ MISSING IN TERRAFORM ⚠️ MCP ERROR
✓ directoryId               (text, optional) ❌ MISSING IN TERRAFORM
✓ directoryType             (dropdown, optional)
✓ userMatch                 (text, optional) ❌ MISSING IN TERRAFORM
✓ providerGroup             (text, optional) ❌ MISSING IN TERRAFORM
✓ providerRegion            (text, optional) ❌ MISSING IN TERRAFORM
✓ fileRef                   (fileRef, optional)
✓ oneTimeCode               (otp, optional)
```

### Terraform Implementation (Lines 24-66)
```go
Schema: map[string]*schema.Schema{
    "pam_hostname":       schemaPamHostnameField(),    // ✅ maps to pamHostname.*
    "pam_settings":       schemaPamSettingsField(),    // ✅ maps to pamSettings
    "directory_type":     schemaDirectoryTypeField(),  // ✅ maps to directoryType
    "login":              schemaLoginField(),          // ⚠️ NOT in MCP (but valid for bind user)
    "password":           schemaPasswordField(""),     // ⚠️ NOT in MCP (but valid for bind user)
    "rotation_scripts":   schemaScriptField(),         // ✅ maps to rotationScripts.*
    "use_ssl":            schemaCheckboxField(),       // ⚠️ WRONG LABEL (see below)
    "distinguished_name": schemaTextField(),           // ✅ maps to distinguishedName (NOT in MCP but valid)
    // ❌ MISSING: "domain_name"
    // ❌ MISSING: "alternative_ips"
    // ❌ MISSING: "directory_id"
    // ❌ MISSING: "user_match"
    // ❌ MISSING: "provider_group"
    // ❌ MISSING: "provider_region"
    "file_ref":           schemaFileRefField(),        // ✅ maps to fileRef
    "totp":               schemaOneTimeCodeField(),    // ✅ maps to oneTimeCode
}
```

### Analysis
**Status:** ❌ **CRITICAL BUGS + MISSING FIELDS**

**CRITICAL: Label Bug**
- Line 170: `field.(*core.Checkbox).Label = "Use SSL"` ❌ WRONG
- Line 335: `getFieldResourceDataWithLabel("checkbox", "fields", secret, "Use SSL")` ❌ WRONG
- Should be: `"useSSL"` (camelCase, matching MCP schema)
- Impact: Field writes fail or read incorrectly

**Missing Fields (6 total):**
1. `domain_name` (text) - Domain name for directory
2. `alternative_ips` (unknown type - MCP shows error!) - Alternative IP addresses
3. `directory_id` (text) - Directory identifier
4. `user_match` (text) - User matching pattern
5. `provider_group` (text) - Cloud provider group
6. `provider_region` (text) - Cloud provider region

**Action Required:**
1. **IMMEDIATE**: Fix `use_ssl` label bug (2 lines)
2. **MEDIUM**: Add 6 missing fields to schema

**MCP Schema Issue:**
- `alternativeIPs` shows: `"Error: Referenced field definition ($ref) not found in fields.json"`
- This is a backend schema bug - report to Keeper team
- Type should likely be `multiline` (based on JSON template)

---

## Summary Table

| Record Type | MCP Fields | TF Fields | Status | Missing in TF | Extra in TF | Bugs |
|-------------|-----------|-----------|--------|---------------|-------------|------|
| pamDatabase | 14 | 13 | ✅ OK | 0 | login, password, connect_database | None |
| pamMachine | 14 | 13 | ❌ INVALID | 0 | login, password, **ssl_verification** | ssl_verification not in MCP |
| pamUser | 11 | 8 | ❌ INCOMPLETE | **privatePEMKey** | 0 | Missing SSH key field |
| pamDirectory | 19 | 10 | ❌ CRITICAL | **6 fields** | login, password, distinguished_name | **use_ssl label bug** |

---

## Action Items Priority

### 🔴 CRITICAL (Must Fix)
1. **Fix pamDirectory useSSL label** (2 lines in resource_pam_directory.go)
   - Line 170: Change `"Use SSL"` → `"useSSL"`
   - Line 335: Change `"Use SSL"` → `"useSSL"`

### 🟡 HIGH (Should Fix)
2. **Investigate ssl_verification in pamMachine**
   - Check Go SDK source code
   - Verify with backend team
   - Remove if invalid OR document if valid

### 🟢 MEDIUM (Nice to Have)
3. **Add privatePEMKey to pamUser**
   - Enables SSH key management
   - Add `private_pem_key` field with schemaSecretField()

4. **Add missing fields to pamDirectory**
   - domain_name, directory_id, user_match
   - provider_group, provider_region
   - alternative_ips (wait for backend fix on type)

5. **Add missing fields to pamDatabase**
   - Consider if provider_group/provider_region needed

### 📋 DOCUMENTATION
6. **Report alternativeIPs schema bug to Keeper backend team**
   - MCP schema shows: "Error: Referenced field definition ($ref) not found"
   - Expected type: multiline (based on JSON template)

---

## Testing Recommendations

After fixes, run these tests:

```bash
# 1. Fix useSSL label bug
export TF_ACC=1
export KEEPER_CREDENTIAL=<xxx>
go test ./secretsmanager -v -run TestAccResourcePamDirectory

# 2. Investigate ssl_verification
go test ./secretsmanager -v -run TestAccResourcePamMachine

# 3. Add privatePEMKey (after implementation)
go test ./secretsmanager -v -run TestAccResourcePamUser

# 4. Full PAM test suite
go test ./secretsmanager -v -run TestAccResourcePam
go test ./secretsmanager -v -run TestAccDataSourcePam
```

---

## Conclusion

**Source of Truth Established:** MCP `get_record_type_schema` is the authoritative source for field validation.

**Next Steps:**
1. Apply critical fixes (useSSL label)
2. Investigate and resolve ssl_verification
3. Add missing fields based on priority
4. Update documentation to reference MCP schema as validation method

**Long-term:**
- Create automated schema validation in CI/CD
- Add schema version tracking
- Monitor MCP schema changes for new fields
