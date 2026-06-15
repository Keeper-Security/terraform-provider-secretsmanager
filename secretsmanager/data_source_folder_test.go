package secretsmanager

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceFolder_parentUid(t *testing.T) {
	testFolderUid := testAcc.getTestFolder()
	if testFolderUid == "" {
		t.Skip("Skipping test - TF_ACC not set or test folder not configured")
	}
	folderName := "tf_acc_test_datasource_folder_parent_uid"
	resourceName := fmt.Sprintf("secretsmanager_folder.%v", folderName)
	dataName := fmt.Sprintf("data.secretsmanager_folder.%v", folderName)

	config := fmt.Sprintf(`
		resource "secretsmanager_folder" "%v" {
			parent_uid   = "%v"
			name         = "%v"
			force_delete = true
		}
		data "secretsmanager_folder" "%v" {
			depends_on = [secretsmanager_folder.%v]
			name       = "%v"
			parent_uid = "%v"
		}
	`, folderName, testFolderUid, folderName,
		folderName, folderName, folderName, testFolderUid)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataName, "name", folderName),
					resource.TestCheckResourceAttr(dataName, "parent_uid", testFolderUid),
					// uid on data source matches the created resource uid
					checkSecretResourceState(resourceName, func(s *terraform.InstanceState) error {
						return nil // existence is sufficient; uid match checked via TestCheckResourceAttrPair below
					}),
					resource.TestCheckResourceAttrPair(dataName, "uid", resourceName, "uid"),
				),
			},
		},
	})
}

func TestAccDataSourceFolder(t *testing.T) {
	// Get test folder UID - this will be populated during test setup
	testFolderUid := testAcc.getTestFolder()
	if testFolderUid == "" {
		t.Skip("Skipping test - TF_ACC not set or test folder not configured")
	}
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
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
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
