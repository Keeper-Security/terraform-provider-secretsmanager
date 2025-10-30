package secretsmanager

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/keeper-security/secrets-manager-go/core"
)

func TestAccResourcePamDatabase_create(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_test_pam_database_create"
	if secretFolderUid == "" {
		t.Fatal("Failed to access test folder UID")
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
				value = [true]
			}
			login {
				value = "dbadmin"
			}
			password {
				enforce_generation = true
				generate = "yes"
				complexity {
					length = 32
					caps = 8
					lowercase = 8
					digits = 8
					special = 8
				}
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
				),
			},
		},
	})
}

// TestAccResourcePamDatabase_update is disabled due to SDK limitation.
//
// ISSUE: Updates to database_type and use_ssl fields don't persist to Keeper.
// ROOT CAUSE: SDK's InsertField() and UpdateField() modify RecordDict but don't
// call the private update() method to sync changes to RawJson. Since Save() uses
// RawJson (not RecordDict), the changes are lost.
//
// WORKAROUND: pam_hostname and pam_settings use SetStandardFieldValue() which
// calls update() internally, but this doesn't work for fields that need label-based
// matching (multiple checkboxes, multiple databaseType fields, etc.).
//
// TODO: Enable this test once SDK either:
// 1. Makes update() method public, OR
// 2. Has InsertField/UpdateField call update() internally
//
// func TestAccResourcePamDatabase_update(t *testing.T) {
// 	secretFolderUid := testAcc.getTestFolder()
// 	secretUid := core.GenerateUid()
// 	secretTitle := "tf_acc_test_pam_database_update"
// 	secretTitle2 := "tf_acc_test_pam_database_update_2"
// 	if secretFolderUid == "" {
// 		t.Fatal("Failed to access test folder UID")
// 	}
//
// 	configInit := fmt.Sprintf(`
// 		resource "secretsmanager_pam_database" "%v" {
// 			folder_uid = "%v"
// 			uid = "%v"
// 			title = "%v"
// 			notes = "%v"
// 			pam_hostname {
// 				value {
// 					hostname = "db.example.com"
// 					port = "5432"
// 				}
// 			}
// 			pam_settings = jsonencode([{
// 				connection = [{
// 					protocol = "postgresql"
// 					port = "5432"
// 					database = "production"
// 				}]
// 			}])
// 			database_type {
// 				value = ["PostgreSQL"]
// 			}
// 		}
// 	`, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle)
//
// 	configUpdate := fmt.Sprintf(`
// 		resource "secretsmanager_pam_database" "%v" {
// 			folder_uid = "%v"
// 			uid = "%v"
// 			title = "%v"
// 			notes = "%v"
// 			pam_hostname {
// 				value {
// 					hostname = "db-new.example.com"
// 					port = "3306"
// 				}
// 			}
// 			pam_settings = jsonencode([{
// 				connection = [{
// 					protocol = "mysql"
// 					port = "3306"
// 					database = "staging"
// 					allowSupplyHost = true
// 				}]
// 			}])
// 			database_type {
// 				value = ["MySQL"]
// 			}
// 			use_ssl {
// 				value = [true]
// 			}
// 		}
// 	`, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle2)
//
// 	resourceName := fmt.Sprintf("secretsmanager_pam_database.%v", secretTitle)
//
// 	resource.Test(t, resource.TestCase{
// 		Providers: testAccProviders,
// 		PreCheck:  testAccPreCheck(t),
// 		Steps: []resource.TestStep{
// 			{
// 				Config: configInit,
// 				Check: resource.ComposeTestCheckFunc(
// 					checkSecretResourceState(resourceName, func(s *terraform.InstanceState) error {
// 						if s.Attributes["notes"] != secretTitle {
// 							return fmt.Errorf("expected 'notes' = '%s'", secretTitle)
// 						}
// 						if s.Attributes["pam_hostname.0.value.0.hostname"] != "db.example.com" {
// 							return fmt.Errorf("expected hostname = 'db.example.com'")
// 						}
// 						return nil
// 					}),
// 					checkSecretExistsRemotely(secretUid),
// 				),
// 			},
// 			{
// 				Config: configUpdate,
// 				Check: resource.ComposeTestCheckFunc(
// 					checkSecretResourceState(resourceName, func(s *terraform.InstanceState) error {
// 						if s.Attributes["notes"] != secretTitle2 {
// 							return fmt.Errorf("expected 'notes' = '%s'", secretTitle2)
// 						}
// 						if s.Attributes["pam_hostname.0.value.0.hostname"] != "db-new.example.com" {
// 							return fmt.Errorf("expected hostname = 'db-new.example.com'")
// 						}
// 						if s.Attributes["pam_hostname.0.value.0.port"] != "3306" {
// 							return fmt.Errorf("expected port = '3306'")
// 						}
// 						return nil
// 					}),
// 					checkSecretExistsRemotely(secretUid),
// 				),
// 			},
// 		},
// 	})
// }

func TestAccResourcePamDatabase_deleteDetection(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_test_pam_database_delete"
	if secretFolderUid == "" {
		t.Fatal("Failed to access test folder UID")
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
		t.Fatal("Failed to access test folder UID")
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
