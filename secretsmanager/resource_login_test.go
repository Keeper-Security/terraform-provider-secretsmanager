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

// NOTE: this test holds special = 0 CONSTANT across both steps and varies only the
// generate flag. Regeneration does occur, because the gate in ApplyFieldChange
// compares the raw generate string and "yes" differs from "true", so the assertion
// is live. What it does not cover is special = 0 itself changing; see
// TestAccResourceLogin_changeSpecialSet for that.
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

// NOTE: this test holds special_set CONSTANT across both steps and varies only the
// generate flag. Regeneration does occur, because the gate in ApplyFieldChange
// compares the raw generate string and "yes" differs from "true", so the assertion
// is live. What it does not cover is special_set itself changing; see
// TestAccResourceLogin_changeSpecialSet for that.
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

// TestAccResourceLogin_changeSpecialSet covers the case the two update tests above do
// not: the complexity constraint itself changing between steps.
//
// TestAccResourceLogin_updateNoSpecial and TestAccResourceLogin_updateSpecialSet both
// hold their constraint constant and vary only the generate flag, so they prove that a
// regenerated password still satisfies an unchanged constraint. Neither proves that a
// regenerated password satisfies a NEW constraint, which is the case a user hits when
// they narrow special_set to suit a system that rejects certain characters.
//
// The password is compared across steps rather than only checked for set membership.
// A membership assertion on its own would pass even if nothing regenerated, because the
// original value already satisfied the original set.
func TestAccResourceLogin_changeSpecialSet(t *testing.T) {
	secretType := "login"
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	_, secretTitle := testAcc.getRecordInfo(secretType)
	if secretUid == "" || secretTitle == "" {
		t.Fatal("Failed to access test data - missing secret UID and/or Title")
	}
	secretTitle += "_resource_change_special_set"

	const oldSet = "!@"
	const newSet = "%^"

	configTemplate := `
		resource "secretsmanager_login" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			password {
				generate = "%v"
				complexity {
					length      = 20
					special     = 2
					special_set = "%v"
				}
			}
		}
	`
	configCreate := fmt.Sprintf(configTemplate, secretTitle, secretFolderUid, secretUid, secretTitle, "yes", oldSet)
	configUpdate := fmt.Sprintf(configTemplate, secretTitle, secretFolderUid, secretUid, secretTitle, "true", newSet)

	resourceName := fmt.Sprintf("secretsmanager_login.%v", secretTitle)

	// specialsOutside returns any non-alphanumeric characters in pwd that are not part
	// of allowed. Asserting positively against the configured set is stricter than
	// listing disallowed characters, and it cannot drift as the SDK default set changes.
	specialsOutside := func(pwd, allowed string) string {
		unexpected := []rune{}
		for _, r := range pwd {
			switch {
			case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
				continue
			case strings.ContainsRune(allowed, r):
				continue
			default:
				unexpected = append(unexpected, r)
			}
		}
		return string(unexpected)
	}

	createdPwd := ""

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: configCreate,
				Check: resource.ComposeTestCheckFunc(
					checkSecretResourceState(resourceName, func(s *terraform.InstanceState) error {
						createdPwd = s.Attributes["password.0.value"]
						if len(createdPwd) != 20 {
							return fmt.Errorf("expected a 20-char password on create, got %d chars", len(createdPwd))
						}
						if unexpected := specialsOutside(createdPwd, oldSet); unexpected != "" {
							return fmt.Errorf("create drew specials outside special_set %q: %q", oldSet, unexpected)
						}
						return nil
					}),
				),
			},
			{
				Config: configUpdate,
				Check: resource.ComposeTestCheckFunc(
					checkSecretResourceState(resourceName, func(s *terraform.InstanceState) error {
						pwd := s.Attributes["password.0.value"]
						if len(pwd) != 20 {
							return fmt.Errorf("expected a 20-char password after update, got %d chars", len(pwd))
						}
						if pwd == createdPwd {
							return fmt.Errorf("password was not regenerated after special_set changed from %q to %q", oldSet, newSet)
						}
						if unexpected := specialsOutside(pwd, newSet); unexpected != "" {
							return fmt.Errorf("update drew specials outside the new special_set %q: %q", newSet, unexpected)
						}
						if !strings.ContainsAny(pwd, newSet) {
							return fmt.Errorf("expected at least one character from the new special_set %q, got %q", newSet, pwd)
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
