package secretsmanager

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceFolder(t *testing.T) {
	accProvider, d := getConfiguredProvider(os.Getenv(envCredential))
	if d.HasError() {
		t.Fail()
	}
	client := *accProvider.client

	folders, err := findFolder("", "", "tf_acc_test_dir", client)
	if err != nil || len(folders) == 0 {
		t.Fail()
	}

	testFolderUid := folders[0].FolderUid
	secretTitle := "tf_acc_test_datasource_folder"
	secretTitleNew := secretTitle + "_new"
	resourceName := fmt.Sprintf("secretsmanager_folder.%v", secretTitleNew)
	config := fmt.Sprintf(`resource "secretsmanager_folder" "%v" {
		parent_uid = "%v"
	 	name = "%v"
		force_delete = true
	}
	data "secretsmanager_folder" "%v" {
		depends_on = [%v]
		name = %v.name
	}
	`, secretTitleNew, testFolderUid, secretTitle, secretTitle, resourceName, resourceName)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkSecretResourceState(resourceName, func(s *terraform.InstanceState) error {
						if s.Attributes["name"] != secretTitle {
							return fmt.Errorf("expected 'name' = '%s'", secretTitle)
						}
						return nil
					}),
					checkFolderExistsRemotely("", secretTitle),
					resource.TestCheckResourceAttr(
						fmt.Sprintf("data.secretsmanager_folder.%v", secretTitle),
						"name",
						secretTitle,
					),
				),
			},
		},
	})
}
