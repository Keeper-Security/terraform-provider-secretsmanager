package secretsmanager

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/keeper-security/secrets-manager-go/core"
)

// TestAccDataSourceLogin_customFields verifies that custom fields written to a
// resource are readable back through the corresponding data source.
func TestAccDataSourceLogin_customFields(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_ds_custom_fields"

	config := fmt.Sprintf(`
		resource "secretsmanager_login" "custom_ds" {
			folder_uid = "%v"
			uid        = "%v"
			title      = "%v"

			custom {
				type  = "text"
				label = "Environment"
				value = "production"
			}
		}

		data "secretsmanager_login" "custom_ds" {
			path       = secretsmanager_login.custom_ds.uid
			depends_on = [secretsmanager_login.custom_ds]
		}
	`, secretFolderUid, secretUid, secretTitle)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.secretsmanager_login.custom_ds", "custom.#", "1"),
					resource.TestCheckResourceAttr("data.secretsmanager_login.custom_ds", "custom.0.type", "text"),
					resource.TestCheckResourceAttr("data.secretsmanager_login.custom_ds", "custom.0.label", "Environment"),
					resource.TestCheckResourceAttr("data.secretsmanager_login.custom_ds", "custom.0.value", "production"),
				),
			},
		},
	})
}
