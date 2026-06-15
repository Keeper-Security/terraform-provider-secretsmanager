package secretsmanager

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceMetadata(t *testing.T) {
	secretType := "login"
	secretUid, secretTitle := testAcc.getRecordInfo(secretType)
	if secretUid == "" || secretTitle == "" {
		t.Fatal("Failed to access test data - missing secret UID and/or Title")
	}

	config := fmt.Sprintf(`
		data "secretsmanager_metadata" "%v" {
			path = "%v"
		}
	`, secretTitle, secretUid)

	resourceName := fmt.Sprintf("data.secretsmanager_metadata.%v", secretTitle)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "uid", secretUid),
					resource.TestCheckResourceAttr(resourceName, "type", secretType),
					resource.TestCheckResourceAttr(resourceName, "title", secretTitle),
					// revision must be a positive integer
					checkSecretResourceState(resourceName, func(s *terraform.InstanceState) error {
						rev, err := strconv.Atoi(s.Attributes["revision"])
						if err != nil || rev <= 0 {
							return fmt.Errorf("expected revision to be a positive integer, got %q", s.Attributes["revision"])
						}
						return nil
					}),
					// no secret field values present in state
					resource.TestCheckNoResourceAttr(resourceName, "password"),
					resource.TestCheckNoResourceAttr(resourceName, "login"),
				),
			},
		},
	})
}
