package secretsmanager

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/keeper-security/secrets-manager-go/core"
)

func TestAccResourcePamUser_create(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_test_pam_user_create"
	if secretFolderUid == "" {
		t.Skip("Skipping test - TF_ACC not set or test folder not configured")
	}

	config := fmt.Sprintf(`
		resource "secretsmanager_pam_user" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "%v"
			login {
				value = "dbadmin"
			}
			password {
				enforce_generation = true
				generate = "yes"
				complexity {
					length = 24
					caps = 5
					lowercase = 5
					digits = 5
					special = 5
				}
			}
			distinguished_name {
				label = "Distinguished Name"
				value = "CN=dbadmin,OU=Users,DC=example,DC=com"
			}
			private_pem_key {
				value = "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA..."
			}
			connect_database {
				value = "production_db"
			}
			managed {
				label = "Managed"
				value = true
			}
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle)

	resourceName := fmt.Sprintf("secretsmanager_pam_user.%v", secretTitle)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkSecretExistsRemotely(secretUid),
					resource.TestCheckResourceAttr(resourceName, "type", "pamUser"),
					resource.TestCheckResourceAttr(resourceName, "title", secretTitle),
					resource.TestCheckResourceAttr(resourceName, "notes", secretTitle),
					resource.TestCheckResourceAttr(resourceName, "login.0.value", "dbadmin"),
					resource.TestCheckResourceAttr(resourceName, "distinguished_name.0.value", "CN=dbadmin,OU=Users,DC=example,DC=com"),
					resource.TestCheckResourceAttr(resourceName, "private_pem_key.0.value", "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA..."),
					resource.TestCheckResourceAttr(resourceName, "connect_database.0.value", "production_db"),
					resource.TestCheckResourceAttr(resourceName, "managed.0.value", "true"),
				),
			},
		},
	})
}

// TestAccResourcePamUser_update tests updating PAM User fields.
// NOTE: This test only updates title and notes which don't use ApplyFieldChange().
// Fields like login, distinguished_name, connect_database, and managed use ApplyFieldChange()
// which doesn't sync RecordDict to RawJson due to an SDK limitation. See resourcePamUserUpdate:337-366.
func TestAccResourcePamUser_update(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_test_pam_user_update"
	secretTitle2 := "tf_acc_test_pam_user_update_2"
	if secretFolderUid == "" {
		t.Skip("Skipping test - TF_ACC not set or test folder not configured")
	}

	configInit := fmt.Sprintf(`
		resource "secretsmanager_pam_user" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "%v"
			login {
				value = "dbadmin"
			}
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle)

	configUpdate := fmt.Sprintf(`
		resource "secretsmanager_pam_user" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "%v"
			login {
				value = "dbadmin"
			}
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle2)

	resourceName := fmt.Sprintf("secretsmanager_pam_user.%v", secretTitle)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
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

func TestAccResourcePamUser_deleteDetection(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_test_pam_user_delete"
	if secretFolderUid == "" {
		t.Skip("Skipping test - TF_ACC not set or test folder not configured")
	}

	config := fmt.Sprintf(`
		resource "secretsmanager_pam_user" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				PreConfig: func() {
					// Delete secret outside of Terraform workspace
					client := *testAccClient()
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

func TestAccResourcePamUser_import(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_test_pam_user_import"
	if secretFolderUid == "" {
		t.Skip("Skipping test - TF_ACC not set or test folder not configured")
	}

	config := fmt.Sprintf(`
		resource "secretsmanager_pam_user" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "%v"
			login {
				value = "dbadmin"
			}
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle)

	resourceName := fmt.Sprintf("secretsmanager_pam_user.%v", secretTitle)

	resource.Test(t, resource.TestCase{
		PreCheck:  testAccPreCheck(t),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
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

func TestAccResourcePamUser_generatePrivatePemKey(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_test_pam_user_generate_pem_key"
	if secretFolderUid == "" {
		t.Skip("Skipping test - TF_ACC not set or test folder not configured")
	}

	config := fmt.Sprintf(`
		resource "secretsmanager_pam_user" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			login {
				value = "svcaccount"
			}
			private_pem_key {
				generate = "yes"
				key_type = "ssh-ed25519"
			}
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle)

	resourceName := fmt.Sprintf("secretsmanager_pam_user.%v", secretTitle)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkSecretExistsRemotely(secretUid),
					resource.TestCheckResourceAttr(resourceName, "type", "pamUser"),
					checkSecretResourceState(resourceName, func(s *terraform.InstanceState) error {
						privKey := s.Attributes["private_pem_key.0.value"]
						if privKey == "" {
							return fmt.Errorf("expected non-empty private_pem_key value")
						}
						if len(privKey) < 20 {
							return fmt.Errorf("private_pem_key value too short: %s", privKey)
						}
						pubKey := s.Attributes["private_pem_key.0.public_key"]
						if pubKey == "" {
							return fmt.Errorf("expected non-empty private_pem_key public_key")
						}
						return nil
					}),
				),
			},
		},
	})
}

func TestAccResourcePamUser_generatePrivatePemKeyWithPassphrase(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_test_pam_user_generate_pem_passphrase"
	if secretFolderUid == "" {
		t.Skip("Skipping test - TF_ACC not set or test folder not configured")
	}

	config := fmt.Sprintf(`
		resource "secretsmanager_pam_user" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			login {
				value = "svcaccount"
			}
			private_key_passphrase {
				generate = "yes"
				complexity {
					length = 20
					caps = 5
					lowercase = 5
					digits = 5
					special = 5
				}
			}
			private_pem_key {
				generate = "yes"
				key_type = "ssh-rsa"
				key_bits = 4096
			}
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle)

	resourceName := fmt.Sprintf("secretsmanager_pam_user.%v", secretTitle)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkSecretExistsRemotely(secretUid),
					resource.TestCheckResourceAttr(resourceName, "type", "pamUser"),
					checkSecretResourceState(resourceName, func(s *terraform.InstanceState) error {
						privKey := s.Attributes["private_pem_key.0.value"]
						if privKey == "" {
							return fmt.Errorf("expected non-empty private_pem_key value")
						}
						pubKey := s.Attributes["private_pem_key.0.public_key"]
						if pubKey == "" {
							return fmt.Errorf("expected non-empty private_pem_key public_key")
						}
						passphrase := s.Attributes["private_key_passphrase.0.value"]
						if passphrase == "" {
							return fmt.Errorf("expected non-empty private_key_passphrase value")
						}
						return nil
					}),
				),
			},
		},
	})
}

// TestAccResourcePamUser_customField verifies that:
//  1. A user-defined custom field can be set alongside private_key_passphrase
//  2. The passphrase does NOT appear in the user-visible custom list (filtered from state)
//  3. Updating the custom field does not wipe the passphrase from the vault
func TestAccResourcePamUser_customField(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_test_pam_user_custom"
	if secretFolderUid == "" {
		t.Skip("Skipping test - TF_ACC not set or test folder not configured")
	}

	configCreate := fmt.Sprintf(`
		resource "secretsmanager_pam_user" "custom" {
			folder_uid = "%v"
			uid        = "%v"
			title      = "%v"

			login { value = "svcaccount" }

			private_key_passphrase {
				value = "s3cr3tpassphrase"
			}

			custom {
				type  = "text"
				label = "Team"
				value = "backend"
			}
		}
	`, secretFolderUid, secretUid, secretTitle)

	configUpdate := fmt.Sprintf(`
		resource "secretsmanager_pam_user" "custom" {
			folder_uid = "%v"
			uid        = "%v"
			title      = "%v"

			login { value = "svcaccount" }

			private_key_passphrase {
				value = "s3cr3tpassphrase"
			}

			custom {
				type  = "text"
				label = "Team"
				value = "platform"
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
					resource.TestCheckResourceAttr("secretsmanager_pam_user.custom", "custom.#", "1"),
					resource.TestCheckResourceAttr("secretsmanager_pam_user.custom", "custom.0.label", "Team"),
					resource.TestCheckResourceAttr("secretsmanager_pam_user.custom", "custom.0.value", "backend"),
					checkSecretResourceState("secretsmanager_pam_user.custom", func(s *terraform.InstanceState) error {
						for k, v := range s.Attributes {
							if v == "Private Key Passphrase" && k != "private_key_passphrase.0.type" {
								return fmt.Errorf("passphrase label leaked into custom list at key %s", k)
							}
						}
						return nil
					}),
				),
			},
			{
				Config: configUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("secretsmanager_pam_user.custom", "custom.0.value", "platform"),
					checkSecretResourceState("secretsmanager_pam_user.custom", func(s *terraform.InstanceState) error {
						passphrase := s.Attributes["private_key_passphrase.0.value"]
						if passphrase == "" {
							return fmt.Errorf("private_key_passphrase was wiped by custom field update")
						}
						return nil
					}),
				),
			},
		},
	})
}
