package secretsmanager

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/keeper-security/secrets-manager-go/core"
)

func TestAccResourcePamMachine_create(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_test_pam_machine_create"
	if secretFolderUid == "" {
		t.Fatal("Failed to access test folder UID")
	}

	config := fmt.Sprintf(`
		resource "secretsmanager_pam_machine" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "%v"
			pam_hostname {
				value {
					hostname = "192.168.1.100"
					port = "22"
				}
			}
			pam_settings = jsonencode([{
				connection = [{
					protocol = "ssh"
					port = "22"
					recordingIncludeKeys = true
					colorScheme = "green_black"
					allowSupplyUser = false
				}]
				portForward = [{
					port = "2222"
					reusePort = true
				}]
			}])
			login {
				value = "admin"
			}
			password {
				enforce_generation = true
				generate = "yes"
				complexity {
					length = 20
					caps = 5
					lowercase = 5
					digits = 5
					special = 5
				}
			}
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle)

	resourceName := fmt.Sprintf("secretsmanager_pam_machine.%v", secretTitle)
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkSecretExistsRemotely(secretUid),
					resource.TestCheckResourceAttr(resourceName, "type", "pamMachine"),
					resource.TestCheckResourceAttr(resourceName, "title", secretTitle),
					resource.TestCheckResourceAttr(resourceName, "notes", secretTitle),
					resource.TestCheckResourceAttr(resourceName, "pam_hostname.0.value.0.hostname", "192.168.1.100"),
					resource.TestCheckResourceAttr(resourceName, "pam_hostname.0.value.0.port", "22"),
				),
			},
		},
	})
}

func TestAccResourcePamMachine_update(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_test_pam_machine_update"
	secretTitle2 := "tf_acc_test_pam_machine_update_2"
	if secretFolderUid == "" {
		t.Fatal("Failed to access test folder UID")
	}

	configInit := fmt.Sprintf(`
		resource "secretsmanager_pam_machine" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "%v"
			pam_hostname {
				value {
					hostname = "192.168.1.100"
					port = "22"
				}
			}
			pam_settings = jsonencode([{
				connection = [{
					protocol = "ssh"
					port = "22"
				}]
			}])
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle)

	configUpdate := fmt.Sprintf(`
		resource "secretsmanager_pam_machine" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "%v"
			pam_hostname {
				value {
					hostname = "192.168.1.200"
					port = "2222"
				}
			}
			pam_settings = jsonencode([{
				connection = [{
					protocol = "rdp"
					port = "3389"
					security = "nla"
					ignoreCert = true
				}]
			}])
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle2)

	resourceName := fmt.Sprintf("secretsmanager_pam_machine.%v", secretTitle)

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
						if s.Attributes["pam_hostname.0.value.0.hostname"] != "192.168.1.100" {
							return fmt.Errorf("expected hostname = '192.168.1.100'")
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
						actualHostname := s.Attributes["pam_hostname.0.value.0.hostname"]
						if actualHostname != "192.168.1.200" {
							return fmt.Errorf("expected hostname = '192.168.1.200', got '%s'", actualHostname)
						}
						actualPort := s.Attributes["pam_hostname.0.value.0.port"]
						if actualPort != "2222" {
							return fmt.Errorf("expected port = '2222', got '%s'", actualPort)
						}
						return nil
					}),
					checkSecretExistsRemotely(secretUid),
				),
			},
		},
	})
}

func TestAccResourcePamMachine_deleteDetection(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_test_pam_machine_delete"
	if secretFolderUid == "" {
		t.Fatal("Failed to access test folder UID")
	}

	config := fmt.Sprintf(`
		resource "secretsmanager_pam_machine" "%v" {
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

func TestAccResourcePamMachine_import(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_test_pam_machine_import"
	if secretFolderUid == "" {
		t.Fatal("Failed to access test folder UID")
	}

	config := fmt.Sprintf(`
		resource "secretsmanager_pam_machine" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "%v"
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle)

	resourceName := fmt.Sprintf("secretsmanager_pam_machine.%v", secretTitle)

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
