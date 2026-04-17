package secretsmanager

import (
	"context"
	"testing"

	fwephemeral "github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestCustomFieldDataSourceSchemaPresence verifies that every record-type data
// source schema contains the "custom" TypeList key with the expected 5
// sub-attributes. No vault required.
func TestCustomFieldDataSourceSchemaPresence(t *testing.T) {
	dataSources := map[string]*schema.Resource{
		"address":              dataSourceAddress(),
		"bank_account":         dataSourceBankAccount(),
		"bank_card":            dataSourceBankCard(),
		"birth_certificate":    dataSourceBirthCertificate(),
		"contact":              dataSourceContact(),
		"database_credentials": dataSourceDatabaseCredentials(),
		"driver_license":       dataSourceDriverLicense(),
		"encrypted_notes":      dataSourceEncryptedNotes(),
		"health_insurance":     dataSourceHealthInsurance(),
		"login":                dataSourceLogin(),
		"membership":           dataSourceMembership(),
		"pam_database":         dataSourcePamDatabase(),
		"pam_directory":        dataSourcePamDirectory(),
		"pam_remote_browser":   dataSourcePamRemoteBrowser(),
		"passport":             dataSourcePassport(),
		"photo":                dataSourcePhoto(),
		"server_credentials":   dataSourceServerCredentials(),
		"software_license":     dataSourceSoftwareLicense(),
		"ssh_keys":             dataSourceSshKeys(),
		"ssn_card":             dataSourceSsnCard(),
	}

	requiredSubKeys := []string{"type", "label", "value", "required", "privacy_screen"}

	for name, res := range dataSources {
		t.Run(name, func(t *testing.T) {
			customSchema, ok := res.Schema["custom"]
			if !ok {
				t.Fatalf("data source %q: \"custom\" key missing from schema", name)
			}
			if customSchema.Type != schema.TypeList {
				t.Errorf("data source %q: \"custom\" is %v, want TypeList", name, customSchema.Type)
			}
			if !customSchema.Computed {
				t.Errorf("data source %q: \"custom\" must be Computed (read-only)", name)
			}
			elem, ok := customSchema.Elem.(*schema.Resource)
			if !ok {
				t.Fatalf("data source %q: \"custom\" Elem is not *schema.Resource", name)
			}
			for _, key := range requiredSubKeys {
				sub, found := elem.Schema[key]
				if !found {
					t.Errorf("data source %q: \"custom\" block missing sub-key %q", name, key)
					continue
				}
				if !sub.Computed {
					t.Errorf("data source %q: \"custom.%s\" must be Computed", name, key)
				}
			}
		})
	}
}

// ephemeralSchemaProvider is the subset of the Plugin Framework ephemeral
// resource interface that exposes the schema.
type ephemeralSchemaProvider interface {
	Schema(context.Context, fwephemeral.SchemaRequest, *fwephemeral.SchemaResponse)
}

// TestCustomFieldEphemeralSchemaPresence verifies that every record-type
// ephemeral resource schema contains the "custom" attribute. No vault required.
func TestCustomFieldEphemeralSchemaPresence(t *testing.T) {
	ephemeralResources := map[string]fwephemeral.EphemeralResource{
		"address":              NewEphemeralAddress(),
		"bank_account":         NewEphemeralBankAccount(),
		"bank_card":            NewEphemeralBankCard(),
		"birth_certificate":    NewEphemeralBirthCertificate(),
		"contact":              NewEphemeralContact(),
		"database_credentials": NewEphemeralDatabaseCredentials(),
		"driver_license":       NewEphemeralDriverLicense(),
		"encrypted_notes":      NewEphemeralEncryptedNotes(),
		"health_insurance":     NewEphemeralHealthInsurance(),
		"login":                NewEphemeralLogin(),
		"membership":           NewEphemeralMembership(),
		"pam_database":         NewEphemeralPamDatabase(),
		"pam_directory":        NewEphemeralPamDirectory(),
		"pam_machine":          NewEphemeralPamMachine(),
		"pam_remote_browser":   NewEphemeralPamRemoteBrowser(),
		"pam_user":             NewEphemeralPamUser(),
		"passport":             NewEphemeralPassport(),
		"photo":                NewEphemeralPhoto(),
		"server_credentials":   NewEphemeralServerCredentials(),
		"software_license":     NewEphemeralSoftwareLicense(),
		"ssh_keys":             NewEphemeralSshKeys(),
		"ssn_card":             NewEphemeralSsnCard(),
	}

	for name, res := range ephemeralResources {
		t.Run(name, func(t *testing.T) {
			sp, ok := res.(ephemeralSchemaProvider)
			if !ok {
				t.Fatalf("ephemeral %q: does not implement Schema()", name)
			}
			var resp fwephemeral.SchemaResponse
			sp.Schema(context.Background(), fwephemeral.SchemaRequest{}, &resp)
			if _, found := resp.Schema.Attributes["custom"]; !found {
				t.Fatalf("ephemeral %q: \"custom\" attribute missing from schema", name)
			}
		})
	}
}
