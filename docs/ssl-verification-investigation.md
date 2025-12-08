# ssl_verification Field Investigation Report
**Date:** 2025-12-08
**Field:** `ssl_verification` in `pamMachine` resource
**Status:** ❌ **INVALID - Should be REMOVED**

## Executive Summary

The `ssl_verification` field in the `pamMachine` resource **does not exist in the official Keeper backend schema** and should be removed from the Terraform provider.

**Evidence:**
1. ✅ MCP `get_record_type_schema` confirms field does NOT exist in backend
2. ✅ Go SDK v1.6.4 has NO reference to "SSL Verification" field
3. ✅ Field is commented out in official examples (suggesting uncertainty)
4. ✅ No other PAM record types have this field

**Recommendation:** Remove from Terraform provider (breaking change → requires version bump)

---

## Investigation Details

### 1. Origin

**When added:** Commit `470998e` (Oct 24, 2025)
**Author:** Max Ustinov
**Commit message:** "KSM-527: Add support for PAM record types"

Excerpt:
```
New PAM-specific field types:
- pamHostname: hostname/port configuration
- checkbox: boolean checkbox fields
...
All implementations follow existing Terraform provider patterns and are based on KSM SDK specifications.
```

**Analysis:** Field was added during initial PAM implementation, likely based on:
- Preliminary specification
- Assumption about needed fields
- Confusion with other SSL-related settings in `pamSettings`

---

### 2. Current Implementation

**Schema Definition:**
```go
// secretsmanager/resource_pam_machine.go:62
"ssl_verification": schemaCheckboxField(),
```

**Create Function:**
```go
// secretsmanager/resource_pam_machine.go:168-177
if fieldData := d.Get("ssl_verification"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
    if field, err := NewFieldFromSchema("checkbox", fieldData); err != nil {
        return diag.FromErr(err)
    } else if field != nil {
        field.(*core.Checkbox).Label = "SSL Verification"  // ❌ No such field in backend
        nrc.Fields = append(nrc.Fields, field)
        if err := SetFieldTypeInSchema(d, "ssl_verification", "checkbox"); err != nil {
            return diag.FromErr(err)
        }
    }
}
```

**Read Function:**
```go
// secretsmanager/resource_pam_machine.go:359-361
sslVerification := getFieldResourceDataWithLabel("checkbox", "fields", secret, "SSL Verification")
if err = d.Set("ssl_verification", sslVerification); err != nil {
    return diag.FromErr(err)
}
```

**Data Source:**
```go
// secretsmanager/data_source_pam_machine.go:49
"ssl_verification": schemaCheckboxField(),

// secretsmanager/data_source_pam_machine.go:117
if err = d.Set("ssl_verification", sslVerification); err != nil {
    return diag.FromErr(err)
}
```

**Provider Type Mapping:**
```go
// secretsmanager/provider.go:176
"ssl_verification":   "checkbox", // Checkbox with label
```

---

### 3. Evidence: MCP Schema (Authoritative Source)

**Query:**
```
mcp__ksm__get_record_type_schema(type="pamMachine")
```

**Result:**
```json
{
  "record_type": "pamMachine",
  "fields": [
    {"name": "pamHostname.hostName", "type": "string", "required": true},
    {"name": "pamHostname.port", "type": "string", "required": true},
    {"name": "pamSettings", "type": "pamSettings", "required": false},
    {"name": "trafficEncryptionSeed", "type": "trafficEncryptionSeed", "required": false},
    {"name": "rotationScripts.command", "type": "string", "required": false},
    {"name": "rotationScripts.fileRef", "type": "string", "required": false},
    {"name": "rotationScripts.recordRef", "type": "string", "required": false},
    {"name": "operatingSystem", "type": "text", "required": false},
    {"name": "instanceName", "type": "text", "required": false},
    {"name": "instanceId", "type": "text", "required": false},
    {"name": "providerGroup", "type": "text", "required": false},
    {"name": "providerRegion", "type": "text", "required": false},
    {"name": "fileRef", "type": "fileRef", "required": false},
    {"name": "oneTimeCode", "type": "otp", "required": false}
  ]
}
```

**Finding:** NO `ssl_verification` or similar field in official schema (14 fields total)

---

### 4. Evidence: Go SDK

**Search command:**
```bash
grep -rn "SSL Verification" /Users/stasschaller/go/pkg/mod/github.com/keeper-security/secrets-manager-go/core@v1.6.4/*.go
```

**Result:** `Not found in Go SDK v1.6.4`

**Finding:** The Go SDK v1.6.4 (latest as of 2025-12-08) has NO reference to "SSL Verification" field

---

### 5. Evidence: Official Examples

**File:** `examples/resources/pam_machine.tf`

**Lines 114-117:**
```hcl
# Optional: SSL verification
# ssl_verification {
#   value = [true]
# }
```

**Finding:** Field is **commented out** in the official example, suggesting:
- Uncertainty about its validity
- Not commonly used
- May have been added speculatively

---

### 6. Comparison with Other PAM Types

| PAM Type | Has ssl_verification? | Notes |
|----------|----------------------|-------|
| pamDatabase | ❌ No | Has `useSSL` field instead |
| pamMachine | ⚠️ In TF only | NOT in MCP schema |
| pamUser | ❌ No | No SSL-related fields |
| pamDirectory | ❌ No | Has `useSSL` field instead |

**Finding:** Only `pamMachine` has `ssl_verification`, and it's not in the backend schema

**Note:** `pamDatabase` and `pamDirectory` have `useSSL` (camelCase) which IS in MCP schema

---

### 7. Possible Confusion

**Hypothesis:** Developer may have confused two different concepts:

1. **Record-level `useSSL` field** (pamDatabase, pamDirectory)
   - Checkbox field with label "useSSL"
   - Controls SSL/TLS for database/directory connections
   - ✅ **Official field in MCP schema**

2. **Connection-level SSL settings** (in `pamSettings.connection`)
   - RDP: `ignoreCert` - ignore SSL certificate errors
   - SSH: No SSL settings (uses SSH keys)
   - ✅ **Part of pamSettings JSON**

3. **Non-existent `ssl_verification`** (pamMachine only)
   - Checkbox field with label "SSL Verification"
   - ❌ **NOT in MCP schema**
   - ❌ **NOT in Go SDK**
   - ❌ **Appears to be a mistake**

---

## Impact Analysis

### If Field is Removed

**Breaking Change:** Yes
- Users with `ssl_verification` in their Terraform configs will get errors
- Requires provider version bump (e.g., 1.1.8 → 1.2.0)

**Migration Path:**
```hcl
# BEFORE (v1.1.x)
resource "secretsmanager_pam_machine" "example" {
  # ...
  ssl_verification {
    value = [true]
  }
}

# AFTER (v1.2.0+)
resource "secretsmanager_pam_machine" "example" {
  # ...
  # Remove ssl_verification - field does not exist in Keeper backend
  # If you need SSL/TLS control, use pamSettings.connection.ignoreCert for RDP
  pam_settings = jsonencode([{
    connection = [{
      protocol = "rdp"
      ignoreCert = false  # Verify SSL certificates
    }]
  }])
}
```

**Risk Assessment:** LOW
- Field is commented out in examples
- Likely has very few (if any) actual users
- Field doesn't actually do anything in backend (would be stored as custom field)

---

## Recommendation

### Action: REMOVE `ssl_verification` field

**Rationale:**
1. Field does not exist in official backend schema (MCP confirmation)
2. Field does not exist in Go SDK
3. Field provides no actual functionality
4. Keeping it confuses users and creates maintenance burden
5. Low risk of breaking user configs (likely unused)

### Implementation Steps

1. **Remove from resource schema** (`resource_pam_machine.go`):
   - Line 62: Remove schema definition
   - Lines 168-177: Remove create logic
   - Lines 359-361: Remove read logic
   - Check update function (may have change detection)

2. **Remove from data source** (`data_source_pam_machine.go`):
   - Line 49: Remove schema definition
   - Line 117: Remove read logic

3. **Remove from provider** (`provider.go`):
   - Line 176: Remove type mapping

4. **Update example** (`examples/resources/pam_machine.tf`):
   - Lines 114-117: Remove commented-out example

5. **Update documentation**:
   - Remove from field list
   - Add migration note to CHANGELOG

6. **Update tests** (if any reference ssl_verification):
   - `resource_pam_machine_test.go`

### Version Bump

**Current:** v1.1.8
**Proposed:** v1.2.0 (breaking change)

**CHANGELOG Entry:**
```markdown
## [1.2.0] - 2025-XX-XX

### Breaking Changes
- **pamMachine**: Removed `ssl_verification` field as it does not exist in Keeper backend schema
  - This field was added in error during initial PAM implementation
  - Confirmed via MCP schema validation and Go SDK inspection
  - No actual functionality was provided by this field
  - Migration: Remove `ssl_verification` from your Terraform configs
  - For RDP SSL control, use `pamSettings.connection.ignoreCert` instead
```

---

## Alternative: Keep as Custom Field

**Option:** Document as "legacy custom field" and mark deprecated

**Pros:**
- No breaking change
- Maintains backward compatibility
- Users can transition gradually

**Cons:**
- Confuses users ("why is this here?")
- Maintenance burden
- Inconsistent with MCP schema
- No actual backend functionality

**Verdict:** NOT RECOMMENDED - cleaner to remove

---

## Conclusion

**The `ssl_verification` field in `pamMachine` should be REMOVED** in the next major version (v1.2.0) because:

1. ✅ Confirmed NOT in Keeper backend (MCP schema)
2. ✅ Confirmed NOT in Go SDK v1.6.4
3. ✅ Commented out in official examples
4. ✅ Unique to pamMachine (no other PAM types have it)
5. ✅ Low risk of breaking user configs
6. ✅ Improves schema accuracy and maintainability

**Next Steps:**
1. Create PR to remove field
2. Update documentation and examples
3. Add CHANGELOG entry
4. Bump version to 1.2.0
5. Communicate breaking change in release notes
