package secretsmanager

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/keeper-security/secrets-manager-go/core"
)

// Provider returns the Keeper Secrets Manager Terraform provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"credential": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KEEPER_CREDENTIAL", nil),
				Description: "Credential to use for Secrets Manager authentication. Can also be sourced from the `KEEPER_CREDENTIAL` environment variable.",
			},
		},
		ConfigureContextFunc: providerConfigure,
		DataSourcesMap: map[string]*schema.Resource{
			"secretsmanager_field":                dataSourceField(),
			"secretsmanager_login":                dataSourceLogin(),
			"secretsmanager_general":              dataSourceGeneral(),
			"secretsmanager_bank_account":         dataSourceBankAccount(),
			"secretsmanager_address":              dataSourceAddress(),
			"secretsmanager_bank_card":            dataSourceBankCard(),
			"secretsmanager_birth_certificate":    dataSourceBirthCertificate(),
			"secretsmanager_contact":              dataSourceContact(),
			"secretsmanager_driver_license":       dataSourceDriverLicense(),
			"secretsmanager_encrypted_notes":      dataSourceEncryptedNotes(),
			"secretsmanager_file":                 dataSourceFile(),
			"secretsmanager_health_insurance":     dataSourceHealthInsurance(),
			"secretsmanager_membership":           dataSourceMembership(),
			"secretsmanager_passport":             dataSourcePassport(),
			"secretsmanager_photo":                dataSourcePhoto(),
			"secretsmanager_server_credentials":   dataSourceServerCredentials(),
			"secretsmanager_software_license":     dataSourceSoftwareLicense(),
			"secretsmanager_ssn_card":             dataSourceSsnCard(),
			"secretsmanager_ssh_keys":             dataSourceSshKeys(),
			"secretsmanager_database_credentials": dataSourceDatabaseCredentials(),
		},
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	creds := d.Get("credential").(string)
	if strings.TrimSpace(creds) == "" {
		return nil, diag.Errorf("empty credential")
	}

	config := core.NewMemoryKeyValueStorage(creds)
	if config.Get(core.KEY_APP_KEY) == "" || config.Get(core.KEY_CLIENT_ID) == "" || config.Get(core.KEY_PRIVATE_KEY) == "" {
		return nil, diag.Errorf("bad credential: %s", creds)
	}

	client := core.NewSecretsManagerFromConfig(config)
	return providerMeta{client}, diags
}

type providerMeta struct {
	client *core.SecretsManager
}

func getTotpCode(totpUrl string) (code string, seconds int, err error) {
	if totp, err := core.GetTotpCode(totpUrl); err == nil {
		return totp.Code, totp.TimeLeft, nil
	} else {
		return "", 0, err
	}
}

func getAddressItemData(secret *core.Record) []interface{} {
	fields := secret.GetFieldsByType("address")
	if len(fields) == 0 {
		return []interface{}{}
	}

	items := []interface{}{}
	if values, ok := fields[0]["value"].([]interface{}); ok && len(values) > 0 {
		if vmap, ok := values[0].(map[string]interface{}); ok {
			item := map[string]interface{}{}
			if val, ok := vmap["street1"].(string); ok {
				item["street1"] = val
			}
			if val, ok := vmap["street2"].(string); ok {
				item["street2"] = val
			}
			if val, ok := vmap["city"].(string); ok {
				item["city"] = val
			}
			if val, ok := vmap["state"].(string); ok {
				item["state"] = val
			}
			if val, ok := vmap["zip"].(string); ok {
				item["zip"] = val
			}
			if val, ok := vmap["country"].(string); ok {
				item["country"] = val
			}
			items = []interface{}{item}
		}
	}

	return items
}

func getAddressRefItemData(secret *core.Record, uid string) []interface{} {
	items := []interface{}{
		map[string]interface{}{
			"uid": uid,
		},
	}
	item := items[0].(map[string]interface{})

	if fields := secret.GetFieldsByType("address"); len(fields) > 0 {
		if values, ok := fields[0]["value"].([]interface{}); ok && len(values) > 0 {
			if vmap, ok := values[0].(map[string]interface{}); ok {
				if val, ok := vmap["street1"].(string); ok {
					item["street1"] = val
				}
				if val, ok := vmap["street2"].(string); ok {
					item["street2"] = val
				}
				if val, ok := vmap["city"].(string); ok {
					item["city"] = val
				}
				if val, ok := vmap["state"].(string); ok {
					item["state"] = val
				}
				if val, ok := vmap["zip"].(string); ok {
					item["zip"] = val
				}
				if val, ok := vmap["country"].(string); ok {
					item["country"] = val
				}
			}
		}
	}
	return items
}

func getBankAccountItemData(secret *core.Record) []interface{} {
	fields := secret.GetFieldsByType("bankAccount")
	if len(fields) == 0 {
		return []interface{}{}
	}

	items := []interface{}{}
	if values, ok := fields[0]["value"].([]interface{}); ok && len(values) > 0 {
		if vmap, ok := values[0].(map[string]interface{}); ok {
			item := map[string]interface{}{}
			if val, ok := vmap["accountType"].(string); ok {
				item["account_type"] = val
			}
			if val, ok := vmap["otherType"].(string); ok {
				item["other_type"] = val
			}
			if val, ok := vmap["routingNumber"].(string); ok {
				item["routing_number"] = val
			}
			if val, ok := vmap["accountNumber"].(string); ok {
				item["account_number"] = val
			}
			items = []interface{}{item}
		}
	}

	return items
}

func getCardRefItemData(secret *core.Record, uid string) []interface{} {
	items := []interface{}{
		map[string]interface{}{
			"uid": uid,
		},
	}
	item := items[0].(map[string]interface{})

	cardItems := getPaymentCardItemData(secret)
	item["payment_card"] = cardItems

	cardholderName := secret.GetFieldValueByType("text")
	item["cardholder_name"] = cardholderName

	pinCode := secret.GetFieldValueByType("pinCode")
	item["pin_code"] = pinCode

	return items
}

func getFileItemsData(fileItems []*core.KeeperFile) []interface{} {
	if len(fileItems) == 0 {
		return []interface{}{}
	}

	fis := make([]interface{}, len(fileItems))

	for i, fileItem := range fileItems {
		fi := map[string]interface{}{}

		fi["uid"] = fileItem.Uid
		fi["title"] = fileItem.Title
		fi["name"] = fileItem.Name
		fi["type"] = fileItem.Type
		fi["size"] = fileItem.Size

		// TF timestamp() uses RFC3339
		timestamp := time.Unix(int64(fileItem.LastModified/1000), 0).Format(time.RFC3339)
		fi["last_modified"] = timestamp
		fi["url"] = fileItem.GetUrl()

		fileData := fileItem.GetFileData()
		fi["content_base64"] = base64.StdEncoding.EncodeToString(fileData)

		fis[i] = fi
	}

	return fis
}

func getHostItemData(secret *core.Record) []interface{} {
	fields := secret.GetFieldsByType("host")
	if len(fields) == 0 {
		return []interface{}{}
	}

	items := []interface{}{}
	if values, ok := fields[0]["value"].([]interface{}); ok && len(values) > 0 {
		if vmap, ok := values[0].(map[string]interface{}); ok {
			item := map[string]interface{}{}
			if val, ok := vmap["hostName"].(string); ok {
				item["host_name"] = val
			}
			if val, ok := vmap["port"].(string); ok {
				item["port"] = val
			}
			items = []interface{}{item}
		}
	}

	return items
}

func getKeyPairItemData(secret *core.Record) []interface{} {
	fields := secret.GetFieldsByType("keyPair")
	if len(fields) == 0 {
		return []interface{}{}
	}

	items := []interface{}{}
	if values, ok := fields[0]["value"].([]interface{}); ok && len(values) > 0 {
		if vmap, ok := values[0].(map[string]interface{}); ok {
			item := map[string]interface{}{}
			if val, ok := vmap["publicKey"].(string); ok {
				item["public_key"] = val
			}
			if val, ok := vmap["privateKey"].(string); ok {
				item["private_key"] = val
			}
			items = []interface{}{item}
		}
	}

	return items
}

func getNameItemData(secret *core.Record) []interface{} {
	fields := secret.GetFieldsByType("name")
	if len(fields) == 0 {
		return []interface{}{}
	}

	items := []interface{}{}
	if values, ok := fields[0]["value"].([]interface{}); ok && len(values) > 0 {
		if vmap, ok := values[0].(map[string]interface{}); ok {
			item := map[string]interface{}{}
			if val, ok := vmap["first"].(string); ok {
				item["first"] = val
			}
			if val, ok := vmap["middle"].(string); ok {
				item["middle"] = val
			}
			if val, ok := vmap["last"].(string); ok {
				item["last"] = val
			}
			items = []interface{}{item}
		}
	}

	return items
}

func getPaymentCardItemData(secret *core.Record) []interface{} {
	fields := secret.GetFieldsByType("paymentCard")
	if len(fields) == 0 {
		return []interface{}{}
	}

	items := []interface{}{}
	if values, ok := fields[0]["value"].([]interface{}); ok && len(values) > 0 {
		if vmap, ok := values[0].(map[string]interface{}); ok {
			item := map[string]interface{}{}
			if val, ok := vmap["cardNumber"].(string); ok {
				item["card_number"] = val
			}
			if val, ok := vmap["cardExpirationDate"].(string); ok {
				item["card_expiration_date"] = val
			}
			if val, ok := vmap["cardSecurityCode"].(string); ok {
				item["card_security_code"] = val
			}
			items = []interface{}{item}
		}
	}

	return items
}

func getPhoneItemData(secret *core.Record) []interface{} {
	fields := secret.GetFieldsByType("phone")
	if len(fields) == 0 {
		return []interface{}{}
	}

	items := []interface{}{}
	if values, ok := fields[0]["value"].([]interface{}); ok && len(values) > 0 {
		if vmap, ok := values[0].(map[string]interface{}); ok {
			item := map[string]interface{}{}
			if val, ok := vmap["region"].(string); ok {
				item["region"] = val
			}
			if val, ok := vmap["number"].(string); ok {
				item["number"] = val
			}
			if val, ok := vmap["ext"].(string); ok {
				item["ext"] = val
			}
			if val, ok := vmap["type"].(string); ok {
				item["type"] = val
			}
			items = []interface{}{item}
		}
	}

	return items
}

func getRecord(path string, title string, client core.SecretsManager) (secret *core.Record, e error) {
	defer func() {
		if r := recover(); r != nil {
			secret = nil
			switch x := r.(type) {
			case string:
				e = errors.New(x)
			case error:
				e = x
			default:
				e = fmt.Errorf("error in provider - getRecord: %v", r)
			}
		}
	}()

	secret = nil
	title = strings.TrimSpace(title)
	path = strings.TrimSpace(path)
	if title != "" && path == "*" { // find by title requested
		secrets, err := client.GetSecrets([]string{})
		if err != nil {
			return nil, err
		}
		if len(secrets) == 0 {
			return nil, fmt.Errorf("record not found - title: %s", title)
		}
		for _, r := range secrets {
			if r.Title() == title {
				if secret == nil {
					secret = r
				} else {
					return secret, fmt.Errorf("more that one records match the search query - title: %s", title)
				}
			}
		}
		if secret == nil {
			return nil, fmt.Errorf("record not found - title: %s", title)
		}
		return secret, nil
	} else {
		secrets, err := client.GetSecrets([]string{path})
		if err != nil {
			return nil, err
		}
		if len(secrets) == 0 {
			return nil, fmt.Errorf("record not found - UID: %s", path)
		}
		if len(secrets) > 1 {
			return nil, fmt.Errorf("expected 1 record - found %d records for UID: %s", len(secrets), path)
		}
		return secrets[0], nil
	}
}
