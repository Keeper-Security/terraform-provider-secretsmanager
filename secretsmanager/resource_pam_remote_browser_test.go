package secretsmanager

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/keeper-security/secrets-manager-go/core"
)

func TestAccResourcePamRemoteBrowser_create(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_test_pam_remote_browser_create"
	if secretFolderUid == "" {
		t.Skip("Skipping test - TF_ACC not set or test folder not configured")
	}

	config := fmt.Sprintf(`
		resource "secretsmanager_pam_remote_browser" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			notes = "%v"
			rbi_url {
				value = "https://internal-app.example.com"
			}
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle)

	resourceName := fmt.Sprintf("secretsmanager_pam_remote_browser.%v", secretTitle)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkSecretExistsRemotely(secretUid),
					resource.TestCheckResourceAttr(resourceName, "type", "pamRemoteBrowser"),
					resource.TestCheckResourceAttr(resourceName, "title", secretTitle),
					resource.TestCheckResourceAttr(resourceName, "notes", secretTitle),
				),
			},
		},
	})
}

func TestAccResourcePamRemoteBrowser_create_no_uid(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretTitle := "tf_acc_test_pam_rb_no_uid"
	if secretFolderUid == "" {
		t.Skip("Skipping test - TF_ACC not set or test folder not configured")
	}

	config := fmt.Sprintf(`
		resource "secretsmanager_pam_remote_browser" "%v" {
			folder_uid = "%v"
			title = "%v"
			rbi_url {
				value = "https://admin.example.com"
			}
		}
	`, secretTitle, secretFolderUid, secretTitle)

	resourceName := fmt.Sprintf("secretsmanager_pam_remote_browser.%v", secretTitle)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "pamRemoteBrowser"),
					resource.TestCheckResourceAttr(resourceName, "title", secretTitle),
					resource.TestCheckResourceAttrSet(resourceName, "uid"),
				),
			},
		},
	})
}

func TestAccResourcePamRemoteBrowser_import(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_test_pam_rb_import"
	if secretFolderUid == "" {
		t.Skip("Skipping test - TF_ACC not set or test folder not configured")
	}

	config := fmt.Sprintf(`
		resource "secretsmanager_pam_remote_browser" "%v" {
			folder_uid = "%v"
			uid = "%v"
			title = "%v"
			rbi_url {
				value = "https://portal.example.com"
			}
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle)

	resourceName := fmt.Sprintf("secretsmanager_pam_remote_browser.%v", secretTitle)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkSecretExistsRemotely(secretUid),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:        true,
				ImportStateVerify:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					return secretUid, nil
				},
			},
		},
	})
}

func TestAccResourcePamRemoteBrowser_customField(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_test_pam_remote_browser_custom"
	if secretFolderUid == "" {
		t.Skip("Skipping test - TF_ACC not set or test folder not configured")
	}

	config := fmt.Sprintf(`
		resource "secretsmanager_pam_remote_browser" "custom" {
			folder_uid = "%v"
			uid        = "%v"
			title      = "%v"

			custom {
				type  = "text"
				label = "Owner"
				value = "infra-team"
			}
		}
	`, secretFolderUid, secretUid, secretTitle)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkSecretExistsRemotely(secretUid),
					resource.TestCheckResourceAttr("secretsmanager_pam_remote_browser.custom", "custom.#", "1"),
					resource.TestCheckResourceAttr("secretsmanager_pam_remote_browser.custom", "custom.0.label", "Owner"),
					resource.TestCheckResourceAttr("secretsmanager_pam_remote_browser.custom", "custom.0.value", "infra-team"),
				),
			},
		},
	})
}
