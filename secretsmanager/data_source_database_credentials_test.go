package secretsmanager

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDatabaseCredentials(t *testing.T) {
	secretType := "databaseCredentials"
	secretUid, secretTitle := testAcc.getRecordInfo(secretType)
	if secretUid == "" || secretTitle == "" {
		t.Fatal("Failed to access test data - missing secret UID and/or Title")
	}

	config := fmt.Sprintf(`
		data "secretsmanager_database_credentials" "%v" {
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
						fmt.Sprintf("data.secretsmanager_database_credentials.%v", secretTitle),
						"type",
						secretType,
					),
					resource.TestCheckResourceAttr(
						fmt.Sprintf("data.secretsmanager_database_credentials.%v", secretTitle),
						"title",
						secretTitle,
					),
					resource.TestCheckResourceAttr(
						fmt.Sprintf("data.secretsmanager_database_credentials.%v", secretTitle),
						"notes",
						secretTitle,
					),
				),
			},
		},
	})
}
