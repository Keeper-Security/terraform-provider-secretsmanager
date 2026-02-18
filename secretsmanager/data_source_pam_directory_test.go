package secretsmanager

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/keeper-security/secrets-manager-go/core"
)

func TestAccDataSourcePamDirectory(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_test_ds_pam_directory"
	if secretFolderUid == "" {
		t.Skip("Skipping test - TF_ACC not set or test folder not configured")
	}

	config := fmt.Sprintf(`
		resource "secretsmanager_pam_directory" "%v" {
			folder_uid = "%v"
			uid        = "%v"
			title      = "%v"
			notes      = "%v"
			pam_hostname {
				value {
					hostname = "ad.corp.example.com"
					port     = "636"
				}
			}
			directory_type = "Active Directory"
			distinguished_name {
				label = "Distinguished Name"
				value = "DC=corp,DC=example,DC=com"
			}
		}

		data "secretsmanager_pam_directory" "%v" {
			path = secretsmanager_pam_directory.%v.uid
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle,
		secretTitle, secretTitle)

	dataName := fmt.Sprintf("data.secretsmanager_pam_directory.%v", secretTitle)
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataName, "type", "pamDirectory"),
					resource.TestCheckResourceAttr(dataName, "title", secretTitle),
					resource.TestCheckResourceAttr(dataName, "notes", secretTitle),
					resource.TestCheckResourceAttr(dataName, "pam_hostname.0.value.0.hostname", "ad.corp.example.com"),
					resource.TestCheckResourceAttr(dataName, "pam_hostname.0.value.0.port", "636"),
					resource.TestCheckResourceAttr(dataName, "directory_type", "Active Directory"),
					resource.TestCheckResourceAttr(dataName, "distinguished_name.0.value", "DC=corp,DC=example,DC=com"),
				),
			},
		},
	})
}
