package secretsmanager

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/keeper-security/secrets-manager-go/core"
)

func TestAccDataSourcePamRemoteBrowser(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_test_pam_rb_datasource"
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
				value = "https://test-rbi.example.com"
			}
		}

		data "secretsmanager_pam_remote_browser" "%v" {
			path = secretsmanager_pam_remote_browser.%v.uid
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle, secretTitle, secretTitle)

	dataSourceName := fmt.Sprintf("data.secretsmanager_pam_remote_browser.%v", secretTitle)
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "type", "pamRemoteBrowser"),
					resource.TestCheckResourceAttr(dataSourceName, "title", secretTitle),
					resource.TestCheckResourceAttr(dataSourceName, "notes", secretTitle),
				),
			},
		},
	})
}
