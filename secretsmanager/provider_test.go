package secretsmanager

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	envCredential = "KEEPER_CREDENTIAL"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider
var testAcc *testAccValues

type testAccValues struct {
	credential string
	secrets    map[string]map[string]string
}

func (testAccValues) validate() error {
	if testAcc.credential == "" {
		return fmt.Errorf("make sure you set environment variables: %s", envCredential)
	}
	return nil
}

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"secretsmanager": testAccProvider,
	}

	testAcc = &testAccValues{
		credential: os.Getenv(envCredential),
		secrets: map[string]map[string]string{
			"field":               {"uid": "*/field/login", "title": "tf_acc_test_field"},
			"login":               {"uid": "*", "title": "tf_acc_test_login"},
			"general":             {"uid": "*", "title": "tf_acc_test_general"},
			"bankAccount":         {"uid": "*", "title": "tf_acc_test_bank_account"},
			"address":             {"uid": "*", "title": "tf_acc_test_address"},
			"bankCard":            {"uid": "*", "title": "tf_acc_test_bank_card"},
			"birthCertificate":    {"uid": "*", "title": "tf_acc_test_birth_certificate"},
			"contact":             {"uid": "*", "title": "tf_acc_test_contact"},
			"driverLicense":       {"uid": "*", "title": "tf_acc_test_driver_license"},
			"encryptedNotes":      {"uid": "*", "title": "tf_acc_test_encrypted_notes"},
			"file":                {"uid": "*", "title": "tf_acc_test_file"},
			"healthInsurance":     {"uid": "*", "title": "tf_acc_test_health_insurance"},
			"membership":          {"uid": "*", "title": "tf_acc_test_membership"},
			"passport":            {"uid": "*", "title": "tf_acc_test_passport"},
			"photo":               {"uid": "*", "title": "tf_acc_test_photo"},
			"serverCredentials":   {"uid": "*", "title": "tf_acc_test_server_credentials"},
			"softwareLicense":     {"uid": "*", "title": "tf_acc_test_software_license"},
			"ssnCard":             {"uid": "*", "title": "tf_acc_test_ssn_card"},
			"sshKeys":             {"uid": "*", "title": "tf_acc_test_ssh_keys"},
			"databaseCredentials": {"uid": "*", "title": "tf_acc_test_database_credentials"},
		},
	}
}

func (testAccValues) getRecordInfo(recordType string) (uid string, title string) {
	if secret, ok := testAcc.secrets[recordType]; ok {
		if uid, ok = secret["uid"]; ok {
			if title, ok = secret["title"]; ok {
				return
			}
		}
	}
	return "", ""
}

// func client() *ksm.SecretsManager {
// 	return testAccProvider.Meta().(providerMeta).client
// }

func testAccPreCheck(t *testing.T) func() {
	return func() {
		err := testAcc.validate()
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
