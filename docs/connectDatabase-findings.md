# connectDatabase Field Investigation
**Date:** 2025-12-08
**Status:** 🔴 **TWO ISSUES FOUND**

## Summary

Investigation revealed **TWO separate problems** with the `connect_database` field:

1. **pamDatabase**: Field exists in Terraform but NOT in MCP schema → **REMOVE**
2. **pamUser**: Field exists in both but has wrong label → **FIX LABEL**

---

## MCP Schema Verification

### pamDatabase MCP Schema (14 fields)
```json
{
  "fields": [
    "pamHostname.hostName",
    "pamHostname.port",
    "useSSL",
    "pamSettings",
    "trafficEncryptionSeed",
    "rotationScripts.*",
    "databaseId",
    "databaseType",
    "providerGroup",
    "providerRegion",
    "fileRef",
    "oneTimeCode"
  ]
}
```
**Result:** ❌ NO `connectDatabase` field

### pamUser MCP Schema (11 fields)
```json
{
  "fields": [
    "login",
    "password",
    "rotationScripts.*",
    "privatePEMKey",
    "distinguishedName",
    "connectDatabase",  // ✅ FIELD EXISTS!
    "managed",
    "fileRef",
    "oneTimeCode"
  ]
}
```
**Result:** ✅ HAS `connectDatabase` field (camelCase)

---

## Problem #1: pamDatabase Has Invalid Field

### Current Implementation (WRONG)

**resource_pam_database.go:**
```go
// Line 62 - Schema
"connect_database": schemaTextField(),

// Line 165-174 - Create
if fieldData := d.Get("connect_database"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
    if field, err := NewFieldFromSchema("text", fieldData); err != nil {
        return diag.FromErr(err)
    } else if field != nil {
        field.(*core.Text).Label = "Connect Database"  // ❌ Field doesn't exist in backend
        nrc.Fields = append(nrc.Fields, field)
        ...
    }
}

// Line 350 - Read
connectDatabase := getFieldResourceDataWithLabel("text", "fields", secret, "Connect Database")
if err = d.Set("connect_database", connectDatabase); err != nil {
    return diag.FromErr(err)
}
```

**data_source_pam_database.go:**
```go
// Line 49 - Schema
"connect_database": schemaTextField(),

// Line 116-117 - Read
connectDatabase := getFieldResourceDataWithLabel("text", "fields", secret, "Connect Database")
if err = d.Set("connect_database", connectDatabase); err != nil {
    return diag.FromErr(err)
}
```

### Fix Required

**Action:** Remove `connect_database` field entirely from pamDatabase

**Reason:** Field does not exist in MCP backend schema

**Impact:** Breaking change (field never released - v1.1.7 has no PAM)

**Files to modify:**
- `secretsmanager/resource_pam_database.go` (3 locations)
- `secretsmanager/data_source_pam_database.go` (2 locations)

---

## Problem #2: pamUser Has Wrong Label

### Current Implementation (WRONG)

**resource_pam_user.go:**
```go
// Line 59 - Schema
"connect_database": schemaTextField(),  // ✅ Schema OK

// Line 145-152 - Create
if fieldData := d.Get("connect_database"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
    if field, err := NewFieldFromSchema("text", fieldData); err != nil {
        return diag.FromErr(err)
    } else if field != nil {
        field.(*core.Text).Label = "Connect Database"  // ❌ WRONG (has space)
        nrc.Fields = append(nrc.Fields, field)
        ...
    }
}

// Line 291-292 - Read
connectDatabase := getFieldResourceDataWithLabel("text", "fields", secret, "Connect Database")  // ❌ WRONG
if err = d.Set("connect_database", connectDatabase); err != nil {
    return diag.FromErr(err)
}

// Line 357-358 - Update (uses ApplyFieldChange - might work?)
if d.HasChange("connect_database") {
    if _, err := ApplyFieldChange("fields", "connect_database", d, secret); err != nil {
        return diag.FromErr(err)
    }
}
```

**data_source_pam_user.go:**
```go
// Line 47 - Schema
"connect_database": schemaTextField(),  // ✅ Schema OK

// Line 99-100 - Read
connectDatabase := getFieldResourceDataWithLabel("text", "fields", secret, "Connect Database")  // ❌ WRONG
if err = d.Set("connect_database", connectDatabase); err != nil {
    return diag.FromErr(err)
}
```

### Fix Required

**Action:** Change label from `"Connect Database"` → `"connectDatabase"`

**Reason:** MCP schema uses camelCase field name

**Impact:** Bug fix (same as useSSL issue)

**Files to modify:**
- `secretsmanager/resource_pam_user.go` (lines 149, 291)
- `secretsmanager/data_source_pam_user.go` (line 99)

---

## Comparison with Similar Fields

| Record Type | Field | TF Label | MCP Field Name | Status |
|-------------|-------|----------|----------------|--------|
| pamDatabase | connect_database | "Connect Database" | ❌ Does not exist | **REMOVE** |
| pamUser | connect_database | "Connect Database" | `connectDatabase` | **FIX LABEL** |
| pamDirectory | use_ssl | ~~"Use SSL"~~ | `useSSL` | ✅ Fixed (commit 7276cf0) |
| pamDatabase | use_ssl | "useSSL" | `useSSL` | ✅ Correct |

**Pattern:** Backend uses camelCase, not human-readable labels with spaces

---

## Root Cause

Same issue as the useSSL bug:
- Developer used human-readable labels instead of backend field names
- Copy-paste between pamDatabase and pamUser without validation
- pamDatabase incorrectly has field that doesn't exist in backend

---

## Action Plan

### Step 1: Remove from pamDatabase

**resource_pam_database.go:**
```diff
# Line 62 - Remove from schema
-			"connect_database": schemaTextField(),

# Lines 165-174 - Remove from create function
-	if fieldData := d.Get("connect_database"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
-		if field, err := NewFieldFromSchema("text", fieldData); err != nil {
-			return diag.FromErr(err)
-		} else if field != nil {
-			field.(*core.Text).Label = "Connect Database"
-			nrc.Fields = append(nrc.Fields, field)
-			if err := SetFieldTypeInSchema(d, "connect_database", "text"); err != nil {
-				return diag.FromErr(err)
-			}
-		}
-	}

# Lines 348-351 - Remove from read function
-	connectDatabase := getFieldResourceDataWithLabel("text", "fields", secret, "Connect Database")
-	if err = d.Set("connect_database", connectDatabase); err != nil {
-		return diag.FromErr(err)
-	}
```

**data_source_pam_database.go:**
```diff
# Line 49 - Remove from schema
-			"connect_database": schemaTextField(),

# Lines 116-117 - Remove from read function
-	connectDatabase := getFieldResourceDataWithLabel("text", "fields", secret, "Connect Database")
-	if err = d.Set("connect_database", connectDatabase); err != nil {
-		return diag.FromErr(err)
-	}
```

### Step 2: Fix label in pamUser

**resource_pam_user.go:**
```diff
# Line 149 - Fix label in create
-			field.(*core.Text).Label = "Connect Database"
+			field.(*core.Text).Label = "connectDatabase"

# Line 291 - Fix label in read
-	connectDatabase := getFieldResourceDataWithLabel("text", "fields", secret, "Connect Database")
+	connectDatabase := getFieldResourceDataWithLabel("text", "fields", secret, "connectDatabase")
```

**data_source_pam_user.go:**
```diff
# Line 99 - Fix label in read
-	connectDatabase := getFieldResourceDataWithLabel("text", "fields", secret, "Connect Database")
+	connectDatabase := getFieldResourceDataWithLabel("text", "fields", secret, "connectDatabase")
```

### Step 3: Update test file (if needed)

**resource_pam_user_test.go line 45:**
```diff
-				label = "Connect Database"
+				label = "connectDatabase"
```
*Note: Actually, the label in test config is ignored - code overrides it. But update for consistency.*

---

## Testing

```bash
# Build
go build ./...

# Test pamUser (after label fix)
export TF_ACC=1
export KEEPER_CREDENTIAL=<xxx>
go test ./secretsmanager -v -run TestAccResourcePamUser

# Test pamDatabase (after field removal)
go test ./secretsmanager -v -run TestAccResourcePamDatabase
```

---

## Commit Plan

**Commit 1: Remove connect_database from pamDatabase**
```
fix(pamDatabase): remove invalid connect_database field
```

**Commit 2: Fix label in pamUser**
```
fix(pamUser): correct connectDatabase label from 'Connect Database' to 'connectDatabase'
```

---

## Impact Assessment

### pamDatabase Field Removal
- **Breaking:** No (never released - v1.1.7 has no PAM)
- **Risk:** Low (field doesn't work anyway - not in backend)

### pamUser Label Fix
- **Breaking:** No (bug fix - field wasn't working correctly)
- **Risk:** Low (same as useSSL fix)

---

## Conclusion

Two separate issues found:
1. ✅ pamDatabase should NOT have connect_database field → Remove
2. ✅ pamUser should have field but with correct label → Fix "Connect Database" → "connectDatabase"

Both fixes align provider with MCP backend schema.
