package secretsmanager

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/keeper-security/secrets-manager-go/core"
)

// TestAccResourceLogin_customFieldSimple tests create, read, update, and delete
// of simple-type custom fields (text, secret) on a login resource.
func TestAccResourceLogin_customFieldSimple(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_custom_simple"

	configCreate := fmt.Sprintf(`
		resource "secretsmanager_login" "custom_simple" {
			folder_uid = "%v"
			uid        = "%v"
			title      = "%v"

			custom {
				type  = "text"
				label = "Environment"
				value = "production"
			}
			custom {
				type  = "secret"
				label = "ApiKey"
				value = "hunter2"
			}
		}
	`, secretFolderUid, secretUid, secretTitle)

	configUpdate := fmt.Sprintf(`
		resource "secretsmanager_login" "custom_simple" {
			folder_uid = "%v"
			uid        = "%v"
			title      = "%v"

			custom {
				type  = "text"
				label = "Environment"
				value = "staging"
			}
			custom {
				type  = "secret"
				label = "ApiKey"
				value = "hunter2"
			}
			custom {
				type  = "text"
				label = "Region"
				value = "us-east-1"
			}
		}
	`, secretFolderUid, secretUid, secretTitle)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: configCreate,
				Check: resource.ComposeTestCheckFunc(
					checkSecretExistsRemotely(secretUid),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_simple", "custom.#", "2"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_simple", "custom.0.label", "Environment"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_simple", "custom.0.value", "production"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_simple", "custom.1.label", "ApiKey"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_simple", "custom.1.value", "hunter2"),
				),
			},
			{
				Config: configUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("secretsmanager_login.custom_simple", "custom.#", "3"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_simple", "custom.0.value", "staging"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_simple", "custom.2.label", "Region"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_simple", "custom.2.value", "us-east-1"),
				),
			},
		},
	})
}

// TestAccResourceLogin_customFieldPhone tests a complex phone custom field
// using jsonencode() for the value.
func TestAccResourceLogin_customFieldPhone(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_custom_phone"

	config := fmt.Sprintf(`
		resource "secretsmanager_login" "custom_phone" {
			folder_uid = "%v"
			uid        = "%v"
			title      = "%v"

			custom {
				type  = "phone"
				label = "WorkPhone"
				value = jsonencode({
					region = "US"
					number = "555-867-5309"
					ext    = "42"
					type   = "Work"
				})
			}
		}
	`, secretFolderUid, secretUid, secretTitle)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkSecretExistsRemotely(secretUid),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_phone", "custom.#", "1"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_phone", "custom.0.label", "WorkPhone"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_phone", "custom.0.type", "phone"),
				),
			},
		},
	})
}

// TestAccResourceLogin_customFieldMultiValuePhone tests a phone field with two entries
// using the jsonencode([{...},{...}]) array syntax, and verifies no perpetual diff.
func TestAccResourceLogin_customFieldMultiValuePhone(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_custom_multi_phone"

	config := fmt.Sprintf(`
		resource "secretsmanager_login" "custom_multi_phone" {
			folder_uid = "%v"
			uid        = "%v"
			title      = "%v"

			custom {
				type  = "phone"
				label = "ContactNumbers"
				value = jsonencode([
					{ region = "US", number = "555-1234", type = "Work" },
					{ region = "US", number = "555-5678", type = "Mobile" }
				])
			}
		}
	`, secretFolderUid, secretUid, secretTitle)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkSecretExistsRemotely(secretUid),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_multi_phone", "custom.#", "1"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_multi_phone", "custom.0.type", "phone"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_multi_phone", "custom.0.label", "ContactNumbers"),
					resource.TestCheckResourceAttrSet("secretsmanager_login.custom_multi_phone", "custom.0.value"),
				),
			},
			{
				// Verify no perpetual diff — multi-value round-trip must be stable
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}

// TestAccResourceLogin_customFieldSecurityQuestion tests a securityQuestion field
// using jsonencode({question, answer}).
func TestAccResourceLogin_customFieldSecurityQuestion(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_custom_secq"

	config := fmt.Sprintf(`
		resource "secretsmanager_login" "custom_secq" {
			folder_uid = "%v"
			uid        = "%v"
			title      = "%v"

			custom {
				type  = "securityQuestion"
				label = "RecoveryQuestion"
				value = jsonencode({
					question = "What was the name of your first pet?"
					answer   = "Fluffy"
				})
			}
		}
	`, secretFolderUid, secretUid, secretTitle)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkSecretExistsRemotely(secretUid),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_secq", "custom.#", "1"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_secq", "custom.0.type", "securityQuestion"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_secq", "custom.0.label", "RecoveryQuestion"),
					resource.TestCheckResourceAttrSet("secretsmanager_login.custom_secq", "custom.0.value"),
				),
			},
			{
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}

// TestAccResourceLogin_customFieldCheckbox tests a checkbox field, verifying that
// the bool value round-trips correctly as "true"/"false" string in state.
func TestAccResourceLogin_customFieldCheckbox(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_custom_checkbox"

	config := fmt.Sprintf(`
		resource "secretsmanager_login" "custom_checkbox" {
			folder_uid = "%v"
			uid        = "%v"
			title      = "%v"

			custom {
				type  = "checkbox"
				label = "MFAEnabled"
				value = "true"
			}
		}
	`, secretFolderUid, secretUid, secretTitle)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkSecretExistsRemotely(secretUid),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_checkbox", "custom.#", "1"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_checkbox", "custom.0.type", "checkbox"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_checkbox", "custom.0.label", "MFAEnabled"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_checkbox", "custom.0.value", "true"),
				),
			},
			{
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}

// TestAccResourceLogin_customFieldBirthDate tests birthDate and expirationDate variants,
// which share the same []int64 millisecond structure as date but were previously missing
// from the date-to-YYYY-MM-DD conversion in the read path.
func TestAccResourceLogin_customFieldBirthDate(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_custom_birthdate"

	config := fmt.Sprintf(`
		resource "secretsmanager_login" "custom_birthdate" {
			folder_uid = "%v"
			uid        = "%v"
			title      = "%v"

			custom {
				type  = "birthDate"
				label = "DateOfBirth"
				value = "1990-05-15"
			}
			custom {
				type  = "expirationDate"
				label = "LicenseExpiry"
				value = "2027-12-31"
			}
		}
	`, secretFolderUid, secretUid, secretTitle)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkSecretExistsRemotely(secretUid),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_birthdate", "custom.#", "2"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_birthdate", "custom.0.type", "birthDate"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_birthdate", "custom.0.value", "1990-05-15"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_birthdate", "custom.1.type", "expirationDate"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_birthdate", "custom.1.value", "2027-12-31"),
				),
			},
			{
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}

// TestAccResourceLogin_customFieldHostAndBankAccount tests host and bankAccount types,
// both of which are structured object types added in the expanded type coverage pass.
func TestAccResourceLogin_customFieldHostAndBankAccount(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_custom_host_bank"

	config := fmt.Sprintf(`
		resource "secretsmanager_login" "custom_host_bank" {
			folder_uid = "%v"
			uid        = "%v"
			title      = "%v"

			custom {
				type  = "host"
				label = "PrimaryDB"
				value = jsonencode({
					hostName = "db.example.com"
					port     = "5432"
				})
			}
			custom {
				type  = "bankAccount"
				label = "Escrow"
				value = jsonencode({
					accountType   = "Checking"
					routingNumber = "021000021"
					accountNumber = "9876543210"
				})
			}
		}
	`, secretFolderUid, secretUid, secretTitle)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkSecretExistsRemotely(secretUid),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_host_bank", "custom.#", "2"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_host_bank", "custom.0.type", "host"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_host_bank", "custom.0.label", "PrimaryDB"),
					resource.TestCheckResourceAttrSet("secretsmanager_login.custom_host_bank", "custom.0.value"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_host_bank", "custom.1.type", "bankAccount"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_host_bank", "custom.1.label", "Escrow"),
					resource.TestCheckResourceAttrSet("secretsmanager_login.custom_host_bank", "custom.1.value"),
				),
			},
			{
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}

// TestAccResourceLogin_customFieldDate tests a date custom field using YYYY-MM-DD format.
// RFC3339 is also accepted on write but the read path always returns YYYY-MM-DD.
func TestAccResourceLogin_customFieldDate(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_custom_date"

	config := fmt.Sprintf(`
		resource "secretsmanager_login" "custom_date" {
			folder_uid = "%v"
			uid        = "%v"
			title      = "%v"

			custom {
				type  = "date"
				label = "ExpiresAt"
				value = "2025-06-01"
			}
		}
	`, secretFolderUid, secretUid, secretTitle)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkSecretExistsRemotely(secretUid),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_date", "custom.#", "1"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_date", "custom.0.label", "ExpiresAt"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_date", "custom.0.type", "date"),
				),
			},
		},
	})
}
