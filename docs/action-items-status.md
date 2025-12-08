# PAM Fields Action Items - Status Update
**Updated:** 2025-12-08
**Original Plan:** `/plans/pam-fields-verification-2025-12-08.md`

## ✅ Completed (High Priority)

### 1. Fix pamDirectory useSSL label inconsistency
- **Status:** ✅ DONE
- **Commit:** 7276cf0
- **Files:** resource_pam_directory.go (lines 170, 335), data_source_pam_directory.go (line 124)

### 2. Investigate ssl_verification on pamMachine
- **Status:** ✅ DONE
- **Result:** Field is INVALID - does not exist in backend
- **Action:** Removed in commit 1746711
- **Documentation:** docs/ssl-verification-investigation.md

---

## 🟡 Pending (Testing - Requires Credentials)

### 3. Test fixes with actual pamDirectory record creation/read
- **Status:** ⏳ PENDING (needs TF_ACC credentials)
- **Command:**
  ```bash
  export TF_ACC=1
  export KEEPER_CREDENTIAL=<xxx>
  go test ./secretsmanager -v -run TestAccResourcePamDirectory
  ```
- **Blocker:** None - tests compile, just needs credentials to run
- **Priority:** HIGH (verify fixes work)

### 4. Test all PAM resources
- **Status:** ⏳ PENDING (needs TF_ACC credentials)
- **Command:**
  ```bash
  go test ./secretsmanager -v -run TestAccResourcePam
  go test ./secretsmanager -v -run TestAccDataSourcePam
  ```
- **Priority:** HIGH (pre-release verification)

---

## 🟢 Optional (Medium Priority - Future Enhancement)

### 5. Verify connectDatabase field on pamDatabase
- **Status:** 🔵 OPTIONAL
- **Notes:**
  - Field exists in Terraform but NOT in MCP schema
  - May be legacy/undocumented field
  - No user complaints yet
- **Action:** Test with actual database records to see if field is populated
- **Decision:** Can defer to post-v1.1.8

### 6. Add missing pamDirectory fields (6 fields)
- **Status:** 🔵 OPTIONAL (no user requests yet)
- **Fields:**
  - domain_name (text)
  - directory_id (text)
  - user_match (text)
  - provider_group (text)
  - provider_region (text)
  - alternative_ips (multiline - has MCP schema bug)
- **Decision:** Can defer to v1.1.9 or v1.2.0
- **Reason:** No blocking user needs identified

### 7. Add private_pem_key field to pamUser
- **Status:** 🔵 OPTIONAL (no user requests yet)
- **Field:** privatePEMKey (secret type)
- **Use Case:** SSH key authentication for PAM users
- **Decision:** Can defer to v1.1.9 or v1.2.0
- **Reason:** No blocking user needs identified

---

## 🔍 Investigation Items (Lower Priority)

### 8. Create test pamDirectory record in vault
- **Status:** 🔵 OPTIONAL
- **Purpose:** Manual verification of all fields
- **Command:**
  ```bash
  keeper record-add \
    --record-type pamDirectory \
    --title "Test PAM Directory" \
    --fields "pamHostname:ldap.example.com:636" \
    --fields "useSSL:true" \
    --fields "directoryType:Active Directory" \
    --fields "domainName:example.com"

  terraform import secretsmanager_pam_directory.test <record-uid>
  ```
- **Decision:** Can be done during acceptance testing

### 9. Update record-templates repo
- **Status:** 🔵 OPTIONAL
- **Fields to add:**
  - connectDatabase to pamUser.json (confirmed valid in MCP schema)
  - Possibly other fields found valid but missing from templates
- **Decision:** Requires coordination with template maintainers
- **Action:** Create issue in record-templates repo

---

## Summary by Priority

### 🔴 BLOCKING v1.1.8 Release
**None** - All critical fixes completed

### 🟡 RECOMMENDED Before Release
- [ ] Run acceptance tests with TF_ACC credentials (items #3, #4)
- Estimated time: 10-15 minutes

### 🟢 Post-Release Enhancements
- [ ] Add missing pamDirectory fields (item #6)
- [ ] Add privatePEMKey to pamUser (item #7)
- [ ] Verify connectDatabase field (item #5)
- [ ] Update record-templates repo (item #9)

### 🔵 Nice to Have
- [ ] Manual vault testing (item #8)

---

## Recommendation

**For v1.1.8 Release:**
1. ✅ All critical fixes are DONE and committed
2. ⏳ Run acceptance tests if you have TF_ACC credentials handy
3. ✅ If tests pass, ready to release
4. 🟢 Defer optional enhancements to v1.1.9/v1.2.0

**Timeline:**
- **Now:** Ready to merge to master (all fixes done)
- **Before tagging v1.1.8:** Run acceptance tests (10-15 min)
- **After release:** Plan feature additions for next version

---

## Decision Required

**Question:** Do you have TF_ACC credentials available to run acceptance tests now?

**Option A - Yes:**
```bash
export TF_ACC=1
export KEEPER_CREDENTIAL=<your-credential>
go test ./secretsmanager -v -run TestAccResourcePamDirectory
go test ./secretsmanager -v -run TestAccResourcePamMachine
```
→ Verify fixes, then release

**Option B - No:**
→ Release v1.1.8 with fixes (tests compile cleanly, high confidence)
→ Run acceptance tests post-release if issues arise

**Recommendation:** Option A if credentials readily available, Option B acceptable given build success and investigation thoroughness.
