package secretsmanager

import (
	"encoding/json"
	"fmt"
	"regexp"
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

// TestAccResourceLogin_customFieldSimpleVariants tests simple string types that share
// the core.Text write path: url, email, and multiline.
func TestAccResourceLogin_customFieldSimpleVariants(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_custom_simple_variants"

	config := fmt.Sprintf(`
		resource "secretsmanager_login" "custom_simple_variants" {
			folder_uid = "%v"
			uid        = "%v"
			title      = "%v"

			custom {
				type  = "url"
				label = "Dashboard"
				value = "https://example.com/dashboard"
			}
			custom {
				type  = "email"
				label = "Alerts"
				value = "alerts@example.com"
			}
			custom {
				type  = "multiline"
				label = "Notes"
				value = "line one\nline two"
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
					resource.TestCheckResourceAttr("secretsmanager_login.custom_simple_variants", "custom.#", "3"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_simple_variants", "custom.0.type", "url"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_simple_variants", "custom.0.value", "https://example.com/dashboard"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_simple_variants", "custom.1.type", "email"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_simple_variants", "custom.1.value", "alerts@example.com"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_simple_variants", "custom.2.type", "multiline"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_simple_variants", "custom.2.value", "line one\nline two"),
				),
			},
			{
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}

// TestAccResourceLogin_customFieldName tests the name complex type (first/middle/last).
func TestAccResourceLogin_customFieldName(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_custom_name"

	config := fmt.Sprintf(`
		resource "secretsmanager_login" "custom_name" {
			folder_uid = "%v"
			uid        = "%v"
			title      = "%v"

			custom {
				type  = "name"
				label = "Owner"
				value = jsonencode({
					first  = "Jane"
					middle = "Q"
					last   = "Doe"
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
					resource.TestCheckResourceAttr("secretsmanager_login.custom_name", "custom.#", "1"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_name", "custom.0.type", "name"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_name", "custom.0.label", "Owner"),
					resource.TestCheckResourceAttrSet("secretsmanager_login.custom_name", "custom.0.value"),
				),
			},
			{
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}

// TestAccResourceLogin_customFieldAddress tests the address complex type.
func TestAccResourceLogin_customFieldAddress(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_custom_address"

	config := fmt.Sprintf(`
		resource "secretsmanager_login" "custom_address" {
			folder_uid = "%v"
			uid        = "%v"
			title      = "%v"

			custom {
				type  = "address"
				label = "HQ"
				value = jsonencode({
					street1 = "123 Main St"
					street2 = "Suite 400"
					city    = "San Francisco"
					state   = "CA"
					country = "US"
					zip     = "94105"
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
					resource.TestCheckResourceAttr("secretsmanager_login.custom_address", "custom.#", "1"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_address", "custom.0.type", "address"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_address", "custom.0.label", "HQ"),
					resource.TestCheckResourceAttrSet("secretsmanager_login.custom_address", "custom.0.value"),
				),
			},
			{
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}

// TestAccResourceLogin_customFieldPaymentCardRoundTrip is the regression test for KSM-888:
// paymentCard keys must use camelCase (cardNumber, cardExpirationDate, cardSecurityCode)
// to match what the read path serializes via toJsonDefault() / json.Marshal on the SDK
// struct. The original write path read snake_case keys (card_number, card_expiration_date,
// card_security_code), which never matched any real key in the user's jsonencode output,
// causing the card to be written empty and a perpetual plan diff on every apply.
//
// Failure mode before fix: PlanOnly step detects diff (state "{}"; config has card data).
// Passing mode after fix: both write and read use camelCase; round-trip is stable.
func TestAccResourceLogin_customFieldPaymentCardRoundTrip(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_custom_payment_card_rt"

	config := fmt.Sprintf(`
		resource "secretsmanager_login" "custom_payment_card_rt" {
			folder_uid = "%v"
			uid        = "%v"
			title      = "%v"

			custom {
				type  = "paymentCard"
				label = "Corporate Card"
				value = jsonencode({
					cardNumber         = "4111111111111111"
					cardExpirationDate = "12/2027"
					cardSecurityCode   = "123"
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
					resource.TestCheckResourceAttr("secretsmanager_login.custom_payment_card_rt", "custom.#", "1"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_payment_card_rt", "custom.0.type", "paymentCard"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_payment_card_rt", "custom.0.label", "Corporate Card"),
					resource.TestCheckResourceAttrSet("secretsmanager_login.custom_payment_card_rt", "custom.0.value"),
				),
			},
			{
				// Regression guard: state must equal config — no perpetual diff.
				// Before fix: write path read card_number (snake_case, not found) →
				// wrote empty card → read back "{}" → diff on every plan.
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}

// TestAccResourceLogin_customFieldDate tests a date custom field using YYYY-MM-DD format.
// Only YYYY-MM-DD is accepted — RFC3339 is rejected with a clear error (KSM-889).
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

// TestAccResourceLogin_customFieldImportAndComputedMetadata tests two related things:
//
//  1. Import: a login resource with custom fields that have required=true and
//     privacy_screen=true can be imported and all fields round-trip correctly.
//
//  2. No perpetual diff: after import (or create), a config that OMITS required
//     and privacy_screen produces no plan diff. This is the regression guard for
//     the Computed: true fix — without it, every plan would show a diff because
//     Terraform would treat the omitted field as "set to false" vs vault's "true".
func TestAccResourceLogin_customFieldImportAndComputedMetadata(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_custom_import"
	resourceName := "secretsmanager_login.custom_import"

	// configExplicit sets required and privacy_screen explicitly.
	// This is what creates the vault record with those flags set.
	configExplicit := fmt.Sprintf(`
		resource "secretsmanager_login" "custom_import" {
			folder_uid = "%v"
			uid        = "%v"
			title      = "%v"

			custom {
				type           = "text"
				label          = "Environment"
				value          = "production"
				required       = true
				privacy_screen = true
			}
		}
	`, secretFolderUid, secretUid, secretTitle)

	// configOmitted is identical but omits required and privacy_screen entirely.
	// With Computed: true, Terraform should accept the vault values and show no diff.
	configOmitted := fmt.Sprintf(`
		resource "secretsmanager_login" "custom_import" {
			folder_uid = "%v"
			uid        = "%v"
			title      = "%v"

			custom {
				type  = "text"
				label = "Environment"
				value = "production"
			}
		}
	`, secretFolderUid, secretUid, secretTitle)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				// Step 1: create with explicit required/privacy_screen
				Config: configExplicit,
				Check: resource.ComposeTestCheckFunc(
					checkSecretExistsRemotely(secretUid),
					resource.TestCheckResourceAttr(resourceName, "custom.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "custom.0.type", "text"),
					resource.TestCheckResourceAttr(resourceName, "custom.0.label", "Environment"),
					resource.TestCheckResourceAttr(resourceName, "custom.0.required", "true"),
					resource.TestCheckResourceAttr(resourceName, "custom.0.privacy_screen", "true"),
				),
			},
			{
				// Step 2: import by UID — verify all fields including required and
				// privacy_screen round-trip correctly from the vault.
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// value is Sensitive: true — Terraform stores it in state but masks
				// it in output. ImportStateVerify compares raw state values, so it
				// should still match. If the vault returns the value in a different
				// format after import, add it to ImportStateVerifyIgnore.
			},
			{
				// Step 3: PlanOnly with config that OMITS required and privacy_screen.
				// This is the perpetual-diff regression guard: if Computed: true is
				// working correctly, omitting these fields in HCL should produce no
				// diff (vault values are authoritative). If Computed were missing,
				// Terraform would want to set them to false on every plan.
				Config:   configOmitted,
				PlanOnly: true,
			},
		},
	})
}

// TestAccResourceLogin_customFieldKeyPair tests the keyPair complex type.
// A real ed25519 key pair is generated in-process and injected into the config
// so we verify actual key material round-trips correctly through the vault.
func TestAccResourceLogin_customFieldKeyPair(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_custom_keypair"

	kp, err := GenerateSSHKeyPair(SSHKeyTypeED25519, 0, "")
	if err != nil {
		t.Fatalf("failed to generate key pair for test: %v", err)
	}

	// Marshal to JSON so newlines and special chars are properly escaped
	// for embedding as an HCL string literal.
	valueJSON, err := json.Marshal(map[string]string{
		"publicKey":  kp.PublicKey,
		"privateKey": kp.PrivateKey,
	})
	if err != nil {
		t.Fatalf("failed to marshal key pair to JSON: %v", err)
	}

	config := fmt.Sprintf(`
		resource "secretsmanager_login" "custom_keypair" {
			folder_uid = "%v"
			uid        = "%v"
			title      = "%v"

			custom {
				type  = "keyPair"
				label = "DeployKey"
				value = %q
			}
		}
	`, secretFolderUid, secretUid, secretTitle, string(valueJSON))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkSecretExistsRemotely(secretUid),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_keypair", "custom.#", "1"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_keypair", "custom.0.type", "keyPair"),
					resource.TestCheckResourceAttr("secretsmanager_login.custom_keypair", "custom.0.label", "DeployKey"),
					resource.TestCheckResourceAttrSet("secretsmanager_login.custom_keypair", "custom.0.value"),
				),
			},
			{
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}

// TestAccResourceLogin_customFieldCheckboxInvalidInput is the regression test for KSM-889
// (checkbox). Values like "yes", "1", "on" were silently treated as false — because
// strings.ToLower(value) == "true" is false for them — then read back as "false", causing
// a perpetual diff between config and state.
//
// After KSM-889: apply returns an error immediately, directing the user to use "true" or "false".
// Failure mode before fix: no error raised → ExpectError step fails (expected error, got none).
// Passing mode after fix: error raised → ExpectError step passes.
func TestAccResourceLogin_customFieldCheckboxInvalidInput(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_custom_checkbox_invalid"

	config := fmt.Sprintf(`
		resource "secretsmanager_login" "custom_checkbox_invalid" {
			folder_uid = "%v"
			uid        = "%v"
			title      = "%v"

			custom {
				type  = "checkbox"
				label = "Feature"
				value = "yes"
			}
		}
	`, secretFolderUid, secretUid, secretTitle)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`invalid checkbox value`),
			},
		},
	})
}

// TestAccResourceLogin_customFieldDateInvalidFormat is the regression test for KSM-889
// (date). RFC3339 values like "2026-03-20T14:30:00Z" were silently accepted, stored as
// epoch ms, but the read path always returned "2026-03-20" (YYYY-MM-DD only, time
// component discarded). This caused a perpetual diff: config kept the RFC3339 string
// while state showed YYYY-MM-DD.
//
// After KSM-889: apply returns an error immediately, directing the user to use YYYY-MM-DD.
// Applies equally to date, birthDate, and expirationDate — all share the same write path.
// Failure mode before fix: no error raised → ExpectError step fails (expected error, got none).
// Passing mode after fix: error raised → ExpectError step passes.
func TestAccResourceLogin_customFieldDateInvalidFormat(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_custom_date_invalid"

	config := fmt.Sprintf(`
		resource "secretsmanager_login" "custom_date_invalid" {
			folder_uid = "%v"
			uid        = "%v"
			title      = "%v"

			custom {
				type  = "date"
				label = "Review"
				value = "2026-03-20T14:30:00Z"
			}
		}
	`, secretFolderUid, secretUid, secretTitle)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`invalid date value`),
			},
		},
	})
}
