# useSSL Label Bug Fix Summary
**Date:** 2025-12-08
**Status:** ✅ **COMPLETED**
**Type:** Bug fix (non-breaking)

## Problem

The `use_ssl` field in `pamDirectory` resource and data source was using the wrong label format when communicating with Keeper backend:

- **Incorrect:** `"Use SSL"` (human-readable, with space)
- **Correct:** `"useSSL"` (camelCase, per MCP schema)

This caused field writes to fail or read incorrectly from the vault.

## Root Cause

During initial PAM implementation (commit 470998e, Oct 2025), the label was mistakenly formatted as human-readable text instead of using the backend's camelCase field naming convention.

**Evidence from MCP Schema:**
```json
{
  "name": "useSSL",
  "type": "checkbox",
  "required": false
}
```

**Comparison with pamDatabase:**
- ✅ `pamDatabase` correctly uses `"useSSL"` (line 122)
- ❌ `pamDirectory` incorrectly used `"Use SSL"` (lines 170, 335)

## Changes Made

### Files Modified (3 locations)

#### 1. `secretsmanager/resource_pam_directory.go`

**Line 170 (Create function):**
```diff
- field.(*core.Checkbox).Label = "Use SSL"
+ field.(*core.Checkbox).Label = "useSSL"
```

**Line 335 (Read function):**
```diff
- useSSL := getFieldResourceDataWithLabel("checkbox", "fields", secret, "Use SSL")
+ useSSL := getFieldResourceDataWithLabel("checkbox", "fields", secret, "useSSL")
```

#### 2. `secretsmanager/data_source_pam_directory.go`

**Line 124 (Read function):**
```diff
- useSSL := getFieldResourceDataWithLabel("checkbox", "fields", secret, "Use SSL")
+ useSSL := getFieldResourceDataWithLabel("checkbox", "fields", secret, "useSSL")
```

## Testing

### Build Verification
```bash
go build ./...
# ✅ SUCCESS - No syntax errors
```

### Test Compilation
```bash
go test ./secretsmanager -run TestAccResourcePamDirectory -v
# ✅ SUCCESS - Tests compile and skip gracefully (TF_ACC not set)
```

### Acceptance Tests (Manual)
To fully verify the fix, run with credentials:
```bash
export TF_ACC=1
export KEEPER_CREDENTIAL=<xxx>
go test ./secretsmanager -v -run TestAccResourcePamDirectory
```

Expected: All tests should pass with correct field reading/writing.

## Impact

### Users Affected
- Users with `pamDirectory` resources using `use_ssl` field
- Users importing existing `pamDirectory` records from vault

### Breaking Change?
**No** - This is a bug fix that makes the provider work correctly with the backend. Users should not experience breaking changes; the field should now work as intended.

### Migration Required?
**No** - Existing Terraform configurations don't need changes. The field name (`use_ssl`) remains the same; only the internal backend label mapping is fixed.

## Example Usage (No Change)

User Terraform configs remain unchanged:
```hcl
resource "secretsmanager_pam_directory" "ldap" {
  folder_uid = "<folder-uid>"
  title = "LDAP Server"

  pam_hostname {
    value {
      hostname = "ldap.example.com"
      port = "636"
    }
  }

  use_ssl {
    value = [true]  # Now works correctly!
  }

  directory_type = "Active Directory"
}
```

## Related Issues

### Documentation Note
The example file `examples/resources/pam_directory.tf` contains:
```hcl
use_ssl {
  label = "Use SSL"  # Misleading but harmless (code overrides)
  value = [true]
}
```

**Recommendation:** Remove `label = "Use SSL"` from example to avoid confusion. The code overrides this with the correct `"useSSL"` label, so user-provided labels are ignored anyway.

### Consistency Check
Other PAM checkbox fields verified for correct labels:
- ✅ `pamDatabase.useSSL` - uses `"useSSL"` (correct)
- ✅ `pamUser.managed` - uses `"managed"` (correct per schema)

## Commit Details

**Branch:** release-v1.1.8 (or appropriate branch)

**Commit Message:**
```
fix(pamDirectory): correct useSSL label to match backend schema

The use_ssl checkbox field was using incorrect label "Use SSL"
instead of the camelCase "useSSL" required by the backend API.
This caused field operations to fail or behave incorrectly.

- Fixed resource_pam_directory.go lines 170, 335
- Fixed data_source_pam_directory.go line 124
- Verified against MCP schema: field name is "useSSL" (camelCase)
- Consistent with pamDatabase which correctly uses "useSSL"

No user-facing breaking changes - existing configs work unchanged.
```

**Files Changed:**
```
secretsmanager/resource_pam_directory.go     | 2 +-
secretsmanager/data_source_pam_directory.go  | 1 +-
2 files changed, 2 insertions(+), 2 deletions(-)
```

## Verification Checklist

- [x] Code changes applied (3 locations)
- [x] Build succeeds (`go build`)
- [x] Tests compile successfully
- [x] MCP schema verified (field is `"useSSL"`)
- [x] Consistent with pamDatabase implementation
- [ ] Acceptance tests pass (requires TF_ACC + credentials)
- [ ] PR created and reviewed
- [ ] Documentation updated if needed

## Next Steps

1. **Run full acceptance tests** with TF_ACC credentials
2. **Create PR** with commit message above
3. **Update CHANGELOG** for next release (v1.1.9 or v1.2.0)
4. **Consider updating** example file to remove misleading label

## References

- [PAM Schema Validation Report](./pam-schema-validation-report.md)
- [PAM Field Mapping](./pam-field-mapping.md)
- MCP Schema: `get_record_type_schema(type="pamDirectory")`
- Original PAM commit: 470998e
