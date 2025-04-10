package secretsmanager

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceFolders(t *testing.T) {
	config := `data "secretsmanager_folders" "folders" { }`

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkListNotEmpty("data.secretsmanager_folders.folders", "folders.#"),
				),
			},
		},
	})
}

func checkListNotEmpty(resourceName, attributeName string) resource.TestCheckFunc {
	// note: attributeName must return count - ex. "folders.#"
	return func(s *terraform.State) error {
		resourceState := s.RootModule().Resources[resourceName]
		if resourceState == nil {
			return fmt.Errorf("resource '%v' not in tf state", resourceName)
		}

		value := resourceState.Primary.Attributes[attributeName]
		count, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("failed to parse count: %v", err)
		}
		if count < 1 {
			return fmt.Errorf("expected non empty list, got empty list for %s.%s", resourceName, attributeName)
		}
		return nil
	}
}
