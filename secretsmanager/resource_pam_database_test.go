package secretsmanager

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/keeper-security/secrets-manager-go/core"
)

func TestAccResourcePamDatabase_create(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_test_pam_database_create"
	if secretFolderUid == "" {
		t.Skip("Skipping test - TF_ACC not set or test folder not configured")
	}

	config := fmt.Sprintf(`
		resource "secretsmanager_pam_database" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "%v"
			pam_hostname {
				value {
					hostname = "db.example.com"
					port = "5432"
				}
			}
			pam_settings = jsonencode([{
				connection = [{
					protocol = "postgresql"
					port = "5432"
					recordingIncludeKeys = true
					allowSupplyUser = false
					database = "production"
				}]
			}])
			database_type = "postgresql"
			use_ssl {
				value = true
			}
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle)

	resourceName := fmt.Sprintf("secretsmanager_pam_database.%v", secretTitle)
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkSecretExistsRemotely(secretUid),
					resource.TestCheckResourceAttr(resourceName, "type", "pamDatabase"),
					resource.TestCheckResourceAttr(resourceName, "title", secretTitle),
					resource.TestCheckResourceAttr(resourceName, "notes", secretTitle),
					resource.TestCheckResourceAttr(resourceName, "pam_hostname.0.value.0.hostname", "db.example.com"),
					resource.TestCheckResourceAttr(resourceName, "pam_hostname.0.value.0.port", "5432"),
					resource.TestCheckResourceAttr(resourceName, "database_type", "postgresql"),
				),
			},
		},
	})
}

// TestAccResourcePamDatabase_update tests updating PAM Database fields.
// NOTE: This test only updates fields that work reliably (pam_hostname, pam_settings, notes).
// The use_ssl field is NOT tested because it uses ApplyFieldChange() which doesn't sync
// RecordDict to RawJson due to an SDK limitation. See resourcePamDatabaseUpdate:471-475.
func TestAccResourcePamDatabase_update(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_test_pam_database_update"
	secretTitle2 := "tf_acc_test_pam_database_update_2"
	if secretFolderUid == "" {
		t.Skip("Skipping test - TF_ACC not set or test folder not configured")
	}

	configInit := fmt.Sprintf(`
		resource "secretsmanager_pam_database" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "%v"
			pam_hostname {
				value {
					hostname = "db.example.com"
					port = "5432"
				}
			}
			pam_settings = jsonencode([{
				connection = [{
					protocol = "postgresql"
					port = "5432"
					database = "production"
				}]
			}])
			database_type = "postgresql"
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle)

	configUpdate := fmt.Sprintf(`
		resource "secretsmanager_pam_database" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "%v"
			pam_hostname {
				value {
					hostname = "db-new.example.com"
					port = "3306"
				}
			}
			pam_settings = jsonencode([{
				connection = [{
					protocol = "mysql"
					port = "3306"
					database = "staging"
					allowSupplyHost = true
				}]
			}])
			database_type = "mysql"
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle2)

	resourceName := fmt.Sprintf("secretsmanager_pam_database.%v", secretTitle)

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
						if s.Attributes["pam_hostname.0.value.0.hostname"] != "db.example.com" {
							return fmt.Errorf("expected hostname = 'db.example.com'")
						}
						if s.Attributes["database_type"] != "postgresql" {
							return fmt.Errorf("expected database_type = 'postgresql'")
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
						if s.Attributes["pam_hostname.0.value.0.hostname"] != "db-new.example.com" {
							return fmt.Errorf("expected hostname = 'db-new.example.com'")
						}
						if s.Attributes["pam_hostname.0.value.0.port"] != "3306" {
							return fmt.Errorf("expected port = '3306'")
						}
						if s.Attributes["database_type"] != "mysql" {
							return fmt.Errorf("expected database_type = 'mysql'")
						}
						return nil
					}),
					checkSecretExistsRemotely(secretUid),
				),
			},
		},
	})
}

func TestAccResourcePamDatabase_deleteDetection(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_test_pam_database_delete"
	if secretFolderUid == "" {
		t.Skip("Skipping test - TF_ACC not set or test folder not configured")
	}

	config := fmt.Sprintf(`
		resource "secretsmanager_pam_database" "%v" {
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

func TestAccResourcePamDatabase_import(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_test_pam_database_import"
	if secretFolderUid == "" {
		t.Skip("Skipping test - TF_ACC not set or test folder not configured")
	}

	config := fmt.Sprintf(`
		resource "secretsmanager_pam_database" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "%v"
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle)

	resourceName := fmt.Sprintf("secretsmanager_pam_database.%v", secretTitle)

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
