package secretsmanager

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceFolder_create(t *testing.T) {
	testFolderUid := getTestFolderUid()
	if testFolderUid == "" {
		t.Fail()
	}
	secretTitle := "tf_acc_test_folder_resource_create"
	resourceName := fmt.Sprintf("secretsmanager_folder.%v", secretTitle)
	config := fmt.Sprintf(`resource "secretsmanager_folder" "%v" {
		parent_uid = "%v"
		name = "%v"
		force_delete = true
	}`, secretTitle, testFolderUid, secretTitle)
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkFolderExistsRemotely("", secretTitle),
					resource.TestCheckResourceAttr(resourceName, "name", secretTitle),
				),
			},
		},
	})
}

func TestAccResourceFolder_update(t *testing.T) {
	testFolderUid := getTestFolderUid()
	if testFolderUid == "" {
		t.Fail()
	}
	secretTitle := "tf_acc_test_folder_resource_update"
	resourceName := fmt.Sprintf("secretsmanager_folder.%v", secretTitle)
	configInit := fmt.Sprintf(`resource "secretsmanager_folder" "%v" {
		parent_uid = "%v"
		name = "%v"
		force_delete = true
	}`, secretTitle, testFolderUid, secretTitle)

	secretTitle2 := secretTitle + "2"
	configUpdate := fmt.Sprintf(`resource "secretsmanager_folder" "%v" {
		parent_uid = "%v"
		name = "%v"
		force_delete = true
	}`, secretTitle, testFolderUid, secretTitle2)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: configInit,
				Check: resource.ComposeTestCheckFunc(
					checkSecretResourceState(resourceName, func(s *terraform.InstanceState) error {
						if s.Attributes["name"] != secretTitle {
							return fmt.Errorf("expected 'name' = '%s'", secretTitle)
						}
						return nil
					}),
					checkFolderExistsRemotely("", secretTitle),
				),
			},
			{
				Config: configUpdate,
				Check: resource.ComposeTestCheckFunc(
					checkSecretResourceState(resourceName, func(s *terraform.InstanceState) error {
						if s.Attributes["name"] != secretTitle2 {
							return fmt.Errorf("expected 'name' = '%s'", secretTitle2)
						}
						return nil
					}),
					checkFolderExistsRemotely("", secretTitle2),
				),
			},
		},
	})
}

func TestAccResourceFolder_deleteDetection(t *testing.T) {
	testFolderUid := getTestFolderUid()
	if testFolderUid == "" {
		t.Fail()
	}
	secretTitle := "tf_acc_test_folder_resource_delete"
	config := fmt.Sprintf(`resource "secretsmanager_folder" "%v" {
		parent_uid = "%v"
		name = "%v"
		force_delete = true
	}`, secretTitle, testFolderUid, secretTitle)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				PreConfig: func() {
					// Delete secret outside of Terraform workspace
					client := *testAccProvider.Meta().(providerMeta).client
					folders, err := findFolder("", "", secretTitle, client)
					if err != nil || len(folders) == 0 {
						t.Fail()
					}
					if err := deleteFolder(folders[0].FolderUid, true, client); err != nil {
						t.Fail()
					}
				},
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true, // The externally deleted secret should be planned in for recreation
			},
		},
	})
}

func TestAccResourceFolder_import(t *testing.T) {
	testFolderUid := getTestFolderUid()
	if testFolderUid == "" {
		t.Fail()
	}
	secretTitle := "tf_acc_test_folder_resource_import"
	resourceName := fmt.Sprintf("secretsmanager_folder.%v", secretTitle)
	config := fmt.Sprintf(`resource "secretsmanager_folder" "%v" {
		parent_uid = "%v"
		name = "%v"
	}`, secretTitle, testFolderUid, secretTitle)

	resource.Test(t, resource.TestCase{
		PreCheck:  testAccPreCheck(t),
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// No need to create/delete - Import step takes care of these
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// {
			// 	PreConfig: func() {
			// 		// Delete folder outside of Terraform workspace
			// 		client := *testAccProvider.Meta().(providerMeta).client
			// 		folders, err := findFolder(testFolderUid, "", secretTitle, client)
			// 		if err != nil || len(folders) == 0 {
			// 			t.Fail()
			// 		}
			// 		if err := deleteFolder(folders[0].FolderUid, true, client); err != nil {
			// 			t.Fail()
			// 		}
			// 	},
			// 	Config:   `data "secretsmanager_folder" "test" { name = "tf_acc_test_dir" }`,
			// 	PlanOnly: true,
			// },
		},
	})
}

func getTestFolderUid() string {
	accProvider, d := getConfiguredProvider(os.Getenv(envCredential))
	if d.HasError() {
		return ""
	}

	testFolderName := "tf_acc_test_dir"
	client := *accProvider.client
	folders, err := findFolder("", "", testFolderName, client)
	if err != nil || len(folders) == 0 {
		return ""
	}

	return folders[0].FolderUid
}
