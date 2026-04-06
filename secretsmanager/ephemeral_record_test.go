package secretsmanager

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEphemeralRecord(t *testing.T) {
	secretType := "login"
	secretUid, secretTitle := testAcc.getRecordInfo(secretType)
	if secretUid == "" || secretTitle == "" {
		t.Fatal("Failed to access test data - missing secret UID and/or Title")
	}

	config := fmt.Sprintf(`
		ephemeral "secretsmanager_record" "%v" {
			path = "%v"
			title = "%v"
		}
	`, secretTitle, secretUid, secretTitle)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
			},
		},
	})
}

func TestAccEphemeralRecordByTitle(t *testing.T) {
	secretType := "login"
	_, secretTitle := testAcc.getRecordInfo(secretType)
	if secretTitle == "" {
		t.Fatal("Failed to access test data - missing secret Title")
	}

	config := fmt.Sprintf(`
		ephemeral "secretsmanager_record" "%v" {
			path = "*"
			title = "%v"
		}
	`, secretTitle, secretTitle)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
			},
		},
	})
}

func TestAccEphemeralRecordInvalidPath(t *testing.T) {
	config := `
		ephemeral "secretsmanager_record" "bad" {
			path = "NONEXISTENT_UID_XXXXX"
		}
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`Error reading secret`),
			},
		},
	})
}
