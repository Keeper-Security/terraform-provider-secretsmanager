package secretsmanager

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
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
			"secretsmanager_address":              dataSourceAddress(),
			"secretsmanager_bank_account":         dataSourceBankAccount(),
			"secretsmanager_bank_card":            dataSourceBankCard(),
			"secretsmanager_birth_certificate":    dataSourceBirthCertificate(),
			"secretsmanager_contact":              dataSourceContact(),
			"secretsmanager_database_credentials": dataSourceDatabaseCredentials(),
			"secretsmanager_driver_license":       dataSourceDriverLicense(),
			"secretsmanager_encrypted_notes":      dataSourceEncryptedNotes(),
			"secretsmanager_field":                dataSourceField(),
			"secretsmanager_file":                 dataSourceFile(),
			"secretsmanager_general":              dataSourceGeneral(),
			"secretsmanager_health_insurance":     dataSourceHealthInsurance(),
			"secretsmanager_login":                dataSourceLogin(),
			"secretsmanager_membership":           dataSourceMembership(),
			"secretsmanager_passport":             dataSourcePassport(),
			"secretsmanager_photo":                dataSourcePhoto(),
			"secretsmanager_server_credentials":   dataSourceServerCredentials(),
			"secretsmanager_software_license":     dataSourceSoftwareLicense(),
			"secretsmanager_ssh_keys":             dataSourceSshKeys(),
			"secretsmanager_ssn_card":             dataSourceSsnCard(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"secretsmanager_address":              resourceAddress(),
			"secretsmanager_bank_account":         resourceBankAccount(),
			"secretsmanager_bank_card":            resourceBankCard(),
			"secretsmanager_birth_certificate":    resourceBirthCertificate(),
			"secretsmanager_contact":              resourceContact(),
			"secretsmanager_database_credentials": resourceDatabaseCredentials(),
			"secretsmanager_driver_license":       resourceDriverLicense(),
			"secretsmanager_encrypted_notes":      resourceEncryptedNotes(),
			"secretsmanager_file":                 resourceFile(),
			"secretsmanager_health_insurance":     resourceHealthInsurance(),
			"secretsmanager_login":                resourceLogin(),
			"secretsmanager_membership":           resourceMembership(),
			"secretsmanager_passport":             resourcePassport(),
			"secretsmanager_photo":                resourcePhoto(),
			"secretsmanager_server_credentials":   resourceServerCredentials(),
			"secretsmanager_software_license":     resourceSoftwareLicense(),
			"secretsmanager_ssh_keys":             resourceSshKeys(),
			"secretsmanager_ssn_card":             resourceSsnCard(),
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

	client := core.NewSecretsManager(&core.ClientOptions{Config: config})
	return providerMeta{client}, diags
}

type providerMeta struct {
	client *core.SecretsManager
}

// map attribute names from schema to field types in record v3
var mapSchemaToRecordFieldName map[string]string = map[string]string{
	"account_number":    "accountNumber", // Text
	"address":           "address",
	"address_ref":       "addressRef",
	"bank_account":      "bankAccount",
	"birth_date":        "birthDate",
	"card_ref":          "cardRef",
	"company":           "text", // Text
	"date":              "date",
	"email":             "email",
	"expiration_date":   "expirationDate",
	"file_ref":          "fileRef",
	"group_number":      "groupNumber", // Text
	"host":              "host",
	"key_pair":          "keyPair",
	"license_number":    "licenseNumber",
	"login":             "login",
	"multiline":         "multiline",
	"name":              "name",
	"note":              "note", // SecureNote
	"one_time_code":     "oneTimeCode",
	"password":          "password",
	"payment_card":      "paymentCard",
	"phone":             "phone",
	"pin_code":          "pinCode",
	"secret":            "secret",
	"security_question": "securityQuestion",
	"text":              "text",
	"title":             "title", // Text
	"url":               "url",
	// schema attributes that use field label instead of field type
	// "company":            "text",          // contact
	"cardholder_name":       "text",          // bankCard
	"db_type":               "text",          // databaseCredentials
	"driver_license_number": "accountNumber", // driverLicense
	"totp":                  "oneTimeCode",   // login/general, bankAccount
	"passport_number":       "accountNumber", // passport
	"date_issued":           "date",          // passport
	"activation_date":       "date",          // softwareLicense
	"passphrase":            "password",      // sshKeys
	"identity_number":       "accountNumber", // ssnCard
}

/*
var mapFieldTypeToFieldValueType map[string]string = map[string]string{
	"accountNumber":    "text", // Lookup: accountNumber
	"address":          "address",
	"addressRef":       "addressRef",  // Lookup: addressRef
	"bankAccount":      "bankAccount", // Lookup: accountNumber
	"birthDate":        "date",        // stored as unix milliseconds
	"cardRef":          "cardRef",     // Lookup: bankCard, Multiple: default
	"company":          "text",        // Lookup: company
	"date":             "date",
	"email":            "email", // Lookup: email, Multiple: optional
	"expirationDate":   "date",
	"fileRef":          "fileRef", // Multiple: default
	"groupNumber":      "text",
	"host":             "host", // Lookup: host, Multiple: optional
	"keyPair":          "privateKey",
	"licenseNumber":    "multiline",
	"login":            "login", // Lookup: login
	"multiline":        "multiline",
	"name":             "name", // Lookup: name
	"note":             "multiline",
	"oneTimeCode":      "otp",
	"password":         "password",
	"paymentCard":      "paymentCard",
	"phone":            "phone", // Lookup: phone, Multiple: optional
	"pinCode":          "secret",
	"secret":           "text",
	"securityQuestion": "securityQuestion", // Multiple: default
	"text":             "text",
	"title":            "text",
	"url":              "url", // Multiple: optional
}
*/

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
		// fi["url"] = fileItem.GetUrl() // use content_base64 to access file content

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

func getFieldDicts(fieldType, section string, recordDict map[string]interface{}) []interface{} {
	result := []interface{}{}
	if flds, found := recordDict[section]; found {
		if fi, ok := flds.([]interface{}); ok {
			for _, fiv := range fi {
				if fmap, ok := fiv.(map[string]interface{}); ok {
					if fv, found := fmap["type"]; found && fv == fieldType {
						result = append(result, fmap)
					}
				}
			}
		}
	}
	return result
}

func getFieldResourceData(fieldType, section string, secret *core.Record) interface{} {
	if flds := getFieldDicts(fieldType, section, secret.RecordDict); len(flds) > 0 {
		if fieldSchema := parseFieldFromDataJson(flds[0]); fieldSchema != nil {
			ftSchema := map[string]interface{}{
				"type": fieldType,
			}
			if _, found := fieldSchema.fields["label"]; found {
				ftSchema["label"] = fieldSchema.Label
			}
			if _, found := fieldSchema.fields["required"]; found {
				ftSchema["required"] = fieldSchema.Required
			}
			if _, found := fieldSchema.fields["privacyScreen"]; found {
				ftSchema["privacy_screen"] = fieldSchema.PrivacyScreen
			}
			if _, found := fieldSchema.fields["enforceGeneration"]; found {
				ftSchema["enforce_generation"] = fieldSchema.EnforceGeneration
			}
			if _, found := fieldSchema.fields["complexity"]; found && len(fieldSchema.Complexity) > 0 {
				complexity := fieldSchema.Complexity[0]
				cmap := map[string]interface{}{
					"length":    complexity.Length,
					"caps":      complexity.Caps,
					"lowercase": complexity.Lowercase,
					"digits":    complexity.Digits,
					"special":   complexity.Special,
				}
				ftSchema["complexity"] = []interface{}{cmap}
			}
			if _, found := fieldSchema.fields["value"]; found {
				fis := []interface{}{}
				if fi, ok := fieldSchema.Value.([]interface{}); ok {
					for _, fiv := range fi {
						if str, ok := fiv.(string); ok && str != "" {
							// simple value - string
							ftSchema["value"] = str
						} else if num, ok := fiv.(float64); ok {
							// simple value - Int64 (converted to float/float64 by JSON)
							ftSchema["value"] = int64(num)
						} else if fmap, ok := fiv.(map[string]interface{}); ok && len(fmap) > 0 {
							// complex value - map struct fields to schema
							fv := map[string]interface{}{}
							switch fieldType {
							case "address":
								if str, ok := fmap["street1"]; ok {
									fv["street1"] = str
								}
								if str, ok := fmap["street2"]; ok {
									fv["street2"] = str
								}
								if str, ok := fmap["city"]; ok {
									fv["city"] = str
								}
								if str, ok := fmap["state"]; ok {
									fv["state"] = str
								}
								if str, ok := fmap["country"]; ok {
									fv["country"] = str
								}
								if str, ok := fmap["zip"]; ok {
									fv["zip"] = str
								}
							case "bankAccount":
								if str, ok := fmap["accountType"]; ok {
									fv["account_type"] = str
								}
								if str, ok := fmap["routingNumber"]; ok {
									fv["routing_number"] = str
								}
								if str, ok := fmap["accountNumber"]; ok {
									fv["account_number"] = str
								}
								if str, ok := fmap["otherType"]; ok {
									fv["other_type"] = str
								}
							case "host":
								if str, ok := fmap["hostName"]; ok {
									fv["host_name"] = str
								}
								if str, ok := fmap["port"]; ok {
									fv["port"] = str
								}
							case "keyPair":
								if str, ok := fmap["publicKey"]; ok {
									fv["public_key"] = str
								}
								if str, ok := fmap["privateKey"]; ok {
									fv["private_key"] = str
								}
							case "name":
								if str, ok := fmap["first"]; ok {
									fv["first"] = str
								}
								if str, ok := fmap["middle"]; ok {
									fv["middle"] = str
								}
								if str, ok := fmap["last"]; ok {
									fv["last"] = str
								}
							case "paymentCard":
								if str, ok := fmap["cardNumber"]; ok {
									fv["card_number"] = str
								}
								if str, ok := fmap["cardExpirationDate"]; ok {
									fv["card_expiration_date"] = str
								}
								if str, ok := fmap["cardSecurityCode"]; ok {
									fv["card_security_code"] = str
								}
							case "phone":
								if str, ok := fmap["region"]; ok {
									fv["region"] = str
								}
								if str, ok := fmap["number"]; ok {
									fv["number"] = str
								}
								if str, ok := fmap["ext"]; ok {
									fv["ext"] = str
								}
								if str, ok := fmap["type"]; ok {
									fv["type"] = str
								}
							case "securityQuestion":
								if str, ok := fmap["question"]; ok {
									fv["question"] = str
								}
								if str, ok := fmap["answer"]; ok {
									fv["answer"] = str
								}
							default:
								fv = nil
							}
							if len(fv) > 0 {
								fis = append(fis, fv)
							}
						}
					}
				}
				if len(fis) > 0 {
					ftSchema["value"] = []interface{}{fis[0]}
				}
			}
			return []interface{}{ftSchema}
		}
	}
	return []interface{}{}
}

func getFileItemsResourceData(secret *core.Record) []interface{} {
	if flds := getFieldDicts("fileRef", "fields", secret.RecordDict); len(flds) > 0 {
		fieldSchema := parseFieldFromDataJson(flds[0])
		ftsFileref := map[string]interface{}{
			"type": "fileRef",
		}
		if _, found := fieldSchema.fields["label"]; found {
			ftsFileref["label"] = fieldSchema.Label
		}
		if _, found := fieldSchema.fields["required"]; found {
			ftsFileref["required"] = fieldSchema.Required
		}
		if _, found := fieldSchema.fields["value"]; found {
			fis := []interface{}{}
			if fi, ok := fieldSchema.Value.([]interface{}); ok {
				for _, fiv := range fi {
					if fref, ok := fiv.(string); ok && fref != "" {
						fi := map[string]interface{}{"uid": fref}
						for _, fileItem := range secret.Files {
							if fref == fileItem.Uid {
								// fi["uid"] = fileItem.Uid
								fi["title"] = fileItem.Title
								fi["name"] = fileItem.Name
								fi["type"] = fileItem.Type
								fi["size"] = fileItem.Size

								// TF timestamp() uses RFC3339
								timestamp := time.Unix(int64(fileItem.LastModified/1000), 0).Format(time.RFC3339)
								fi["last_modified"] = timestamp
								// fi["url"] = fileItem.GetUrl() // use content_base64 to access file content

								fileData := fileItem.GetFileData()
								fi["content_base64"] = base64.StdEncoding.EncodeToString(fileData)
								break
							}
						}
						fis = append(fis, fi)
					}
				}
			}
			if len(fis) > 0 {
				ftsFileref["value"] = fis
			}
		}
		return []interface{}{ftsFileref}
	}
	return []interface{}{}
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

/*
// deprecated - use NewRecordCreate
func getTemplateRecord(folderUid string, recordType string, templateTitle string, client core.SecretsManager) (secret *core.Record, e error) {
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

	e = nil
	secret = nil
	folderUid = strings.TrimSpace(folderUid)
	recordType = strings.TrimSpace(recordType)
	if folderUid != "" && recordType != "" {
		secrets, err := client.GetSecrets([]string{})
		if err != nil {
			return nil, err
		}
		// find folder by looking up record with templateTitle
		if folderUid == "*" && templateTitle != "" {
			for _, r := range secrets {
				if r.Title() == templateTitle {
					folderUid = r.FolderUid()
					break
				}
			}
		}
		// lookup template record of specified type in requested folder
		for _, rec := range secrets {
			typeMatch := rec.Type() == recordType
			folderMatch := rec.FolderUid() == folderUid
			if typeMatch && folderMatch {
				secret = rec
				break
			} else if folderMatch && secret == nil {
				secret = rec // different record type but in requested folder
			}
		}
	}

	if secret == nil {
		e = fmt.Errorf("template record not found")
	}

	return secret, e
}
*/

func getTemplateFolder(folderUid string, client core.SecretsManager) (fuid string, e error) {
	defer func() {
		if r := recover(); r != nil {
			fuid = ""
			switch x := r.(type) {
			case string:
				e = errors.New(x)
			case error:
				e = x
			default:
				e = fmt.Errorf("error in provider - getTemplateFolder: %v", r)
			}
		}
	}()

	e = nil
	fuid = ""
	folderUid = strings.TrimSpace(folderUid)
	if folderUid == "" || folderUid == "*" {
		secrets, err := client.GetSecrets([]string{})
		if err != nil {
			return "", err
		}
		for _, r := range secrets {
			if r.FolderUid() != "" {
				fuid = r.FolderUid()
				break
			}
		}
	} else {
		fuid = folderUid
	}

	if fuid == "" || fuid == "*" {
		e = fmt.Errorf("template folder not found")
	}

	return fuid, e
}

// getStringListData splits a string into list using the separator and skipping empty parts
func GetStringListData(data string, separator string) []interface{} {
	if data == "" {
		return []interface{}{}
	}

	if separator == "" {
		separator = " "
	}

	stringItems := strings.Split(data, separator)

	// remove empty parts and convert to interfaces
	items := make([]interface{}, 0, len(stringItems))
	for _, v := range stringItems {
		if v != "" {
			items = append(items, v)
		}
	}

	return items
}

// converts string to bool, returns default value if conversion fails
func StrToBoolDef(boolString string, defaultValue bool) bool {
	if value, err := strconv.ParseBool(boolString); err == nil {
		return value
	}
	return defaultValue
}

type genericFieldSchema struct {
	fields            map[string]struct{}
	Type              string                     `json:"type"`
	Label             string                     `json:"label,omitempty"`
	Required          bool                       `json:"required,omitempty"`
	PrivacyScreen     bool                       `json:"privacy_screen,omitempty"`
	EnforceGeneration bool                       `json:"enforce_generation,omitempty"`
	Complexity        []*core.PasswordComplexity `json:"complexity,omitempty"`
	Value             interface{}                `json:"value,omitempty"`
	ValueString       []string
}

// note: in schema Complexity is TypeList with MaxItems=1 but in record it is just a single object
type genericFieldJson struct {
	fields            map[string]struct{}
	Type              string                   `json:"type"`
	Label             string                   `json:"label,omitempty"`
	Required          bool                     `json:"required,omitempty"`
	PrivacyScreen     bool                     `json:"privacyScreen,omitempty"`
	EnforceGeneration bool                     `json:"enforceGeneration,omitempty"`
	Complexity        *core.PasswordComplexity `json:"complexity,omitempty"`
	Value             interface{}              `json:"value,omitempty"`
	ValueString       []string
}

func GetGenericFieldSchemaValueInt64(field *genericFieldSchema) (value int64, ok bool) {
	// for simple field values of type schema.TypeInt
	if val, ok := field.Value.(int64); ok {
		return val, true
	}
	// JSON stores numbers as float - interface {}(float64)
	if val, ok := field.Value.(float64); ok {
		return int64(val), true
	}
	return 0, false
}

func GetGenericFieldSchemaValueString(field *genericFieldSchema) (value string, ok bool) {
	// for simple field values of type schema.TypeString vs ex. fileRef ["UID1", "UID2"]
	if str, ok := field.Value.(string); ok {
		return str, true
	}
	return "", false
}

func GetGenericFieldSchemaValueStringList(field *genericFieldSchema) (value []string, ok bool) {
	// file_ref values are schema.TypeList mapped to fileRef ["UID1", "UID2"]
	if si, ok := field.Value.([]interface{}); ok && len(si) > 0 {
		values := []string{}
		for _, sii := range si {
			if msi, ok := sii.(map[string]interface{}); ok && len(msi) > 0 {
				if vi, found := msi["uid"]; found {
					if val, ok := vi.(string); ok && val != "" {
						values = append(values, val)
					}
				}
			}
		}
		return values, true
	}
	return nil, false
}

func GetGenericFieldSchemaValue(fieldType string, field *genericFieldSchema) []interface{} {
	// for complex fields return the expected field type object as interface{}
	result := []interface{}{}

	if _, found := field.fields["value"]; !found || field.Value == nil {
		return nil
	}

	if si, ok := field.Value.([]interface{}); ok && len(si) > 0 {
		for _, sii := range si {
			if msi, ok := sii.(map[string]interface{}); ok && len(msi) > 0 {
				switch fieldType {
				case "address":
					address := core.Address{}
					if v, found := msi["street1"]; found {
						address.Street1 = v.(string)
					}
					if v, found := msi["street2"]; found {
						address.Street2 = v.(string)
					}
					if v, found := msi["city"]; found {
						address.City = v.(string)
					}
					if v, found := msi["state"]; found {
						address.State = v.(string)
					}
					if v, found := msi["country"]; found {
						address.Country = v.(string)
					}
					if v, found := msi["zip"]; found {
						address.Zip = v.(string)
					}
					result = append(result, address)
				case "bankAccount":
					bankAccount := core.BankAccount{}
					if v, found := msi["account_type"]; found {
						bankAccount.AccountType = v.(string)
					}
					if v, found := msi["routing_number"]; found {
						bankAccount.RoutingNumber = v.(string)
					}
					if v, found := msi["account_number"]; found {
						bankAccount.AccountNumber = v.(string)
					}
					if v, found := msi["other_type"]; found {
						bankAccount.OtherType = v.(string)
					}
					result = append(result, bankAccount)
				case "host":
					host := core.Host{}
					if v, found := msi["host_name"]; found {
						host.Hostname = v.(string)
					}
					if v, found := msi["port"]; found {
						host.Port = v.(string)
					}
					result = append(result, host)
				case "keyPair":
					keyPair := core.KeyPair{}
					if v, found := msi["public_key"]; found {
						keyPair.PublicKey = v.(string)
					}
					if v, found := msi["private_key"]; found {
						keyPair.PrivateKey = v.(string)
					}
					result = append(result, keyPair)
				case "name":
					name := core.Name{}
					if v, found := msi["first"]; found {
						name.First = v.(string)
					}
					if v, found := msi["middle"]; found {
						name.Middle = v.(string)
					}
					if v, found := msi["last"]; found {
						name.Last = v.(string)
					}
					result = append(result, name)
				case "paymentCard":
					paymentCard := core.PaymentCard{}
					if v, found := msi["card_number"]; found {
						paymentCard.CardNumber = v.(string)
					}
					if v, found := msi["card_expiration_date"]; found {
						paymentCard.CardExpirationDate = v.(string)
					}
					if v, found := msi["card_security_code"]; found {
						paymentCard.CardSecurityCode = v.(string)
					}
					result = append(result, paymentCard)
				case "phone":
					phone := core.Phone{}
					if v, found := msi["region"]; found {
						phone.Region = v.(string)
					}
					if v, found := msi["number"]; found {
						phone.Number = v.(string)
					}
					if v, found := msi["ext"]; found {
						phone.Ext = v.(string)
					}
					if v, found := msi["type"]; found {
						phone.Type = v.(string)
					}
					result = append(result, phone)
				case "securityQuestion":
					securityQuestion := core.SecurityQuestion{}
					if v, found := msi["question"]; found {
						securityQuestion.Question = v.(string)
					}
					if v, found := msi["answer"]; found {
						securityQuestion.Answer = v.(string)
					}
					result = append(result, securityQuestion)
				}
			}
		}
	}

	return result
}

func convertFieldJsonToFieldSchema(field *genericFieldJson) *genericFieldSchema {
	// Complexity is TypeList with MaxItems=1 in schema but in record it is just a single object
	if field == nil {
		return nil
	}

	// convert Complexity to TypeList
	var complexity []*core.PasswordComplexity = nil
	if _, found := field.fields["complexity"]; found && field.Complexity != nil {
		complexity = []*core.PasswordComplexity{field.Complexity}
	}

	return &genericFieldSchema{
		fields:            field.fields,
		Type:              field.Type,
		Label:             field.Label,
		Required:          field.Required,
		PrivacyScreen:     field.PrivacyScreen,
		EnforceGeneration: field.EnforceGeneration,
		Complexity:        complexity,
		Value:             field.Value,
		ValueString:       field.ValueString,
	}
}

func parseFieldFromDataJson(data interface{}) *genericFieldSchema {
	if data != nil {
		if jsonData, err := json.Marshal(data); err == nil {
			fieldSchema := genericFieldJson{}
			if err = json.Unmarshal(jsonData, &fieldSchema); err == nil {
				fmap := map[string]interface{}{}
				if err = json.Unmarshal(jsonData, &fmap); err == nil {
					fieldSchema.fields = make(map[string]struct{}, len(fmap))
					for k := range fmap {
						fieldSchema.fields[k] = struct{}{}
					}
				}
				return convertFieldJsonToFieldSchema(&fieldSchema)
			}
		}
	}
	return nil
}

func parseFieldFromResourceDataJson(data interface{}) []*genericFieldSchema {
	result := []*genericFieldSchema{}
	if items, ok := data.([]interface{}); ok {
		for _, imap := range items {
			if jsonData, err := json.Marshal(imap); err == nil {
				fieldSchema := genericFieldSchema{}
				if err = json.Unmarshal(jsonData, &fieldSchema); err == nil {
					fmap := map[string]interface{}{}
					if err = json.Unmarshal(jsonData, &fmap); err == nil {
						fieldSchema.fields = make(map[string]struct{}, len(fmap))
						for k := range fmap {
							fieldSchema.fields[k] = struct{}{}
						}
					}
					result = append(result, &fieldSchema)
				}
			}
		}
	}
	return result
}

// NewFieldFromSchema parses simple fields into their corresponding type
func NewFieldFromSchema(fieldType string, fieldData interface{}) (newField interface{}, err error) {
	if data := parseFieldFromResourceDataJson(fieldData); len(data) > 0 {
		switch fieldType {
		case "accountNumber":
			field := &core.AccountNumber{KeeperRecordField: core.KeeperRecordField{Type: "accountNumber"}}
			if _, found := data[0].fields["label"]; found {
				field.Label = data[0].Label
			}
			if _, found := data[0].fields["required"]; found {
				field.Required = data[0].Required
			}
			if _, found := data[0].fields["privacy_screen"]; found {
				field.PrivacyScreen = data[0].PrivacyScreen
			}
			if _, found := data[0].fields["value"]; found {
				if v, ok := GetGenericFieldSchemaValueString(data[0]); ok && v != "" {
					field.Value = append(field.Value, v)
				}
			}
			return field, nil
		case "address":
			field := &core.Addresses{KeeperRecordField: core.KeeperRecordField{Type: "address"}}
			if _, found := data[0].fields["label"]; found {
				field.Label = data[0].Label
			}
			if _, found := data[0].fields["required"]; found {
				field.Required = data[0].Required
			}
			if _, found := data[0].fields["privacy_screen"]; found {
				field.PrivacyScreen = data[0].PrivacyScreen
			}
			if _, found := data[0].fields["value"]; found {
				if v := GetGenericFieldSchemaValue(field.Type, data[0]); len(v) > 0 {
					field.Value = append(field.Value, v[0].(core.Address))
				}
			}
			return field, nil
		case "addressRef":
			field := &core.AddressRef{KeeperRecordField: core.KeeperRecordField{Type: "addressRef"}}
			if _, found := data[0].fields["label"]; found {
				field.Label = data[0].Label
			}
			if _, found := data[0].fields["required"]; found {
				field.Required = data[0].Required
			}
			if _, found := data[0].fields["privacy_screen"]; found {
				field.PrivacyScreen = data[0].PrivacyScreen
			}
			if _, found := data[0].fields["value"]; found {
				if v, ok := GetGenericFieldSchemaValueString(data[0]); ok && v != "" {
					field.Value = append(field.Value, v)
				}
			}
			return field, nil
		case "bankAccount":
			field := &core.BankAccounts{KeeperRecordField: core.KeeperRecordField{Type: "bankAccount"}}
			if _, found := data[0].fields["label"]; found {
				field.Label = data[0].Label
			}
			if _, found := data[0].fields["required"]; found {
				field.Required = data[0].Required
			}
			if _, found := data[0].fields["privacy_screen"]; found {
				field.PrivacyScreen = data[0].PrivacyScreen
			}
			if _, found := data[0].fields["value"]; found {
				if v := GetGenericFieldSchemaValue(field.Type, data[0]); len(v) > 0 {
					field.Value = append(field.Value, v[0].(core.BankAccount))
				}
			}
			return field, nil
		case "birthDate":
			field := &core.BirthDate{KeeperRecordField: core.KeeperRecordField{Type: "birthDate"}}
			if _, found := data[0].fields["label"]; found {
				field.Label = data[0].Label
			}
			if _, found := data[0].fields["required"]; found {
				field.Required = data[0].Required
			}
			if _, found := data[0].fields["privacy_screen"]; found {
				field.PrivacyScreen = data[0].PrivacyScreen
			}
			if _, found := data[0].fields["value"]; found {
				if v, ok := GetGenericFieldSchemaValueInt64(data[0]); ok {
					field.Value = append(field.Value, v)
				}
			}
			return field, nil
		case "cardRef":
			field := &core.CardRef{KeeperRecordField: core.KeeperRecordField{Type: "cardRef"}}
			if _, found := data[0].fields["label"]; found {
				field.Label = data[0].Label
			}
			if _, found := data[0].fields["required"]; found {
				field.Required = data[0].Required
			}
			if _, found := data[0].fields["privacy_screen"]; found {
				field.PrivacyScreen = data[0].PrivacyScreen
			}
			if _, found := data[0].fields["value"]; found {
				if v, ok := GetGenericFieldSchemaValueString(data[0]); ok && v != "" {
					field.Value = append(field.Value, v)
				}
			}
			return field, nil
		case "date":
			field := &core.Date{KeeperRecordField: core.KeeperRecordField{Type: "date"}}
			if _, found := data[0].fields["label"]; found {
				field.Label = data[0].Label
			}
			if _, found := data[0].fields["required"]; found {
				field.Required = data[0].Required
			}
			if _, found := data[0].fields["privacy_screen"]; found {
				field.PrivacyScreen = data[0].PrivacyScreen
			}
			if _, found := data[0].fields["value"]; found {
				if v, ok := GetGenericFieldSchemaValueInt64(data[0]); ok {
					field.Value = append(field.Value, v)
				}
			}
			return field, nil
		case "email":
			field := &core.Email{KeeperRecordField: core.KeeperRecordField{Type: "email"}}
			if _, found := data[0].fields["label"]; found {
				field.Label = data[0].Label
			}
			if _, found := data[0].fields["required"]; found {
				field.Required = data[0].Required
			}
			if _, found := data[0].fields["privacy_screen"]; found {
				field.PrivacyScreen = data[0].PrivacyScreen
			}
			if _, found := data[0].fields["value"]; found {
				if v, ok := GetGenericFieldSchemaValueString(data[0]); ok && v != "" {
					field.Value = append(field.Value, v)
				}
			}
			return field, nil
		case "expirationDate":
			field := &core.ExpirationDate{KeeperRecordField: core.KeeperRecordField{Type: "expirationDate"}}
			if _, found := data[0].fields["label"]; found {
				field.Label = data[0].Label
			}
			if _, found := data[0].fields["required"]; found {
				field.Required = data[0].Required
			}
			if _, found := data[0].fields["privacy_screen"]; found {
				field.PrivacyScreen = data[0].PrivacyScreen
			}
			if _, found := data[0].fields["value"]; found {
				if v, ok := GetGenericFieldSchemaValueInt64(data[0]); ok {
					field.Value = append(field.Value, v)
				}
			}
			return field, nil
		case "fileRef":
			field := &core.FileRef{KeeperRecordField: core.KeeperRecordField{Type: "fileRef"}}
			if _, found := data[0].fields["label"]; found {
				field.Label = data[0].Label
			}
			if _, found := data[0].fields["required"]; found {
				field.Required = data[0].Required
			}
			if _, found := data[0].fields["value"]; found {
				if v, ok := GetGenericFieldSchemaValueStringList(data[0]); ok && len(v) > 0 {
					field.Value = append(field.Value, v...)
				}
			}
			return field, nil
		case "host":
			field := &core.Hosts{KeeperRecordField: core.KeeperRecordField{Type: "host"}}
			if _, found := data[0].fields["label"]; found {
				field.Label = data[0].Label
			}
			if _, found := data[0].fields["required"]; found {
				field.Required = data[0].Required
			}
			if _, found := data[0].fields["privacy_screen"]; found {
				field.PrivacyScreen = data[0].PrivacyScreen
			}
			if _, found := data[0].fields["value"]; found {
				if v := GetGenericFieldSchemaValue(field.Type, data[0]); len(v) > 0 {
					field.Value = append(field.Value, v[0].(core.Host))
				}
			}
			return field, nil
		case "keyPair":
			field := &core.KeyPairs{KeeperRecordField: core.KeeperRecordField{Type: "keyPair"}}
			if _, found := data[0].fields["label"]; found {
				field.Label = data[0].Label
			}
			if _, found := data[0].fields["required"]; found {
				field.Required = data[0].Required
			}
			if _, found := data[0].fields["privacy_screen"]; found {
				field.PrivacyScreen = data[0].PrivacyScreen
			}
			if _, found := data[0].fields["value"]; found {
				if v := GetGenericFieldSchemaValue(field.Type, data[0]); len(v) > 0 {
					field.Value = append(field.Value, v[0].(core.KeyPair))
				}
			}
			return field, nil
		case "licenseNumber":
			field := &core.LicenseNumber{KeeperRecordField: core.KeeperRecordField{Type: "licenseNumber"}}
			if _, found := data[0].fields["label"]; found {
				field.Label = data[0].Label
			}
			if _, found := data[0].fields["required"]; found {
				field.Required = data[0].Required
			}
			if _, found := data[0].fields["privacy_screen"]; found {
				field.PrivacyScreen = data[0].PrivacyScreen
			}
			if _, found := data[0].fields["value"]; found {
				if v, ok := GetGenericFieldSchemaValueString(data[0]); ok && v != "" {
					field.Value = append(field.Value, v)
				}
			}
			return field, nil
		case "login":
			field := &core.Login{KeeperRecordField: core.KeeperRecordField{Type: "login"}}
			if _, found := data[0].fields["label"]; found {
				field.Label = data[0].Label
			}
			if _, found := data[0].fields["required"]; found {
				field.Required = data[0].Required
			}
			if _, found := data[0].fields["privacy_screen"]; found {
				field.PrivacyScreen = data[0].PrivacyScreen
			}
			if _, found := data[0].fields["value"]; found {
				if v, ok := GetGenericFieldSchemaValueString(data[0]); ok && v != "" {
					field.Value = append(field.Value, v)
				}
			}
			return field, nil
		case "multiline":
			field := &core.Multiline{KeeperRecordField: core.KeeperRecordField{Type: "multiline"}}
			if _, found := data[0].fields["label"]; found {
				field.Label = data[0].Label
			}
			if _, found := data[0].fields["required"]; found {
				field.Required = data[0].Required
			}
			if _, found := data[0].fields["privacy_screen"]; found {
				field.PrivacyScreen = data[0].PrivacyScreen
			}
			if _, found := data[0].fields["value"]; found {
				if v, ok := GetGenericFieldSchemaValueString(data[0]); ok && v != "" {
					field.Value = append(field.Value, v)
				}
			}
			return field, nil
		case "name":
			field := &core.Names{KeeperRecordField: core.KeeperRecordField{Type: "name"}}
			if _, found := data[0].fields["label"]; found {
				field.Label = data[0].Label
			}
			if _, found := data[0].fields["required"]; found {
				field.Required = data[0].Required
			}
			if _, found := data[0].fields["privacy_screen"]; found {
				field.PrivacyScreen = data[0].PrivacyScreen
			}
			if _, found := data[0].fields["value"]; found {
				if v := GetGenericFieldSchemaValue(field.Type, data[0]); len(v) > 0 {
					field.Value = append(field.Value, v[0].(core.Name))
				}
			}
			return field, nil
		case "note":
			field := &core.SecureNote{KeeperRecordField: core.KeeperRecordField{Type: "note"}}
			if _, found := data[0].fields["label"]; found {
				field.Label = data[0].Label
			}
			if _, found := data[0].fields["required"]; found {
				field.Required = data[0].Required
			}
			if _, found := data[0].fields["privacy_screen"]; found {
				field.PrivacyScreen = data[0].PrivacyScreen
			}
			if _, found := data[0].fields["value"]; found {
				if v, ok := GetGenericFieldSchemaValueString(data[0]); ok && v != "" {
					field.Value = append(field.Value, v)
				}
			}
			return field, nil
		case "oneTimeCode":
			field := &core.OneTimeCode{KeeperRecordField: core.KeeperRecordField{Type: "oneTimeCode"}}
			if _, found := data[0].fields["label"]; found {
				field.Label = data[0].Label
			}
			if _, found := data[0].fields["required"]; found {
				field.Required = data[0].Required
			}
			if _, found := data[0].fields["privacy_screen"]; found {
				field.PrivacyScreen = data[0].PrivacyScreen
			}
			if _, found := data[0].fields["value"]; found {
				if v, ok := GetGenericFieldSchemaValueString(data[0]); ok && v != "" {
					field.Value = append(field.Value, v)
				}
			}
			return field, nil
		case "password":
			field := &core.Password{KeeperRecordField: core.KeeperRecordField{Type: "password"}}
			if _, found := data[0].fields["label"]; found {
				field.Label = data[0].Label
			}
			if _, found := data[0].fields["required"]; found {
				field.Required = data[0].Required
			}
			if _, found := data[0].fields["privacy_screen"]; found {
				field.PrivacyScreen = data[0].PrivacyScreen
			}
			if _, found := data[0].fields["enforce_generation"]; found {
				field.EnforceGeneration = data[0].EnforceGeneration
			}
			if _, found := data[0].fields["complexity"]; found && len(data[0].Complexity) > 0 {
				field.Complexity = data[0].Complexity[0]
			}
			if _, found := data[0].fields["value"]; found {
				if v, ok := GetGenericFieldSchemaValueString(data[0]); ok && v != "" {
					field.Value = append(field.Value, v)
				}
			}
			return field, nil
		case "paymentCard":
			field := &core.PaymentCards{KeeperRecordField: core.KeeperRecordField{Type: "paymentCard"}}
			if _, found := data[0].fields["label"]; found {
				field.Label = data[0].Label
			}
			if _, found := data[0].fields["required"]; found {
				field.Required = data[0].Required
			}
			if _, found := data[0].fields["privacy_screen"]; found {
				field.PrivacyScreen = data[0].PrivacyScreen
			}
			if _, found := data[0].fields["value"]; found {
				if v := GetGenericFieldSchemaValue(field.Type, data[0]); len(v) > 0 {
					field.Value = append(field.Value, v[0].(core.PaymentCard))
				}
			}
			return field, nil
		case "phone":
			field := &core.Phones{KeeperRecordField: core.KeeperRecordField{Type: "phone"}}
			if _, found := data[0].fields["label"]; found {
				field.Label = data[0].Label
			}
			if _, found := data[0].fields["required"]; found {
				field.Required = data[0].Required
			}
			if _, found := data[0].fields["privacy_screen"]; found {
				field.PrivacyScreen = data[0].PrivacyScreen
			}
			if _, found := data[0].fields["value"]; found {
				if v := GetGenericFieldSchemaValue(field.Type, data[0]); len(v) > 0 {
					field.Value = append(field.Value, v[0].(core.Phone))
				}
			}
			return field, nil
		case "pinCode":
			field := &core.PinCode{KeeperRecordField: core.KeeperRecordField{Type: "pinCode"}}
			if _, found := data[0].fields["label"]; found {
				field.Label = data[0].Label
			}
			if _, found := data[0].fields["required"]; found {
				field.Required = data[0].Required
			}
			if _, found := data[0].fields["privacy_screen"]; found {
				field.PrivacyScreen = data[0].PrivacyScreen
			}
			if _, found := data[0].fields["value"]; found {
				if v, ok := GetGenericFieldSchemaValueString(data[0]); ok && v != "" {
					field.Value = append(field.Value, v)
				}
			}
			return field, nil
		case "secret":
			field := &core.Secret{KeeperRecordField: core.KeeperRecordField{Type: "secret"}}
			if _, found := data[0].fields["label"]; found {
				field.Label = data[0].Label
			}
			if _, found := data[0].fields["required"]; found {
				field.Required = data[0].Required
			}
			if _, found := data[0].fields["privacy_screen"]; found {
				field.PrivacyScreen = data[0].PrivacyScreen
			}
			if _, found := data[0].fields["value"]; found {
				if v, ok := GetGenericFieldSchemaValueString(data[0]); ok && v != "" {
					field.Value = append(field.Value, v)
				}
			}
			return field, nil
		case "securityQuestion":
			field := &core.SecurityQuestions{KeeperRecordField: core.KeeperRecordField{Type: "securityQuestion"}}
			if _, found := data[0].fields["label"]; found {
				field.Label = data[0].Label
			}
			if _, found := data[0].fields["required"]; found {
				field.Required = data[0].Required
			}
			if _, found := data[0].fields["privacy_screen"]; found {
				field.PrivacyScreen = data[0].PrivacyScreen
			}
			if _, found := data[0].fields["value"]; found {
				if v := GetGenericFieldSchemaValue(field.Type, data[0]); len(v) > 0 {
					field.Value = append(field.Value, v[0].(core.SecurityQuestion))
				}
			}
			return field, nil
		case "text":
			field := &core.Text{KeeperRecordField: core.KeeperRecordField{Type: "text"}}
			if _, found := data[0].fields["label"]; found {
				field.Label = data[0].Label
			}
			if _, found := data[0].fields["required"]; found {
				field.Required = data[0].Required
			}
			if _, found := data[0].fields["privacy_screen"]; found {
				field.PrivacyScreen = data[0].PrivacyScreen
			}
			if _, found := data[0].fields["value"]; found {
				if v, ok := GetGenericFieldSchemaValueString(data[0]); ok && v != "" {
					field.Value = append(field.Value, v)
				}
			}
			return field, nil
		case "url":
			field := &core.Url{KeeperRecordField: core.KeeperRecordField{Type: "url"}}
			if _, found := data[0].fields["label"]; found {
				field.Label = data[0].Label
			}
			if _, found := data[0].fields["required"]; found {
				field.Required = data[0].Required
			}
			if _, found := data[0].fields["privacy_screen"]; found {
				field.PrivacyScreen = data[0].PrivacyScreen
			}
			if _, found := data[0].fields["value"]; found {
				if v, ok := GetGenericFieldSchemaValueString(data[0]); ok && v != "" {
					field.Value = append(field.Value, v)
				}
			}
			return field, nil
		default:
			return nil, fmt.Errorf("unable to create unknown field type %v", fieldType)
		}
	}
	return nil, nil
}

func NewFieldFromType(fieldType string) (field interface{}, err error) {
	switch fieldType {
	case "login":
		return &core.Login{}, nil
	case "password":
		return &core.Password{}, nil
	case "url":
		return &core.Url{}, nil
	case "fileRef":
		return &core.FileRef{}, nil
	case "oneTimeCode":
		return &core.OneTimeCode{}, nil
	case "name":
		return &core.Names{}, nil
	case "birthDate":
		return &core.BirthDate{}, nil
	case "date":
		return &core.Date{}, nil
	case "expirationDate":
		return &core.ExpirationDate{}, nil
	case "text":
		return &core.Text{}, nil
	case "securityQuestion":
		return &core.SecurityQuestions{}, nil
	case "multiline":
		return &core.Multiline{}, nil
	case "email":
		return &core.Email{}, nil
	case "cardRef":
		return &core.CardRef{}, nil
	case "addressRef":
		return &core.AddressRef{}, nil
	case "pinCode":
		return &core.PinCode{}, nil
	case "phone":
		return &core.Phones{}, nil
	case "secret":
		return &core.Secret{}, nil
	case "note":
		return &core.SecureNote{}, nil
	case "accountNumber":
		return &core.AccountNumber{}, nil
	case "paymentCard":
		return &core.PaymentCards{}, nil
	case "bankAccount":
		return &core.BankAccounts{}, nil
	case "keyPair":
		return &core.KeyPairs{}, nil
	case "host":
		return &core.Hosts{}, nil
	case "address":
		return &core.Addresses{}, nil
	case "licenseNumber":
		return &core.LicenseNumber{}, nil
	default:
		return nil, fmt.Errorf("unable to create unknown field type %v", fieldType)
	}
}

func ParseGeneratePassword(data interface{}) (bool, error) {
	if s, ok := data.([]interface{}); ok && len(s) > 0 {
		if m, ok := s[0].(map[string]interface{}); ok {
			if igen, ok := m["generate"]; ok {
				if sgen, ok := igen.(string); ok {
					if sgen == "" {
						return false, nil
					} else if sgen == "true" || sgen == "yes" {
						return true, nil
					} else {
						return false, fmt.Errorf("generate = %s - expected one of ('true', 'yes' or '')", sgen)
					}
				} else {
					return false, errors.New("generate should be a string ('true', 'yes' or '')")
				}
			}
		}
	}
	// generate = false when not present in schema
	return false, nil
}

func SetFieldTypeInSchema(d *schema.ResourceData, fieldName, fieldType string) error {
	if fieldData := d.Get(fieldName); fieldData != nil {
		if s, ok := fieldData.([]interface{}); ok && len(s) > 0 {
			if m, ok := s[0].(map[string]interface{}); ok {
				if itype, ok := m["type"]; ok {
					if _, ok := itype.(string); ok {
						// if stype != "" && stype != fieldType { return false, errors.New("field type is already set incorrectly") }
						m["type"] = fieldType
						return d.Set(fieldName, fieldData)
					}
				}
			}
		}
	}
	return errors.New("failed to set field type to " + fieldType)
}

func validateUid(uid string) bool {
	if ruid := strings.TrimSpace(uid); ruid != "" {
		base64UrlSafeRegexp := `^(?:[A-Za-z\d\-_]{4})*(?:[A-Za-z\d\-_]{3}=?|[A-Za-z\d\-_]{2}(?:==)?)?$`
		if matched, err := regexp.MatchString(base64UrlSafeRegexp, ruid); err == nil && matched {
			if numBytes := len(core.Base64ToBytes(ruid)); numBytes == 16 {
				return true
			}
		}
	}

	return false
}

func validateComplexity(length, uppercase, lowercase, digits, special int) error {
	sumWeights := uppercase + lowercase + digits + special
	if length >= 8 && length <= 100 {
		if uppercase >= 0 && uppercase <= length &&
			lowercase >= 0 && lowercase <= length &&
			digits >= 0 && digits <= length &&
			special >= 0 && special <= length &&
			sumWeights <= length {
			return nil
		}
	}

	return fmt.Errorf("expected - length in [8..100],"+
		" charset_len in [0..length], and the sum of all lengths in [0..length],"+
		" got length: %v, sum: %v = %v + %v + %v + %v",
		length, sumWeights, uppercase, lowercase, digits, special)
}

/*
func validateDiagFunc(fn schema.SchemaValidateFunc) schema.SchemaValidateDiagFunc {
	return func(v interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		warnings, errors := fn(v, fmt.Sprintf("%#v", path))
		for _, warning := range warnings {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Warning,
				Summary:       warning,
				Detail:        warning,
				AttributePath: path,
			})
		}
		for _, err := range errors {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       err.Error(),
				Detail:        err.Error(),
				AttributePath: path,
			})
		}
		return diags
	}
}
*/

func getSchemaAndRecordFieldNames(fieldName string) (schemaName string, recordName string) {
	if recordFieldName, found := mapSchemaToRecordFieldName[fieldName]; found {
		// fieldName is schemaFieldName
		return fieldName, recordFieldName
	} else {
		// check if fieldName is recordFieldName
		for schemaFieldName, recordFieldName := range mapSchemaToRecordFieldName {
			if fieldName == recordFieldName {
				return schemaFieldName, recordFieldName
			}
		}
	}

	return "", ""
}

func ApplyFieldChange(section, name string, d *schema.ResourceData, record *core.Record) (int, error) {
	modified := 0

	if d == nil || record == nil {
		return modified, fmt.Errorf("apply change expects both schema and record to be non empty")
	}
	if section != "fields" && section != "custom" {
		return modified, fmt.Errorf("apply change expects field section to be one of ['fields', 'custom'] - got '%s'", section)
	}
	if name = strings.TrimSpace(name); name == "" {
		return modified, fmt.Errorf("apply change expects field name to be non empty")
	}

	schemaFieldName, recordFieldName := getSchemaAndRecordFieldNames(name)
	if schemaFieldName == "" {
		return modified, fmt.Errorf("apply change was unable to find schema field name for field '%s'", name)
	}
	if recordFieldName == "" {
		return modified, fmt.Errorf("apply change was unable to find record field name for field '%s'", name)
	}

	if fieldData, exists := d.GetOk(schemaFieldName); exists {
		// field present in configuration
		if fieldData != nil && len(fieldData.([]interface{})) > 0 {
			// check if password field needs to be re/generated
			generate := false
			if recordFieldName == "password" {
				if oldf, newf := d.GetChange("password"); newf != nil && len(newf.([]interface{})) > 0 {
					oldg, newg := "", ""
					if fmap, ok := newf.([]interface{})[0].(map[string]interface{}); ok {
						if og, found := fmap["generate"]; found {
							newg = og.(string)
						}
					}
					if oldf != nil && len(oldf.([]interface{})) > 0 {
						if fmap, ok := oldf.([]interface{})[0].(map[string]interface{}); ok {
							if og, found := fmap["generate"]; found {
								oldg = og.(string)
							}
						}
					}
					generate = newg != "" && newg != oldg
				}
			}
			if field, err := NewFieldFromSchema(recordFieldName, fieldData); err != nil {
				return modified, err
			} else if field == nil {
				return modified, fmt.Errorf("apply change failed to convert schema '%s' to field '%s' from field data: '%v'", schemaFieldName, recordFieldName, fieldData)
			} else {
				if generate {
					if generated, err := applyGeneratePassword(fieldData, field); err != nil {
						return modified, err
					} else if generated {
						if err := d.Set("password", fieldData); err != nil {
							return modified, err
						}
					}
				}
				// upsert
				if record.FieldExists(section, recordFieldName) {
					if err := record.UpdateField(section, field); err != nil {
						return modified, err
					}
					modified++
				} else {
					if err := record.InsertField(section, field); err != nil {
						return modified, err
					}
					modified++
				}
			}
		} else {
			return modified, fmt.Errorf("apply change failed to get field data from schema - field: '%s', data: '%v'", schemaFieldName, fieldData)
		}
	} else {
		// field not present in configuration - remove from record if present
		modified = record.RemoveField("fields", recordFieldName, false)
	}

	return modified, nil
}

func mergePassword(schemaField interface{}, recordField interface{}) {
	// password field must merge with schema to pull data not stored in record like generate=true
	// merge schema only attributes back into the new value before schema update
	if schemaField != nil && recordField != nil {
		var generate interface{} = nil
		if sfi, ok := schemaField.([]interface{}); ok && len(sfi) > 0 {
			if sfmap, ok := sfi[0].(map[string]interface{}); ok {
				if sfg, found := sfmap["generate"]; found {
					generate = sfg
				}
			}
		}
		if generate != nil {
			if sfi, ok := recordField.([]interface{}); ok && len(sfi) > 0 {
				if sfmap, ok := sfi[0].(map[string]interface{}); ok {
					sfmap["generate"] = generate
				}
			}
		}
	}
}

func applyGeneratePassword(fieldData interface{}, field interface{}) (generated bool, e error) {
	if fv, ok := field.(*core.Password); ok {
		complexity := core.PasswordComplexity{Length: 16}
		if fv.Complexity != nil {
			if err := validateComplexity(fv.Complexity.Length,
				fv.Complexity.Caps,
				fv.Complexity.Lowercase,
				fv.Complexity.Digits,
				fv.Complexity.Special); err != nil {
				return false, err
			}
			complexity = *fv.Complexity
		}
		if generate, _ := ParseGeneratePassword(fieldData); generate {
			if pwd, err := core.GeneratePassword(complexity.Length,
				complexity.Lowercase,
				complexity.Caps,
				complexity.Digits,
				complexity.Special); err != nil {
				return false, err
			} else {
				if len(fv.Value) > 0 {
					fv.Value = fv.Value[0:0]
				}
				// update secret
				fv.Value = append(fv.Value, pwd)
				// update schema
				if fmap, ok := fieldData.([]interface{})[0].(map[string]interface{}); ok {
					fmap["value"] = pwd
				}
				return true, nil
			}
		}
	} else {
		return false, fmt.Errorf("applyGeneratePassword expects field to be of type *core.Password")
	}
	return false, nil
}
