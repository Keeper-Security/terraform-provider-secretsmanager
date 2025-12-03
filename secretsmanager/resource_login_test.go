package secretsmanager

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/keeper-security/secrets-manager-go/core"
)

func TestAccResourceLogin_create(t *testing.T) {
	secretType := "login"
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	_, secretTitle := testAcc.getRecordInfo(secretType)
	if secretUid == "" || secretTitle == "" {
		t.Fatal("Failed to access test data - missing secret UID and/or Title")
	}
	secretTitle += "_resource_create"

	config := fmt.Sprintf(`
		resource "secretsmanager_login" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "%v"
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
					length = 32
					caps = 8
					lowercase = 8
					digits = 8
					special = 8
				}
				#value = "to_be_generated"
			}
			url {
				label = "MyUrl"
				required = true
				privacy_screen = true
				value = "https://192.168.1.1/"
			}
			totp {
				label = "MyTOTP"
				required = true
				privacy_screen = true
				value = "otpauth://totp/Acme:Buster?secret=6I4PI5EUKS66GPRY5TMLJJP25MAYWAVL&issuer=Acme&algorithm=SHA1&digits=6&period=30"
			}
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle)

	resourceName := fmt.Sprintf("secretsmanager_login.%v", secretTitle)
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

func TestAccResourceLogin_generate(t *testing.T) {
	secretType := "login"
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	_, secretTitle := testAcc.getRecordInfo(secretType)
	if secretUid == "" || secretTitle == "" {
		t.Fatal("Failed to access test data - missing secret UID and/or Title")
	}
	secretTitle += "_resource_generate"

	configInit := fmt.Sprintf(`
		resource "secretsmanager_login" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "%v"
			password {
				generate = "yes"
				complexity {
					length = 32
				}
			}
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle)

	configLengthUpdate := fmt.Sprintf(`
		resource "secretsmanager_login" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "%v"
			password {
				generate = "true"
				complexity {
					length = 16
				}
			}
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle)

	resourceName := fmt.Sprintf("secretsmanager_login.%v", secretTitle)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: configInit,
				Check: resource.ComposeTestCheckFunc(
					checkSecretResourceState(resourceName, func(s *terraform.InstanceState) error {
						if len(s.Attributes["password.0.value"]) != 32 {
							return fmt.Errorf("expected 'value' to contain a 32 char password")
						}
						return nil
					}),
					checkSecretExistsRemotely(secretUid),
				),
			},
			{
				Config: configLengthUpdate,
				Check: resource.ComposeTestCheckFunc(
					checkSecretResourceState(resourceName, func(s *terraform.InstanceState) error {
						if len(s.Attributes["password.0.value"]) != 16 {
							return fmt.Errorf("expected 'value' to contain a 16 char password")
						}
						return nil
					}),
					checkSecretExistsRemotely(secretUid),
				),
			},
		},
	})
}

func TestAccResourceLogin_deleteDetection(t *testing.T) {
	secretType := "login"
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	_, secretTitle := testAcc.getRecordInfo(secretType)
	if secretUid == "" || secretTitle == "" {
		t.Fatal("Failed to access test data - missing secret UID and/or Title")
	}
	secretTitle += "_resource_delete"

	config := fmt.Sprintf(`
		resource "secretsmanager_login" "%v" {
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

func TestAccResourceLogin_import(t *testing.T) {
	secretType := "login"
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	_, secretTitle := testAcc.getRecordInfo(secretType)
	if secretUid == "" || secretTitle == "" {
		t.Fatal("Failed to access test data - missing secret UID and/or Title")
	}
	secretTitle += "_resource_import"

	config := fmt.Sprintf(`
		resource "secretsmanager_login" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "%v"
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle)

	resourceName := fmt.Sprintf("secretsmanager_login.%v", secretTitle)

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
func TestAccResourceLogin_customFields(t *testing.T) {
	secretType := "login"
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	_, secretTitle := testAcc.getRecordInfo(secretType)
	if secretUid == "" || secretTitle == "" {
		t.Fatal("Failed to access test data - missing secret UID and/or Title")
	}
	secretTitle += "_resource_custom_fields"

	configCreate := fmt.Sprintf(`
		resource "secretsmanager_login" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "Test custom fields"
			login {
				value = "testuser"
			}
			password {
				value = "testpass123"
			}
			custom {
				type = "text"
				label = "Department"
				value = "Engineering"
			}
			custom {
				type = "text"
				label = "Employee ID"
				value = "EMP001"
			}
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle)

	configUpdate := fmt.Sprintf(`
		resource "secretsmanager_login" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "Test custom fields updated"
			login {
				value = "testuser"
			}
			password {
				value = "testpass123"
			}
			custom {
				type = "text"
				label = "Department"
				value = "DevOps"
			}
			custom {
				type = "text"
				label = "Employee ID"
				value = "EMP002"
			}
			custom {
				type = "text"
				label = "Cost Center"
				value = "CC123"
			}
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle)

	resourceName := fmt.Sprintf("secretsmanager_login.%v", secretTitle)
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
					resource.TestCheckResourceAttr(resourceName, "custom.0.label", "Department"),
					resource.TestCheckResourceAttr(resourceName, "custom.0.value", "Engineering"),
					resource.TestCheckResourceAttr(resourceName, "custom.1.label", "Employee ID"),
					resource.TestCheckResourceAttr(resourceName, "custom.1.value", "EMP001"),
				),
			},
			{
				Config: configUpdate,
				Check: resource.ComposeTestCheckFunc(
					checkSecretExistsRemotely(secretUid),
					resource.TestCheckResourceAttr(resourceName, "custom.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "custom.0.label", "Department"),
					resource.TestCheckResourceAttr(resourceName, "custom.0.value", "DevOps"),
					resource.TestCheckResourceAttr(resourceName, "custom.1.label", "Employee ID"),
					resource.TestCheckResourceAttr(resourceName, "custom.1.value", "EMP002"),
					resource.TestCheckResourceAttr(resourceName, "custom.2.label", "Cost Center"),
					resource.TestCheckResourceAttr(resourceName, "custom.2.value", "CC123"),
				),
			},
		},
	})
}
