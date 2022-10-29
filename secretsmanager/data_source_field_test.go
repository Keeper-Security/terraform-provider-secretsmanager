package secretsmanager

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceField(t *testing.T) {
	secretUid, secretTitle := testAcc.getRecordInfo("field")
	if secretUid == "" || secretTitle == "" {
		t.Fatal("Failed to access test data - missing secret UID and/or Title")
	}

	config := fmt.Sprintf(`
		data "secretsmanager_field" "%v" {
			path = "%v"
			title = "%v"
		}
		`, secretTitle, secretUid, secretTitle)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						fmt.Sprintf("data.secretsmanager_field.%v", secretTitle),
						"title",
						secretTitle,
					),
					resource.TestCheckResourceAttr(
						fmt.Sprintf("data.secretsmanager_field.%v", secretTitle),
						"value",
						secretTitle,
					),
				),
			},
		},
	})
}
