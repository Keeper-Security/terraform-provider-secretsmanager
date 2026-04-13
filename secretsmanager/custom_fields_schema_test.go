package secretsmanager

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestCustomFieldSchemaPresence verifies that every record resource schema
// contains the "custom" TypeList key with the expected 5 sub-attributes.
//
// This is a structural check — it does not require a live vault. It catches:
//   - Resources where schemaCustomField() was accidentally omitted
//   - Sub-attribute regressions (missing or renamed keys)
//
// Runtime wiring (HasChange/d.Get) is verified by per-resource acceptance tests
// and the grep in CI.
func TestCustomFieldSchemaPresence(t *testing.T) {
	resources := map[string]*schema.Resource{
		"address":              resourceAddress(),
		"bank_account":         resourceBankAccount(),
		"bank_card":            resourceBankCard(),
		"birth_certificate":    resourceBirthCertificate(),
		"contact":              resourceContact(),
		"database_credentials": resourceDatabaseCredentials(),
		"driver_license":       resourceDriverLicense(),
		"encrypted_notes":      resourceEncryptedNotes(),
		"health_insurance":     resourceHealthInsurance(),
		"login":                resourceLogin(),
		"membership":           resourceMembership(),
		"pam_database":         resourcePamDatabase(),
		"pam_directory":        resourcePamDirectory(),
		"pam_machine":          resourcePamMachine(),
		"pam_remote_browser":   resourcePamRemoteBrowser(),
		"pam_user":             resourcePamUser(),
		"passport":             resourcePassport(),
		"photo":                resourcePhoto(),
		"server_credentials":   resourceServerCredentials(),
		"software_license":     resourceSoftwareLicense(),
		"ssh_keys":             resourceSshKeys(),
		"ssn_card":             resourceSsnCard(),
	}

	requiredSubKeys := []string{"type", "label", "value", "required", "privacy_screen"}

	for name, res := range resources {
		t.Run(name, func(t *testing.T) {
			customSchema, ok := res.Schema["custom"]
			if !ok {
				t.Fatalf("resource %q: \"custom\" key missing from schema", name)
			}
			if customSchema.Type != schema.TypeList {
				t.Errorf("resource %q: \"custom\" is %v, want TypeList", name, customSchema.Type)
			}
			elem, ok := customSchema.Elem.(*schema.Resource)
			if !ok {
				t.Fatalf("resource %q: \"custom\" Elem is not *schema.Resource", name)
			}
			for _, key := range requiredSubKeys {
				if _, found := elem.Schema[key]; !found {
					t.Errorf("resource %q: \"custom\" block missing sub-key %q", name, key)
				}
			}
		})
	}
}
