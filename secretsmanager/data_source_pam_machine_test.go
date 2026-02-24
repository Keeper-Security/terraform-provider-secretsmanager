package secretsmanager

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/keeper-security/secrets-manager-go/core"
)

func TestAccDataSourcePamMachine(t *testing.T) {
	secretFolderUid := testAcc.getTestFolder()
	secretUid := core.GenerateUid()
	secretTitle := "tf_acc_test_ds_pam_machine"
	if secretFolderUid == "" {
		t.Skip("Skipping test - TF_ACC not set or test folder not configured")
	}

	config := fmt.Sprintf(`
		resource "secretsmanager_pam_machine" "%v" {
			folder_uid = "%v"
			uid        = "%v"
			title      = "%v"
			notes      = "%v"
			pam_hostname {
				value {
					hostname = "192.168.1.100"
					port     = "22"
				}
			}
			pam_settings = jsonencode([{
				connection = [{
					protocol = "ssh"
					port = "22"
				}]
			}])
			login {
				value = "svc_machine"
			}
			password {
				value = "StrongMachinePass123!"
			}
			operating_system {
				label = "Operating System"
				value = "Linux"
			}
			private_pem_key {
				value = "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA..."
			}
			private_key_passphrase {
				value = "TestPassphrase#123"
			}
			ssl_verification {
				label = "SSL Verification"
				value = true
			}
			totp {
				value = "otpauth://totp/keeper:machine?secret=JBSWY3DPEHPK3PXP&issuer=Keeper"
			}
			rotation_scripts {
				value {
					command = "echo hello"
				}
			}
			instance_name {
				label = "Instance Name"
				value = "web-server-01"
			}
			instance_id {
				label = "Instance Id"
				value = "i-1234567890abcdef0"
			}
			provider_group {
				label = "Provider Group"
				value = "production-servers"
			}
			provider_region {
				label = "Provider Region"
				value = "us-east-1"
			}
		}

		data "secretsmanager_pam_machine" "%v" {
			path = secretsmanager_pam_machine.%v.uid
		}
	`, secretTitle, secretFolderUid, secretUid, secretTitle, secretTitle,
		secretTitle, secretTitle)

	dataName := fmt.Sprintf("data.secretsmanager_pam_machine.%v", secretTitle)
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataName, "type", "pamMachine"),
					resource.TestCheckResourceAttr(dataName, "title", secretTitle),
					resource.TestCheckResourceAttr(dataName, "notes", secretTitle),
					resource.TestCheckResourceAttr(dataName, "folder_uid", secretFolderUid),
					resource.TestCheckResourceAttr(dataName, "pam_hostname.0.value.0.hostname", "192.168.1.100"),
					resource.TestCheckResourceAttr(dataName, "pam_hostname.0.value.0.port", "22"),
					resource.TestCheckResourceAttrSet(dataName, "pam_settings"),
					resource.TestCheckResourceAttr(dataName, "login.0.value", "svc_machine"),
					resource.TestCheckResourceAttr(dataName, "password.0.value", "StrongMachinePass123!"),
					resource.TestCheckResourceAttr(dataName, "operating_system.0.value", "Linux"),
					resource.TestCheckResourceAttr(dataName, "private_pem_key.0.value", "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA..."),
					resource.TestCheckResourceAttr(dataName, "private_key_passphrase.0.value", "TestPassphrase#123"),
					resource.TestCheckResourceAttr(dataName, "ssl_verification.0.value", "true"),
					resource.TestCheckResourceAttr(dataName, "rotation_scripts.0.value.0.command", "echo hello"),
					resource.TestCheckResourceAttr(dataName, "instance_name.0.value", "web-server-01"),
					resource.TestCheckResourceAttr(dataName, "instance_id.0.value", "i-1234567890abcdef0"),
					resource.TestCheckResourceAttr(dataName, "provider_group.0.value", "production-servers"),
					resource.TestCheckResourceAttr(dataName, "provider_region.0.value", "us-east-1"),
					resource.TestCheckResourceAttr(dataName, "totp.0.value", "otpauth://totp/keeper:machine?secret=JBSWY3DPEHPK3PXP&issuer=Keeper"),
				),
			},
		},
	})
}
