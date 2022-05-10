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
		provider "secretsmanager" {
			credential = "%v"
		}

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
				value = "MyLogin"
			}
			password {
				label = "MyPass"
				required = true
				privacy_screen = true
				enforce_generation = true
				generate = "yes"
				complexity {
					length = 20
					caps = 5
					lowercase = 5
					digits = 5
					special = 5
				}
				#value = "to_be_generated"
			}
			url {
				label = "MyUrl"
				required = true
				privacy_screen = true
				value = "https://192.168.1.1/"
			}
			card_ref {
				label = "MyCardRef"
				required = true
				privacy_screen = true
				value = "KEndJhRb5TzYShFu-zJqDA"
			}
			totp {
				label = "MyTOTP"
				required = true
				privacy_screen = true
				value = "otpauth://totp/Acme:Buster?secret=6I4PI5EUKS66GPRY5TMLJJP25MAYWAVL&issuer=Acme&algorithm=SHA1&digits=6&period=30"
			}
		}
	`, testAcc.credential, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle)

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
					resource.TestCheckResourceAttr(resourceName, "notes", secretTitle),
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
	secretTitle2 := secretTitle + "2"

	configInit := fmt.Sprintf(`
		provider "secretsmanager" {
			credential = "%v"
		}
		resource "secretsmanager_bank_account" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "%v"
		}
	`, testAcc.credential, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle)

	configUpdate := fmt.Sprintf(`
		provider "secretsmanager" {
			credential = "%v"
		}
		resource "secretsmanager_bank_account" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "%v"
		}
	`, testAcc.credential, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle2)

	resourceName := fmt.Sprintf("secretsmanager_bank_account.%v", secretTitle)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: configInit,
				Check: resource.ComposeTestCheckFunc(
					checkSecretResourceState(resourceName, func(s *terraform.InstanceState) error {
						if s.Attributes["notes"] != secretTitle {
							return fmt.Errorf("expected 'notes' = '%s'", secretTitle)
						}
						return nil
					}),
					checkSecretExistsRemotely(secretUid),
				),
			},
			{
				Config: configUpdate,
				Check: resource.ComposeTestCheckFunc(
					checkSecretResourceState(resourceName, func(s *terraform.InstanceState) error {
						if s.Attributes["notes"] != secretTitle2 {
							return fmt.Errorf("expected 'notes' = '%s'", secretTitle2)
						}
						return nil
					}),
					checkSecretExistsRemotely(secretUid),
				),
			},
		},
	})
}

/*
func TestAccResourceBankAccount_deleteDetection(t *testing.T) {
	secretType := "bankAccount"
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	_, secretTitle := testAcc.getRecordInfo(secretType)
	if secretUid == "" || secretTitle == "" {
		t.Fatal("Failed to access test data - missing secret UID and/or Title")
	}
	secretTitle += "_resource_delete"

	config := fmt.Sprintf(`
		provider "secretsmanager" {
			credential = "%v"
		}
		resource "secretsmanager_bank_account" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
		}
	`, testAcc.credential, secretTitle, secretFolderUid, secretUid, secretTitle)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				PreConfig: func() {
					// Delete secret outside of Terraform workspace
					client := *testAccProvider.Meta().(providerMeta).client
					err := client.Delete(secretUid)
					assert.OK(t, err)
				},
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true, // The externally deleted secret should be planned in for recreation
			},
		},
	})
}
*/

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
		provider "secretsmanager" {
			credential = "%v"
		}

		resource "secretsmanager_bank_account" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "%v"
		}
	`, testAcc.credential, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle)

	resourceName := fmt.Sprintf("secretsmanager_bank_account.%v", secretTitle)

	resource.Test(t, resource.TestCase{
		PreCheck:  testAccPreCheck(t),
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
