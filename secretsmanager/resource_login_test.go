package secretsmanager

import (
	"fmt"
	"strings"
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
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 testAccPreCheck(t),
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
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 testAccPreCheck(t),
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
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 testAccPreCheck(t),
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

func TestAccResourceLogin_generateNoSpecial(t *testing.T) {
	secretType := "login"
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	_, secretTitle := testAcc.getRecordInfo(secretType)
	if secretUid == "" || secretTitle == "" {
		t.Fatal("Failed to access test data - missing secret UID and/or Title")
	}
	secretTitle += "_resource_no_special"

	config := fmt.Sprintf(`
		resource "secretsmanager_login" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			password {
				generate = "yes"
				complexity {
					length  = 20
					special = 0
				}
			}
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle)

	resourceName := fmt.Sprintf("secretsmanager_login.%v", secretTitle)
	const defaultSpecialChars = "\"!@#$%()+;<>=?[]{}^.,"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkSecretResourceState(resourceName, func(s *terraform.InstanceState) error {
						pwd := s.Attributes["password.0.value"]
						if len(pwd) != 20 {
							return fmt.Errorf("expected a 20-char password, got %d chars", len(pwd))
						}
						if strings.ContainsAny(pwd, defaultSpecialChars) {
							return fmt.Errorf("expected no special chars (special=0) but password contains one: %q", pwd)
						}
						return nil
					}),
					checkSecretExistsRemotely(secretUid),
				),
			},
		},
	})
}

func TestAccResourceLogin_generateSpecialSet(t *testing.T) {
	secretType := "login"
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	_, secretTitle := testAcc.getRecordInfo(secretType)
	if secretUid == "" || secretTitle == "" {
		t.Fatal("Failed to access test data - missing secret UID and/or Title")
	}
	secretTitle += "_resource_special_set"

	config := fmt.Sprintf(`
		resource "secretsmanager_login" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			password {
				generate = "yes"
				complexity {
					length      = 20
					special     = 2
					special_set = "!@"
				}
			}
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle)

	resourceName := fmt.Sprintf("secretsmanager_login.%v", secretTitle)
	const disallowedSpecialChars = "\"#$%()+;<>=?[]{}^.,"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkSecretResourceState(resourceName, func(s *terraform.InstanceState) error {
						pwd := s.Attributes["password.0.value"]
						if len(pwd) != 20 {
							return fmt.Errorf("expected a 20-char password, got %d chars", len(pwd))
						}
						if strings.ContainsAny(pwd, disallowedSpecialChars) {
							return fmt.Errorf("expected only !@ as special chars but password contains a disallowed special: %q", pwd)
						}
						return nil
					}),
					checkSecretExistsRemotely(secretUid),
				),
			},
		},
	})
}

func TestAccResourceLogin_updateNoSpecial(t *testing.T) {
	secretType := "login"
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	_, secretTitle := testAcc.getRecordInfo(secretType)
	if secretUid == "" || secretTitle == "" {
		t.Fatal("Failed to access test data - missing secret UID and/or Title")
	}
	secretTitle += "_resource_update_no_special"

	// Step 1: create with generate="yes" to establish the resource
	configCreate := fmt.Sprintf(`
		resource "secretsmanager_login" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			password {
				generate = "yes"
				complexity {
					length  = 20
					special = 0
				}
			}
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle)

	// Step 2: change generate "yes"→"true" to trigger Update-path regeneration
	configUpdate := fmt.Sprintf(`
		resource "secretsmanager_login" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			password {
				generate = "true"
				complexity {
					length  = 20
					special = 0
				}
			}
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle)

	resourceName := fmt.Sprintf("secretsmanager_login.%v", secretTitle)
	const defaultSpecialChars = "\"!@#$%()+;<>=?[]{}^.,"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 testAccPreCheck(t),
		Steps: []resource.TestStep{
			{Config: configCreate},
			{
				Config: configUpdate,
				Check: resource.ComposeTestCheckFunc(
					checkSecretResourceState(resourceName, func(s *terraform.InstanceState) error {
						pwd := s.Attributes["password.0.value"]
						if len(pwd) != 20 {
							return fmt.Errorf("expected a 20-char password after update, got %d chars", len(pwd))
						}
						if strings.ContainsAny(pwd, defaultSpecialChars) {
							return fmt.Errorf("expected no special chars (special=0) on update but password contains one: %q", pwd)
						}
						return nil
					}),
				),
			},
		},
	})
}

func TestAccResourceLogin_updateSpecialSet(t *testing.T) {
	secretType := "login"
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	_, secretTitle := testAcc.getRecordInfo(secretType)
	if secretUid == "" || secretTitle == "" {
		t.Fatal("Failed to access test data - missing secret UID and/or Title")
	}
	secretTitle += "_resource_update_special_set"

	configCreate := fmt.Sprintf(`
		resource "secretsmanager_login" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			password {
				generate = "yes"
				complexity {
					length      = 20
					special     = 2
					special_set = "!@"
				}
			}
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle)

	configUpdate := fmt.Sprintf(`
		resource "secretsmanager_login" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			password {
				generate = "true"
				complexity {
					length      = 20
					special     = 2
					special_set = "!@"
				}
			}
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle)

	resourceName := fmt.Sprintf("secretsmanager_login.%v", secretTitle)
	const disallowedSpecialChars = "\"#$%()+;<>=?[]{}^.,"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 testAccPreCheck(t),
		Steps: []resource.TestStep{
			{Config: configCreate},
			{
				Config: configUpdate,
				Check: resource.ComposeTestCheckFunc(
					checkSecretResourceState(resourceName, func(s *terraform.InstanceState) error {
						pwd := s.Attributes["password.0.value"]
						if len(pwd) != 20 {
							return fmt.Errorf("expected a 20-char password after update, got %d chars", len(pwd))
						}
						if strings.ContainsAny(pwd, disallowedSpecialChars) {
							return fmt.Errorf("expected only !@ as special chars on update but password contains a disallowed special: %q", pwd)
						}
						return nil
					}),
				),
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
		PreCheck:                 testAccPreCheck(t),
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

// TestAccResourceLogin_adoptedValueIsNotRegenerated covers KSM-1305: a record that
// already holds a value must not have that value replaced merely because the
// configuration has started declaring generate.
//
// The condition under test is the one terraform import creates. "generate" is never
// persisted to the vault, so after an import the flag reads back empty while the password
// value is present, and any configuration declaring generate therefore looks like a
// change. Reaching that state through import inside the acceptance framework is awkward,
// because an ImportState step verifies against a throwaway state rather than continuing
// with it, so this test reproduces the same code path by starting from a literal value
// and then adding the flag.
//
// Step 1 establishes a literal value with no generate flag.
// Step 2 declares generate for the first time, so the existing value must be preserved.
// Step 3 changes the flag again, which is an unambiguous rotation request, so the value
// must change.
func TestAccResourceLogin_adoptedValueIsNotRegenerated(t *testing.T) {
	secretType := "login"
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	_, secretTitle := testAcc.getRecordInfo(secretType)
	if secretUid == "" || secretTitle == "" {
		t.Fatal("Failed to access test data - missing secret UID and/or Title")
	}
	secretTitle += "_resource_adopted_value"

	const literal = "AdoptedLiteralValue123"

	configLiteral := fmt.Sprintf(`
		resource "secretsmanager_login" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			password {
				value = "%v"
			}
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle, literal)

	configGenerateTemplate := `
		resource "secretsmanager_login" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			password {
				generate = "%v"
				complexity {
					length  = 25
					special = 0
				}
			}
		}
	`
	configAdopt := fmt.Sprintf(configGenerateTemplate, secretTitle, secretFolderUid, secretUid, secretTitle, "yes")
	configRotate := fmt.Sprintf(configGenerateTemplate, secretTitle, secretFolderUid, secretUid, secretTitle, "true")

	resourceName := fmt.Sprintf("secretsmanager_login.%v", secretTitle)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: configLiteral,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "password.0.value", literal),
				),
			},
			{
				// generate is declared for the first time on a record that already holds
				// a value. This is adoption, not a rotation request, so the stored secret
				// must survive untouched.
				Config: configAdopt,
				Check: resource.ComposeTestCheckFunc(
					checkSecretResourceState(resourceName, func(s *terraform.InstanceState) error {
						if pwd := s.Attributes["password.0.value"]; pwd != literal {
							return fmt.Errorf("adopting a record must not regenerate its value: expected %q, got %q", literal, pwd)
						}
						return nil
					}),
				),
			},
			{
				// The flag changes from a recorded value, so this is an explicit rotation.
				Config: configRotate,
				Check: resource.ComposeTestCheckFunc(
					checkSecretResourceState(resourceName, func(s *terraform.InstanceState) error {
						pwd := s.Attributes["password.0.value"]
						if pwd == literal {
							return fmt.Errorf("changing the generate flag on an adopted record must rotate the value, but it is unchanged")
						}
						if len(pwd) != 25 {
							return fmt.Errorf("expected a 25-char generated password after rotation, got %d chars", len(pwd))
						}
						if strings.ContainsAny(pwd, "\"!@#$%()+;<>=?[]{}^.,") {
							return fmt.Errorf("expected no special characters with special = 0, got %q", pwd)
						}
						return nil
					}),
				),
			},
		},
	})
}
