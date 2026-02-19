package secretsmanager

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/keeper-security/secrets-manager-go/core"
)

func TestAccDataSourcePamUser(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_test_ds_pam_user"
	if secretFolderUid == "" {
		t.Skip("Skipping test - TF_ACC not set or test folder not configured")
	}

	config := fmt.Sprintf(`
		resource "secretsmanager_pam_user" "%v" {
			folder_uid = "%v"
			uid        = "%v"
			title      = "%v"
			notes      = "%v"
			login {
				value = "testuser"
			}
			rotation_scripts {
				value {
					command = "echo hello"
				}
			}
			distinguished_name {
				label = "Distinguished Name"
				value = "CN=testuser,OU=Users,DC=example,DC=com"
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

		data "secretsmanager_pam_user" "%v" {
			path = secretsmanager_pam_user.%v.uid
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle,
		secretTitle, secretTitle)

	dataName := fmt.Sprintf("data.secretsmanager_pam_user.%v", secretTitle)
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataName, "type", "pamUser"),
					resource.TestCheckResourceAttr(dataName, "title", secretTitle),
					resource.TestCheckResourceAttr(dataName, "notes", secretTitle),
					resource.TestCheckResourceAttr(dataName, "login.0.value", "testuser"),
					resource.TestCheckResourceAttr(dataName, "rotation_scripts.0.value.0.command", "echo hello"),
					resource.TestCheckResourceAttr(dataName, "distinguished_name.0.value", "CN=testuser,OU=Users,DC=example,DC=com"),
					resource.TestCheckResourceAttr(dataName, "private_pem_key.0.value", "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA..."),
					resource.TestCheckResourceAttr(dataName, "connect_database.0.value", "production_db"),
					resource.TestCheckResourceAttr(dataName, "managed.0.value", "true"),
				),
			},
		},
	})
}
