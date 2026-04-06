package secretsmanager

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/keeper-security/secrets-manager-go/core"
)

const (
	envCredential = "KEEPER_CREDENTIAL"
)

var testAccProtoV6ProviderFactories map[string]func() (tfprotov6.ProviderServer, error)
var testAcc *testAccValues

// testAccKSMClient is a lazily-initialized standalone KSM client for use in
// test helpers and PreConfig hooks. It is independent of the provider lifecycle
// so it works correctly when tests use testAccProtoV6ProviderFactories.
var (
	testAccKSMClient     *core.SecretsManager
	testAccKSMClientOnce sync.Once
)

type testAccValues struct {
	credential string
	secrets    map[string]map[string]string
	folderUid  string
}

func (testAccValues) validate() error {
	if testAcc.credential == "" {
		return fmt.Errorf("make sure you set environment variables: %s", envCredential)
	}
	return nil
}

func init() {
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"secretsmanager": func() (tfprotov6.ProviderServer, error) {
			ctx := context.Background()
			sdkv2Provider := Provider()
			upgradedSdkv2, err := tf5to6server.UpgradeServer(ctx, sdkv2Provider.GRPCProvider)
			if err != nil {
				return nil, err
			}
			servers := []func() tfprotov6.ProviderServer{
				func() tfprotov6.ProviderServer { return upgradedSdkv2 },
				providerserver.NewProtocol6(NewFWProvider()),
			}
			muxServer, err := tf6muxserver.NewMuxServer(ctx, servers...)
			if err != nil {
				return nil, err
			}
			return muxServer.ProviderServer(), nil
		},
	}

	testAcc = &testAccValues{
		credential: os.Getenv(envCredential),
		secrets: map[string]map[string]string{
			"address":             {"uid": "*", "title": "tf_acc_test_address"},
			"bankAccount":         {"uid": "*", "title": "tf_acc_test_bank_account"},
			"bankCard":            {"uid": "*", "title": "tf_acc_test_bank_card"},
			"birthCertificate":    {"uid": "*", "title": "tf_acc_test_birth_certificate"},
			"contact":             {"uid": "*", "title": "tf_acc_test_contact"},
			"databaseCredentials": {"uid": "*", "title": "tf_acc_test_database_credentials"},
			"driverLicense":       {"uid": "*", "title": "tf_acc_test_driver_license"},
			"encryptedNotes":      {"uid": "*", "title": "tf_acc_test_encrypted_notes"},
			"field":               {"uid": "*/field/login", "title": "tf_acc_test_field"},
			"file":                {"uid": "*", "title": "tf_acc_test_file"},
			"healthInsurance":     {"uid": "*", "title": "tf_acc_test_health_insurance"},
			"login":               {"uid": "*", "title": "tf_acc_test_login"},
			"membership":          {"uid": "*", "title": "tf_acc_test_membership"},
			"passport":            {"uid": "*", "title": "tf_acc_test_passport"},
			"photo":               {"uid": "*", "title": "tf_acc_test_photo"},
			"serverCredentials":   {"uid": "*", "title": "tf_acc_test_server_credentials"},
			"softwareLicense":     {"uid": "*", "title": "tf_acc_test_software_license"},
			"sshKeys":             {"uid": "*", "title": "tf_acc_test_ssh_keys"},
			"ssnCard":             {"uid": "*", "title": "tf_acc_test_ssn_card"},
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

func (testAccValues) getTestFolder() string {
	folderUid := strings.TrimSpace(testAcc.folderUid)
	if folderUid == "" || folderUid == "*" {
		creds := strings.TrimSpace(testAcc.credential)
		if creds != "" {
			config := core.NewMemoryKeyValueStorage(creds)
			if config.Get(core.KEY_APP_KEY) != "" && config.Get(core.KEY_CLIENT_ID) != "" && config.Get(core.KEY_PRIVATE_KEY) != "" {
				client := core.NewSecretsManager(&core.ClientOptions{Config: config})
				if fuid, err := getTemplateFolder(folderUid, *client); err == nil && fuid != "" {
					testAcc.folderUid = fuid
				}
			}
		}
	}
	return testAcc.folderUid
}

// testAccClient returns a lazily-initialized standalone KSM client built
// directly from KEEPER_CREDENTIAL. It does not depend on the Terraform
// provider being configured, so it works correctly in tests that use
// testAccProtoV6ProviderFactories.
func testAccClient() *core.SecretsManager {
	testAccKSMClientOnce.Do(func() {
		creds := strings.TrimSpace(testAcc.credential)
		if creds == "" {
			return
		}
		config := core.NewMemoryKeyValueStorage(creds)
		if config.Get(core.KEY_APP_KEY) == "" || config.Get(core.KEY_CLIENT_ID) == "" || config.Get(core.KEY_PRIVATE_KEY) == "" {
			return
		}
		testAccKSMClient = core.NewSecretsManager(&core.ClientOptions{Config: config})
	})
	return testAccKSMClient
}

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

func checkSecretExistsRemotely(uid string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccClient()
		if client == nil {
			return fmt.Errorf("cannot create KSM client from credentials")
		}

		records, err := getSecrets(*client, []string{uid})
		if err != nil {
			return err
		}
		if len(records) == 0 {
			return fmt.Errorf("resource '%v' doesn't exist remotely", uid)
		}

		return nil
	}
}

func checkFolderExistsRemotely(uid, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccClient()
		if client == nil {
			return fmt.Errorf("cannot create KSM client from credentials")
		}

		folders, err := getFolders(*client)
		if err != nil {
			return err
		}

		for _, f := range folders {
			if (uid != "" && f.FolderUid == uid) || (name != "" && name == f.Name) {
				return nil
			}
		}

		return fmt.Errorf("resource '%v' doesn't exist remotely", uid)
	}
}

func checkSecretResourceState(resourceName string, check func(s *terraform.InstanceState) error) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceState := s.RootModule().Resources[resourceName]
		if resourceState == nil {
			return fmt.Errorf("resource '%v' not in tf state", resourceName)
		}

		state := resourceState.Primary
		if state == nil {
			return fmt.Errorf("resource has no primary instance")
		}

		return check(state)
	}
}
