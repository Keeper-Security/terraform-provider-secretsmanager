package secretsmanager

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/keeper-security/secrets-manager-go/core"
)

func TestAccDataSourcePamDatabase(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_test_ds_pam_database"
	if secretFolderUid == "" {
		t.Skip("Skipping test - TF_ACC not set or test folder not configured")
	}

	config := fmt.Sprintf(`
		resource "secretsmanager_pam_database" "%v" {
			folder_uid = "%v"
			uid        = "%v"
			title      = "%v"
			notes      = "%v"
			pam_hostname {
				value {
					hostname = "db.example.com"
					port     = "5432"
				}
			}
			pam_settings = jsonencode([{
				connection = [{
					protocol = "postgresql"
					port = "5432"
				}]
			}])
			database_type = "postgresql"
			use_ssl {
				value = true
			}
			totp {
				value = "otpauth://totp/keeper:database?secret=JBSWY3DPEHPK3PXP&issuer=Keeper"
			}
			rotation_scripts {
				value {
					command = "echo hello"
				}
			}
			database_id {
				value = "db-prod-01"
			}
			provider_group {
				value = "production-servers"
			}
			provider_region {
				value = "us-east-1"
			}
		}

		data "secretsmanager_pam_database" "%v" {
			path = secretsmanager_pam_database.%v.uid
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle,
		secretTitle, secretTitle)

	dataName := fmt.Sprintf("data.secretsmanager_pam_database.%v", secretTitle)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataName, "type", "pamDatabase"),
					resource.TestCheckResourceAttr(dataName, "title", secretTitle),
					resource.TestCheckResourceAttr(dataName, "notes", secretTitle),
					resource.TestCheckResourceAttr(dataName, "folder_uid", secretFolderUid),
					resource.TestCheckResourceAttr(dataName, "pam_hostname.0.value.0.hostname", "db.example.com"),
					resource.TestCheckResourceAttr(dataName, "pam_hostname.0.value.0.port", "5432"),
					resource.TestCheckResourceAttrSet(dataName, "pam_settings"),
					resource.TestCheckResourceAttr(dataName, "database_type", "postgresql"),
					resource.TestCheckResourceAttr(dataName, "use_ssl.0.value", "true"),
					resource.TestCheckResourceAttr(dataName, "rotation_scripts.0.value.0.command", "echo hello"),
					resource.TestCheckResourceAttr(dataName, "database_id.0.value", "db-prod-01"),
					resource.TestCheckResourceAttr(dataName, "provider_group.0.value", "production-servers"),
					resource.TestCheckResourceAttr(dataName, "provider_region.0.value", "us-east-1"),
					resource.TestCheckResourceAttr(dataName, "totp.0.value", "otpauth://totp/keeper:database?secret=JBSWY3DPEHPK3PXP&issuer=Keeper"),
				),
			},
		},
	})
}
