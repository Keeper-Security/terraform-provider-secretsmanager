package secretsmanager

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/keeper-security/secrets-manager-go/core"
)

func TestAccResourcePamDirectory_create(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_test_pam_directory_create"
	if secretFolderUid == "" {
		t.Skip("Skipping test - TF_ACC not set or test folder not configured")
	}

	config := fmt.Sprintf(`
		resource "secretsmanager_pam_directory" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "%v"
			pam_hostname {
				value {
					hostname = "ad.corp.example.com"
					port = "636"
				}
			}
			pam_settings = jsonencode([{
				connection = [{
					protocol = "ldaps"
					port = "636"
					recordingIncludeKeys = false
				}]
			}])
			directory_type = "Active Directory"
			login {
				value = "CN=Admin,CN=Users,DC=corp,DC=example,DC=com"
			}
			password {
				enforce_generation = true
				generate = "yes"
				complexity {
					length = 24
					caps = 6
					lowercase = 6
					digits = 6
					special = 6
				}
			}
			distinguished_name {
				label = "Distinguished Name"
				value = "DC=corp,DC=example,DC=com"
			}
		domain_name {
			label = "domainName"
			value = "corp.example.com"
		}
		directory_id {
			label = "directoryId"
			value = "dir-12345678"
		}
		user_match {
			label = "userMatch"
			value = "sAMAccountName"
		}
		provider_group {
			label = "providerGroup"
			value = "prod-ad-servers"
		}
		provider_region {
			label = "providerRegion"
			value = "us-west-2"
		}
		alternative_ips {
			label = "alternativeIPs"
			value = "10.0.1.5\n10.0.1.6\n10.0.1.7"
		}
			use_ssl {
				value = [true]
			}
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle)

	resourceName := fmt.Sprintf("secretsmanager_pam_directory.%v", secretTitle)
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkSecretExistsRemotely(secretUid),
					resource.TestCheckResourceAttr(resourceName, "type", "pamDirectory"),
					resource.TestCheckResourceAttr(resourceName, "title", secretTitle),
					resource.TestCheckResourceAttr(resourceName, "notes", secretTitle),
					resource.TestCheckResourceAttr(resourceName, "pam_hostname.0.value.0.hostname", "ad.corp.example.com"),
					resource.TestCheckResourceAttr(resourceName, "pam_hostname.0.value.0.port", "636"),
					resource.TestCheckResourceAttr(resourceName, "directory_type", "Active Directory"),
					resource.TestCheckResourceAttr(resourceName, "login.0.value", "CN=Admin,CN=Users,DC=corp,DC=example,DC=com"),
					resource.TestCheckResourceAttr(resourceName, "distinguished_name.0.value", "DC=corp,DC=example,DC=com"),
					resource.TestCheckResourceAttr(resourceName, "domain_name.0.value", "corp.example.com"),
					resource.TestCheckResourceAttr(resourceName, "directory_id.0.value", "dir-12345678"),
					resource.TestCheckResourceAttr(resourceName, "user_match.0.value", "sAMAccountName"),
					resource.TestCheckResourceAttr(resourceName, "provider_group.0.value", "prod-ad-servers"),
					resource.TestCheckResourceAttr(resourceName, "provider_region.0.value", "us-west-2"),
					resource.TestCheckResourceAttr(resourceName, "alternative_ips.0.value", "10.0.1.5\n10.0.1.6\n10.0.1.7"),
					resource.TestCheckResourceAttr(resourceName, "use_ssl.0.value.0", "true"),
				),
			},
		},
	})
}

// TestAccResourcePamDirectory_update tests updating PAM Directory fields.
// NOTE: This test only updates pam_hostname and pam_settings which use SetStandardFieldValue().
// Fields like distinguished_name and use_ssl use ApplyFieldChange() which doesn't sync
// RecordDict to RawJson due to an SDK limitation. See resourcePamDirectoryUpdate for details.
func TestAccResourcePamDirectory_update(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_test_pam_directory_update"
	secretTitle2 := "tf_acc_test_pam_directory_update_2"
	if secretFolderUid == "" {
		t.Skip("Skipping test - TF_ACC not set or test folder not configured")
	}

	configInit := fmt.Sprintf(`
		resource "secretsmanager_pam_directory" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "%v"
			pam_hostname {
				value {
					hostname = "ad.corp.example.com"
					port = "636"
				}
			}
			pam_settings = jsonencode([{
				connection = [{
					protocol = "ldaps"
					port = "636"
				}]
			}])
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle)

	configUpdate := fmt.Sprintf(`
		resource "secretsmanager_pam_directory" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "%v"
			pam_hostname {
				value {
					hostname = "ldap.dev.example.com"
					port = "389"
				}
			}
			pam_settings = jsonencode([{
				connection = [{
					protocol = "ldap"
					port = "389"
					recordingIncludeKeys = false
				}]
			}])
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle2)

	resourceName := fmt.Sprintf("secretsmanager_pam_directory.%v", secretTitle)

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
						if s.Attributes["pam_hostname.0.value.0.hostname"] != "ad.corp.example.com" {
							return fmt.Errorf("expected hostname = 'ad.corp.example.com'")
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
						if actualHostname != "ldap.dev.example.com" {
							return fmt.Errorf("expected hostname = 'ldap.dev.example.com', got '%s'", actualHostname)
						}
						actualPort := s.Attributes["pam_hostname.0.value.0.port"]
						if actualPort != "389" {
							return fmt.Errorf("expected port = '389', got '%s'", actualPort)
						}
						return nil
					}),
					checkSecretExistsRemotely(secretUid),
				),
			},
		},
	})
}

func TestAccResourcePamDirectory_deleteDetection(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_test_pam_directory_delete"
	if secretFolderUid == "" {
		t.Skip("Skipping test - TF_ACC not set or test folder not configured")
	}

	config := fmt.Sprintf(`
		resource "secretsmanager_pam_directory" "%v" {
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

func TestAccResourcePamDirectory_import(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_test_pam_directory_import"
	if secretFolderUid == "" {
		t.Skip("Skipping test - TF_ACC not set or test folder not configured")
	}

	config := fmt.Sprintf(`
		resource "secretsmanager_pam_directory" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "%v"
			pam_hostname {
				value {
					hostname = "ad.example.com"
					port = "636"
				}
			}
			directory_type = "Active Directory"
			login {
				value = "CN=Admin,DC=example,DC=com"
			}
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle)

	resourceName := fmt.Sprintf("secretsmanager_pam_directory.%v", secretTitle)

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
