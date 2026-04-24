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

		# check blocks can reference ephemeral values; a failed assertion surfaces
		# as an apply-time warning and causes terraform test (.tftest.hcl) to fail.
		check "ephemeral_custom_fields" {
			assert {
				condition     = length(ephemeral.secretsmanager_login.custom_eph.custom) == 1
				error_message = "expected 1 custom field, got ${length(ephemeral.secretsmanager_login.custom_eph.custom)}"
			}
			assert {
				condition     = ephemeral.secretsmanager_login.custom_eph.custom[0].type == "text"
				error_message = "expected custom[0].type == \"text\""
			}
			assert {
				condition     = ephemeral.secretsmanager_login.custom_eph.custom[0].label == "Environment"
				error_message = "expected custom[0].label == \"Environment\""
			}
			assert {
				condition     = ephemeral.secretsmanager_login.custom_eph.custom[0].value == "production"
				error_message = "expected custom[0].value == \"production\""
			}
			assert {
				condition     = ephemeral.secretsmanager_login.custom_eph.custom[0].required == false
				error_message = "expected custom[0].required == false"
			}
			assert {
				condition     = ephemeral.secretsmanager_login.custom_eph.custom[0].privacy_screen == false
				error_message = "expected custom[0].privacy_screen == false"
			}
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

// TestAccEphemeralLogin_customFieldComplexType verifies that a complex custom field
// (phone, which getFieldItemsData serialises to JSON) is returned as a JSON string
// through the ephemeral resource — matching the data source contract.
func TestAccEphemeralLogin_customFieldComplexType(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_eph_custom_phone"

	configResource := fmt.Sprintf(`
		resource "secretsmanager_login" "custom_eph_phone" {
			folder_uid = "%v"
			uid        = "%v"
			title      = "%v"

			custom {
				type  = "phone"
				label = "WorkPhone"
				value = jsonencode({
					region = "US"
					number = "555-867-5309"
					ext    = "42"
					type   = "Work"
				})
			}
		}
	`, secretFolderUid, secretUid, secretTitle)

	configWithEphemeral := configResource + fmt.Sprintf(`
		ephemeral "secretsmanager_login" "custom_eph_phone" {
			path = "%v"
		}

		check "ephemeral_custom_phone" {
			assert {
				condition     = ephemeral.secretsmanager_login.custom_eph_phone.custom[0].type == "phone"
				error_message = "expected custom[0].type == \"phone\""
			}
			assert {
				condition     = ephemeral.secretsmanager_login.custom_eph_phone.custom[0].label == "WorkPhone"
				error_message = "expected custom[0].label == \"WorkPhone\""
			}
			# Complex types are returned as a JSON string — verify it is non-empty
			# and contains the expected phone number.
			assert {
				condition     = can(jsondecode(ephemeral.secretsmanager_login.custom_eph_phone.custom[0].value))
				error_message = "expected custom[0].value to be valid JSON for a complex phone field"
			}
			assert {
				condition     = jsondecode(ephemeral.secretsmanager_login.custom_eph_phone.custom[0].value).number == "555-867-5309"
				error_message = "expected phone number 555-867-5309 in custom[0].value JSON"
			}
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
