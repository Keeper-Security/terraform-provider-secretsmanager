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
	_, secretTitle := testAcc.getRecordInfo(secretType)
	if secretTitle == "" {
		t.Fatal("Failed to access test data - missing secret Title")
	}

	config := fmt.Sprintf(`
		data "secretsmanager_metadata" "%v" {
			path  = "*"
			title = "%v"
		}
	`, secretTitle, secretTitle)

	resourceName := fmt.Sprintf("data.secretsmanager_metadata.%v", secretTitle)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", secretType),
					resource.TestCheckResourceAttr(resourceName, "title", secretTitle),
					// title lookup must resolve to a real record UID, not the "*" sentinel
					checkSecretResourceState(resourceName, func(s *terraform.InstanceState) error {
						if uid := s.Attributes["uid"]; uid == "" || uid == "*" {
							return fmt.Errorf("expected a resolved record UID, got %q", uid)
						}
						return nil
					}),
					// revision must be a positive integer
					checkSecretResourceState(resourceName, func(s *terraform.InstanceState) error {
						rev, err := strconv.Atoi(s.Attributes["revision"])
						if err != nil || rev <= 0 {
							return fmt.Errorf("expected revision to be a positive integer, got %q", s.Attributes["revision"])
						}
						return nil
					}),
					// metadata must not expose secret field values
					resource.TestCheckNoResourceAttr(resourceName, "password"),
					resource.TestCheckNoResourceAttr(resourceName, "login"),
				),
			},
		},
	})
}
