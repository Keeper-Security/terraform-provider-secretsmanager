package secretsmanager

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/keeper-security/secrets-manager-go/core"
)

// TestAccEphemeralLogin_customFields verifies that a login record with custom
// fields can be read through the ephemeral resource without error. Uses two
// steps: step 1 creates the resource; step 2 reads it with the ephemeral
// resource. This avoids a pre-apply refresh failure caused by uid being a
// known value at plan time before the record exists.
func TestAccEphemeralLogin_customFields(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_eph_custom_fields"

	configResource := fmt.Sprintf(`
		resource "secretsmanager_login" "custom_eph" {
			folder_uid = "%v"
			uid        = "%v"
			title      = "%v"

			custom {
				type  = "text"
				label = "Environment"
				value = "production"
			}
		}
	`, secretFolderUid, secretUid, secretTitle)

	configWithEphemeral := configResource + fmt.Sprintf(`
		ephemeral "secretsmanager_login" "custom_eph" {
			path = "%v"
		}
	`, secretUid)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 testAccPreCheck(t),
		Steps: []resource.TestStep{
			{Config: configResource},
			{Config: configWithEphemeral},
		},
	})
}
