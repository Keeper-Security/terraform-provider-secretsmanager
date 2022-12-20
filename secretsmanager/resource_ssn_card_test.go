package secretsmanager

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/keeper-security/secrets-manager-go/core"
)

func TestAccResourceSsnCard_create(t *testing.T) {
	secretType := "ssnCard"
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	_, secretTitle := testAcc.getRecordInfo(secretType)
	if secretUid == "" || secretTitle == "" {
		t.Fatal("Failed to access test data - missing secret UID and/or Title")
	}
	secretTitle += "_resource_create"

	config := fmt.Sprintf(`
		resource "secretsmanager_ssn_card" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "%v"
			identity_number {
				label = "MyAccount"
				required = true
				privacy_screen = true
				value = "Account1"
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
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle)

	resourceName := fmt.Sprintf("secretsmanager_ssn_card.%v", secretTitle)
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

func TestAccResourceSsnCard_update(t *testing.T) {
	secretType := "ssnCard"
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	_, secretTitle := testAcc.getRecordInfo(secretType)
	if secretUid == "" || secretTitle == "" {
		t.Fatal("Failed to access test data - missing secret UID and/or Title")
	}
	secretTitle += "_resource_update"
	secretTitle2 := secretTitle + "2"

	configInit := fmt.Sprintf(`
		resource "secretsmanager_ssn_card" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "%v"
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle)

	configUpdate := fmt.Sprintf(`
		resource "secretsmanager_ssn_card" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "%v"
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle2)

	resourceName := fmt.Sprintf("secretsmanager_ssn_card.%v", secretTitle)

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

func TestAccResourceSsnCard_deleteDetection(t *testing.T) {
	secretType := "ssnCard"
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	_, secretTitle := testAcc.getRecordInfo(secretType)
	if secretUid == "" || secretTitle == "" {
		t.Fatal("Failed to access test data - missing secret UID and/or Title")
	}
	secretTitle += "_resource_delete"

	config := fmt.Sprintf(`
		resource "secretsmanager_ssn_card" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle)

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
					if err := deleteRecord(secretUid, client); err != nil {
						t.Fail()
					}
				},
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true, // The externally deleted secret should be planned in for recreation
			},
		},
	})
}

func TestAccResourceSsnCard_import(t *testing.T) {
	secretType := "ssnCard"
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	_, secretTitle := testAcc.getRecordInfo(secretType)
	if secretUid == "" || secretTitle == "" {
		t.Fatal("Failed to access test data - missing secret UID and/or Title")
	}
	secretTitle += "_resource_import"

	config := fmt.Sprintf(`
		resource "secretsmanager_ssn_card" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "%v"
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle)

	resourceName := fmt.Sprintf("secretsmanager_ssn_card.%v", secretTitle)

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
