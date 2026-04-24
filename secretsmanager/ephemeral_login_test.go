package secretsmanager

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEphemeralLogin(t *testing.T) {
	secretType := "login"
	secretUid, secretTitle := testAcc.getRecordInfo(secretType)
	if secretUid == "" || secretTitle == "" {
		t.Fatal("Failed to access test data - missing secret UID and/or Title")
	}

	config := fmt.Sprintf(`
		ephemeral "secretsmanager_login" "%v" {
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

func TestAccEphemeralLoginByTitle(t *testing.T) {
	secretType := "login"
	_, secretTitle := testAcc.getRecordInfo(secretType)
	if secretTitle == "" {
		t.Fatal("Failed to access test data - missing secret Title")
	}

	config := fmt.Sprintf(`
		ephemeral "secretsmanager_login" "%v" {
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

func TestAccEphemeralLoginWrongType(t *testing.T) {
	// Use a bankCard record UID with the login ephemeral resource — should error
	secretUid, secretTitle := testAcc.getRecordInfo("bankCard")
	if secretUid == "" || secretTitle == "" {
		t.Fatal("Failed to access test data - missing bankCard UID and/or Title")
	}

	config := fmt.Sprintf(`
		ephemeral "secretsmanager_login" "wrong_type" {
			path = "%v"
			title = "%v"
		}
	`, secretUid, secretTitle)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`Record Type Mismatch`),
			},
		},
	})
}
