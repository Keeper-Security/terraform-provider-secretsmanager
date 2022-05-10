package secretsmanager

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/keeper-security/secrets-manager-go/core"
)

func TestAccResourceContact_create(t *testing.T) {
	secretType := "contact"
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

		resource "secretsmanager_contact" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "%v"
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
			company {
				label = "MyCompany"
				required = true
				privacy_screen = true
				value = "MyCompany"
			}
			email {
				label = "MyEmail"
				required = true
				privacy_screen = true
				value = "MyEmail"
			}
			phone {
				label = "MyPhone"
				required = true
				privacy_screen = true
				value {
					region = "US"
					number = "202-555-0130"
					ext = "9987"
					type = "Work"
				}
			}
			address_ref {
				label = "MyAddressRef"
				required = true
				privacy_screen = true
				value = "wmP8cuyXOAcJ7jflx7dgNg"
			}
		}
	`, testAcc.credential, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle)

	resourceName := fmt.Sprintf("secretsmanager_contact.%v", secretTitle)
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

func TestAccResourceContact_update(t *testing.T) {
	secretType := "contact"
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
		resource "secretsmanager_contact" "%v" {
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
		resource "secretsmanager_contact" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "%v"
		}
	`, testAcc.credential, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle2)

	resourceName := fmt.Sprintf("secretsmanager_contact.%v", secretTitle)

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
func TestAccResourceContact_deleteDetection(t *testing.T) {
	secretType := "contact"
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
		resource "secretsmanager_contact" "%v" {
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

func TestAccResourceContact_import(t *testing.T) {
	secretType := "contact"
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

		resource "secretsmanager_contact" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "%v"
		}
	`, testAcc.credential, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle)

	resourceName := fmt.Sprintf("secretsmanager_contact.%v", secretTitle)

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
