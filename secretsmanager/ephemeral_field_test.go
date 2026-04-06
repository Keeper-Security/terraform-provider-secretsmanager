package secretsmanager

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEphemeralField(t *testing.T) {
	secretType := "field"
	secretUid, secretTitle := testAcc.getRecordInfo(secretType)
	if secretUid == "" || secretTitle == "" {
		t.Fatal("Failed to access test data - missing secret UID and/or Title")
	}

	config := fmt.Sprintf(`
		ephemeral "secretsmanager_field" "%v" {
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

func TestAccEphemeralFieldByTitle(t *testing.T) {
	secretType := "field"
	_, secretTitle := testAcc.getRecordInfo(secretType)
	if secretTitle == "" {
		t.Fatal("Failed to access test data - missing secret Title")
	}

	config := fmt.Sprintf(`
		ephemeral "secretsmanager_field" "%v" {
			path = "*/field/login"
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

func TestAccEphemeralFieldInvalidPath(t *testing.T) {
	config := `
		ephemeral "secretsmanager_field" "bad" {
			path = "NONEXISTENT_UID_12345/field/login"
		}
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`Error reading field`),
			},
		},
	})
}
