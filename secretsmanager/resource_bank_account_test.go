package secretsmanager

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/keeper-security/secrets-manager-go/core"
)

func TestAccResourceBankAccount_create(t *testing.T) {
	secretType := "bankAccount"
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	_, secretTitle := testAcc.getRecordInfo(secretType)
	if secretUid == "" || secretTitle == "" {
		t.Fatal("Failed to access test data - missing secret UID and/or Title")
	}
	secretTitle += "_resource_create"

	config := fmt.Sprintf(`
		resource "secretsmanager_bank_account" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "%v"
			bank_account {
				label = "My Bank Account"
				required = true
				privacy_screen = true
				value {
					account_type = "Other"
					routing_number = "1234567890"
					account_number = "0987654321"
					other_type = "Investment"
				}
			}
			name {
				label = "John"
				required = true
				privacy_screen = true
				value {
					first = "John"
					middle = "D"
					last = "Doe"
				}
			}
			login {
				label = "MyLogin"
				required = true
				privacy_screen = true
				value = "john.doe@example.com"
			}
			password {
				label = "MyPassword"
				required = true
				privacy_screen = true
				value = "ThisIsAStrongPassword123!"
			}
			url {
				label = "MyURL"
				required = true
				privacy_screen = true
				value = "https://mybank.example.com"
			}
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle, "A longer note field")

	resourceName := fmt.Sprintf("secretsmanager_bank_account.%v", secretTitle)
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkSecretExistsRemotely(secretUid),
					resource.TestCheckResourceAttr(resourceName, "type", secretType),
					resource.TestCheckResourceAttr(resourceName, "title", secretTitle),
					resource.TestCheckResourceAttr(resourceName, "bank_account.0.label", "My Bank Account"),
					resource.TestCheckResourceAttr(resourceName, "bank_account.0.required", "true"),
					resource.TestCheckResourceAttr(resourceName, "bank_account.0.privacy_screen", "true"),
					resource.TestCheckResourceAttr(resourceName, "bank_account.0.value.0.account_type", "Other"),
					resource.TestCheckResourceAttr(resourceName, "bank_account.0.value.0.routing_number", "1234567890"),
					resource.TestCheckResourceAttr(resourceName, "bank_account.0.value.0.account_number", "0987654321"),
					resource.TestCheckResourceAttr(resourceName, "bank_account.0.value.0.other_type", "Investment"),
				),
			},
		},
	})
}

func TestAccResourceBankAccount_update(t *testing.T) {
	secretType := "bankAccount"
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	_, secretTitle := testAcc.getRecordInfo(secretType)
	if secretUid == "" || secretTitle == "" {
		t.Fatal("Failed to access test data - missing secret UID and/or Title")
	}
	secretTitle += "_resource_update"

	config := fmt.Sprintf(`
		resource "secretsmanager_bank_account" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "%v"
			bank_account {
				label = "My Bank Account"
				required = true
				privacy_screen = true
				value {
					account_type = "Other"
					routing_number = "1234567890"
					account_number = "0987654321"
					other_type = "Investment"
				}
			}
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle, "A longer note field")

	configUpdate := fmt.Sprintf(`
		resource "secretsmanager_bank_account" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "%v"
			bank_account {
				label = "My Bank Account"
				required = true
				privacy_screen = true
				value {
					account_type = "Other"
					routing_number = "99999999"
					account_number = "11111111"
					other_type = "Checking"
				}
			}
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle+"_updated", "Updated note field")

	resourceName := fmt.Sprintf("secretsmanager_bank_account.%v", secretTitle)
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkSecretExistsRemotely(secretUid),
					resource.TestCheckResourceAttr(resourceName, "type", secretType),
					resource.TestCheckResourceAttr(resourceName, "title", secretTitle),
					resource.TestCheckResourceAttr(resourceName, "bank_account.0.value.0.routing_number", "1234567890"),
				),
			},
			{
				Config: configUpdate,
				Check: resource.ComposeTestCheckFunc(
					checkSecretExistsRemotely(secretUid),
					resource.TestCheckResourceAttr(resourceName, "title", secretTitle+"_updated"),
					resource.TestCheckResourceAttr(resourceName, "bank_account.0.value.0.routing_number", "99999999"),
					resource.TestCheckResourceAttr(resourceName, "bank_account.0.value.0.account_number", "11111111"),
					resource.TestCheckResourceAttr(resourceName, "bank_account.0.value.0.other_type", "Checking"),
				),
			},
		},
	})
}

func TestAccResourceBankAccount_import(t *testing.T) {
	secretType := "bankAccount"
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	_, secretTitle := testAcc.getRecordInfo(secretType)
	if secretUid == "" || secretTitle == "" {
		t.Fatal("Failed to access test data - missing secret UID and/or Title")
	}
	secretTitle += "_resource_import"

	config := fmt.Sprintf(`
		resource "secretsmanager_bank_account" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			bank_account {
				value {
					routing_number = "1234567890"
					account_number = "0987654321"
				}
			}
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle)

	resourceName := fmt.Sprintf("secretsmanager_bank_account.%v", secretTitle)
	checkFn := func(s *terraform.State) error {
		return checkSecretExistsRemotely(secretUid)(s)
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkFn,
					resource.TestCheckResourceAttr(resourceName, "type", secretType),
					resource.TestCheckResourceAttr(resourceName, "title", secretTitle),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     secretUid,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccResourceBankAccount_customFields tests custom field support for bank_account resources.
// This test is representative of all resources with complex nested value structures, including:
//   - bank_account (this resource)
//   - bank_card
//   - payment_card
//   - contact
//   - health_insurance
//   - membership
//   - software_license
//
// Custom field functionality is identical across all resource types - this test validates
// the pattern works correctly with complex nested field structures.
func TestAccResourceBankAccount_customFields(t *testing.T) {
	secretType := "bankAccount"
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	_, secretTitle := testAcc.getRecordInfo(secretType)
	if secretUid == "" || secretTitle == "" {
		t.Fatal("Failed to access test data - missing secret UID and/or Title")
	}
	secretTitle += "_resource_custom_fields"

	configCreate := fmt.Sprintf(`
		resource "secretsmanager_bank_account" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "Test custom fields"
			bank_account {
				value {
					routing_number = "1234567890"
					account_number = "0987654321"
				}
			}
			custom {
				type = "text"
				label = "Account Manager"
				value = "Jane Smith"
			}
			custom {
				type = "text"
				label = "Branch Code"
				value = "BR001"
			}
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle)

	configUpdate := fmt.Sprintf(`
		resource "secretsmanager_bank_account" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "Test custom fields updated"
			bank_account {
				value {
					routing_number = "1234567890"
					account_number = "0987654321"
				}
			}
			custom {
				type = "text"
				label = "Account Manager"
				value = "John Doe"
			}
			custom {
				type = "text"
				label = "Branch Code"
				value = "BR002"
			}
			custom {
				type = "text"
				label = "Region"
				value = "West Coast"
			}
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle)

	resourceName := fmt.Sprintf("secretsmanager_bank_account.%v", secretTitle)
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: configCreate,
				Check: resource.ComposeTestCheckFunc(
					checkSecretExistsRemotely(secretUid),
					resource.TestCheckResourceAttr(resourceName, "type", secretType),
					resource.TestCheckResourceAttr(resourceName, "title", secretTitle),
					resource.TestCheckResourceAttr(resourceName, "custom.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "custom.0.label", "Account Manager"),
					resource.TestCheckResourceAttr(resourceName, "custom.0.value", "Jane Smith"),
					resource.TestCheckResourceAttr(resourceName, "custom.1.label", "Branch Code"),
					resource.TestCheckResourceAttr(resourceName, "custom.1.value", "BR001"),
				),
			},
			{
				Config: configUpdate,
				Check: resource.ComposeTestCheckFunc(
					checkSecretExistsRemotely(secretUid),
					resource.TestCheckResourceAttr(resourceName, "custom.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "custom.0.label", "Account Manager"),
					resource.TestCheckResourceAttr(resourceName, "custom.0.value", "John Doe"),
					resource.TestCheckResourceAttr(resourceName, "custom.1.label", "Branch Code"),
					resource.TestCheckResourceAttr(resourceName, "custom.1.value", "BR002"),
					resource.TestCheckResourceAttr(resourceName, "custom.2.label", "Region"),
					resource.TestCheckResourceAttr(resourceName, "custom.2.value", "West Coast"),
				),
			},
		},
	})
}
