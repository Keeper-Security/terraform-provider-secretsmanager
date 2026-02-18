package secretsmanager

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/keeper-security/secrets-manager-go/core"
)

func TestAccDataSourcePamMachine(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_test_ds_pam_machine"
	if secretFolderUid == "" {
		t.Skip("Skipping test - TF_ACC not set or test folder not configured")
	}

	config := fmt.Sprintf(`
		resource "secretsmanager_pam_machine" "%v" {
			folder_uid = "%v"
			uid        = "%v"
			title      = "%v"
			notes      = "%v"
			pam_hostname {
				value {
					hostname = "192.168.1.100"
					port     = "22"
				}
			}
			operating_system {
				label = "Operating System"
				value = "Linux"
			}
		}

		data "secretsmanager_pam_machine" "%v" {
			path = secretsmanager_pam_machine.%v.uid
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle,
		secretTitle, secretTitle)

	dataName := fmt.Sprintf("data.secretsmanager_pam_machine.%v", secretTitle)
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataName, "type", "pamMachine"),
					resource.TestCheckResourceAttr(dataName, "title", secretTitle),
					resource.TestCheckResourceAttr(dataName, "notes", secretTitle),
					resource.TestCheckResourceAttr(dataName, "pam_hostname.0.value.0.hostname", "192.168.1.100"),
					resource.TestCheckResourceAttr(dataName, "pam_hostname.0.value.0.port", "22"),
					resource.TestCheckResourceAttr(dataName, "operating_system.0.value", "Linux"),
				),
			},
		},
	})
}
