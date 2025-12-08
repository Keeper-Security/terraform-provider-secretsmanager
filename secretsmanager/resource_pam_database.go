package secretsmanager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/keeper-security/secrets-manager-go/core"
)

func resourcePamDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePamDatabaseCreate,
		ReadContext:   resourcePamDatabaseRead,
		UpdateContext: resourcePamDatabaseUpdate,
		DeleteContext: resourcePamDatabaseDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePamDatabaseImport,
		},
		Schema: map[string]*schema.Schema{
			"folder_uid": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				AtLeastOneOf: []string{"folder_uid", "uid"},
				Description:  "The folder UID where the secret is stored. Ensure the folder is shared to your KSM application with 'Can Edit' permissions.",
			},
			"uid": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"folder_uid", "uid"},
				Description:  "The UID of the new secret (using RFC4648 URL and Filename Safe Alphabet).",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The secret type.",
			},
			"title": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The secret title.",
			},
			"notes": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The secret notes.",
			},
			// PAM Database specific fields
			"pam_hostname":     schemaPamHostnameField(),
			"pam_settings":     schemaPamSettingsField(),
			"use_ssl":          schemaCheckboxField(),
			"login":            schemaLoginField(),
			"password":         schemaPasswordField(""),
			"rotation_scripts": schemaScriptField(),
			"connect_database": schemaTextField(),
			"database_id":      schemaTextField(),
			"database_type":    schemaDatabaseTypeField(),
			"provider_group":   schemaTextField(),
			"provider_region":  schemaTextField(),
			"file_ref":         schemaFileRefField(),
			"custom": schemaCustomField(),
			"totp":             schemaOneTimeCodeField(),
		},
	}
}

func resourcePamDatabaseCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics
	// Validate custom field labels are unique
	if validateDiags := validateUniqueCustomFieldLabels(d); len(validateDiags) > 0 {
		return validateDiags
	}


	uid := strings.TrimSpace(d.Get("uid").(string))
	if uid == "" {
		uid = core.GenerateUid()
	}
	if validUid := validateUid(uid); !validUid {
		return diag.Errorf("invalid UID format - use unpadded base64url encoded value (RFC 4648)")
	}

	folderUid := strings.TrimSpace(d.Get("folder_uid").(string))
	if folderUid == "" {
		return diag.Errorf("'folder_uid' is required to create new resource")
	}

	nrc := core.NewRecordCreate("pamDatabase", "")
	if title := d.Get("title"); title != nil && title.(string) != "" {
		nrc.Title = title.(string)
	}
	if notes := d.Get("notes"); notes != nil && notes.(string) != "" {
		nrc.Notes = notes.(string)
	}

	if fieldData := d.Get("pam_hostname"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("pamHostname", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			nrc.Fields = append(nrc.Fields, field)
			if err := SetFieldTypeInSchema(d, "pam_hostname", "pamHostname"); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	// Handle pam_settings as JSON string
	// Note: pam_settings is a TypeString, not TypeList, so we don't call SetFieldTypeInSchema
	if pamSettingsJSON := d.Get("pam_settings").(string); pamSettingsJSON != "" {
		if field, err := createPamSettingsFieldFromJSON(pamSettingsJSON); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			nrc.Fields = append(nrc.Fields, field)
		}
	}
	if fieldData := d.Get("use_ssl"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("checkbox", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			field.(*core.Checkbox).Label = "useSSL"
			nrc.Fields = append(nrc.Fields, field)
		}
	}
	if fieldData := d.Get("login"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("login", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			nrc.Fields = append(nrc.Fields, field)
			if err := SetFieldTypeInSchema(d, "login", "login"); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if fieldData := d.Get("password"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("password", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			if generated, err := applyGeneratePassword(fieldData, field); err != nil {
				return diag.FromErr(err)
			} else if generated {
				if err := d.Set("password", fieldData); err != nil {
					return diag.FromErr(err)
				}
			}
			nrc.Fields = append(nrc.Fields, field)
			if err := SetFieldTypeInSchema(d, "password", "password"); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if fieldData := d.Get("rotation_scripts"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("script", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			field.(*core.Scripts).Label = "Rotation Scripts"
			nrc.Fields = append(nrc.Fields, field)
			if err := SetFieldTypeInSchema(d, "rotation_scripts", "script"); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if fieldData := d.Get("connect_database"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("text", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			field.(*core.Text).Label = "Connect Database"
			nrc.Fields = append(nrc.Fields, field)
			if err := SetFieldTypeInSchema(d, "connect_database", "text"); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if fieldData := d.Get("database_id"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("text", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			field.(*core.Text).Label = "Database Id"
			nrc.Fields = append(nrc.Fields, field)
			if err := SetFieldTypeInSchema(d, "database_id", "text"); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if databaseType := d.Get("database_type").(string); databaseType != "" {
		field := core.NewDatabaseType(databaseType)
		nrc.Fields = append(nrc.Fields, field)
	}
	if fieldData := d.Get("provider_group"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("text", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			field.(*core.Text).Label = "Provider Group"
			nrc.Fields = append(nrc.Fields, field)
			if err := SetFieldTypeInSchema(d, "provider_group", "text"); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if fieldData := d.Get("provider_region"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("text", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			field.(*core.Text).Label = "Provider Region"
			nrc.Fields = append(nrc.Fields, field)
			if err := SetFieldTypeInSchema(d, "provider_region", "text"); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if fieldData := d.Get("totp"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("oneTimeCode", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			nrc.Fields = append(nrc.Fields, field)
			if err := SetFieldTypeInSchema(d, "totp", "oneTimeCode"); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if fieldData := d.Get("file_ref"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("fileRef", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			nrc.Fields = append(nrc.Fields, field)
			if err := SetFieldTypeInSchema(d, "file_ref", "fileRef"); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	// Process custom fields
	if customData := d.Get("custom"); customData != nil && len(customData.([]interface{})) > 0 {
		for _, customItem := range customData.([]interface{}) {
			if customMap, ok := customItem.(map[string]interface{}); ok {
				fieldType := "text" // default to text
				if ft, ok := customMap["type"].(string); ok && ft != "" {
					fieldType = ft
				}
				
				// Support text, multiline, secret, url, and email field types
				switch fieldType {
				case "text", "multiline", "secret", "url", "email":
					var field interface{}
					switch fieldType {
					case "text":
						field = &core.Text{KeeperRecordField: core.KeeperRecordField{Type: "text"}}
					case "multiline":
						field = &core.Multiline{KeeperRecordField: core.KeeperRecordField{Type: "multiline"}}
					case "secret":
						field = &core.Secret{KeeperRecordField: core.KeeperRecordField{Type: "secret"}}
					case "url":
						field = &core.Url{KeeperRecordField: core.KeeperRecordField{Type: "url"}}
					case "email":
						field = &core.Email{KeeperRecordField: core.KeeperRecordField{Type: "email"}}
					}

					// Set common properties using type assertion
					switch f := field.(type) {
					case *core.Text:
						if label, ok := customMap["label"].(string); ok { f.Label = label }
						if required, ok := customMap["required"].(bool); ok { f.Required = required }
						if privacyScreen, ok := customMap["privacy_screen"].(bool); ok { f.PrivacyScreen = privacyScreen }
						if value, ok := customMap["value"].(string); ok && value != "" { f.Value = []string{value} }
					case *core.Multiline:
						if label, ok := customMap["label"].(string); ok { f.Label = label }
						if required, ok := customMap["required"].(bool); ok { f.Required = required }
						if privacyScreen, ok := customMap["privacy_screen"].(bool); ok { f.PrivacyScreen = privacyScreen }
						if value, ok := customMap["value"].(string); ok && value != "" { f.Value = []string{value} }
					case *core.Secret:
						if label, ok := customMap["label"].(string); ok { f.Label = label }
						if required, ok := customMap["required"].(bool); ok { f.Required = required }
						if privacyScreen, ok := customMap["privacy_screen"].(bool); ok { f.PrivacyScreen = privacyScreen }
						if value, ok := customMap["value"].(string); ok && value != "" { f.Value = []string{value} }
					case *core.Url:
						if label, ok := customMap["label"].(string); ok { f.Label = label }
						if required, ok := customMap["required"].(bool); ok { f.Required = required }
						if privacyScreen, ok := customMap["privacy_screen"].(bool); ok { f.PrivacyScreen = privacyScreen }
						if value, ok := customMap["value"].(string); ok && value != "" { f.Value = []string{value} }
					case *core.Email:
						if label, ok := customMap["label"].(string); ok { f.Label = label }
						if required, ok := customMap["required"].(bool); ok { f.Required = required }
						if privacyScreen, ok := customMap["privacy_screen"].(bool); ok { f.PrivacyScreen = privacyScreen }
						if value, ok := customMap["value"].(string); ok && value != "" { f.Value = []string{value} }
					}

					nrc.Custom = append(nrc.Custom, field)
				}
			}
		}
	}


	if folderUid == "*" {
		if fuid, err := getTemplateFolder(folderUid, client); err != nil {
			return diag.FromErr(err)
		} else {
			folderUid = fuid
		}
	}

	uid, err := createRecord(uid, folderUid, nrc, client)
	if err != nil {
		return diag.FromErr(err)
	}

	if fuid := strings.TrimSpace(d.Get("folder_uid").(string)); fuid == "*" {
		if err = d.Set("folder_uid", folderUid); err != nil {
			return diag.FromErr(err)
		}
	}
	if err = d.Set("uid", uid); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("type", "pamDatabase"); err != nil {
		return diag.FromErr(err)
	}


	d.SetId(uid)
	return diags
}

func resourcePamDatabaseRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics

	uid := strings.TrimSpace(d.Get("uid").(string))
	title := strings.TrimSpace(d.Get("title").(string))
	if uid == "" && title == "" {
		return diag.Errorf("record UID and/or title required to locate the record")
	}

	secret, err := getRecord(uid, title, client)
	if err != nil {
		if strings.HasPrefix(err.Error(), "record not found") {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}

	resourceType := "pamDatabase"
	recordType := secret.Type()
	if recordType != resourceType {
		return diag.Errorf("record type '%s' is not the expected type '%s' for this resource", recordType, resourceType)
	}
	if uid == "" {
		if err = d.Set("uid", secret.Uid); err != nil {
			return diag.FromErr(err)
		}
	}
	fuid := secret.InnerFolderUid()
	if fuid == "" {
		fuid = secret.FolderUid()
	}
	if fuid != "" {
		if err = d.Set("folder_uid", fuid); err != nil {
			return diag.FromErr(err)
		}
	}
	if err = d.Set("type", recordType); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("title", secret.Title()); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("notes", secret.Notes()); err != nil {
		return diag.FromErr(err)
	}

	login := getFieldResourceData("login", "fields", secret)
	if err = d.Set("login", login); err != nil {
		return diag.FromErr(err)
	}
	password := getFieldResourceData("password", "fields", secret)
	mergePassword(d.Get("password"), password)
	if err = d.Set("password", password); err != nil {
		return diag.FromErr(err)
	}
	oneTimeCode := getFieldResourceData("oneTimeCode", "fields", secret)
	if err = d.Set("totp", oneTimeCode); err != nil {
		return diag.FromErr(err)
	}

	// PAM Database specific fields
	pamHostname := getFieldResourceData("pamHostname", "fields", secret)
	if err = d.Set("pam_hostname", pamHostname); err != nil {
		return diag.FromErr(err)
	}
	// Read pam_settings as JSON string
	if pamSettingsFields := secret.GetFieldsByType("pamSettings"); len(pamSettingsFields) > 0 {
		if pamSettingsJSON, err := pamSettingsFieldToJSON(pamSettingsFields[0]); err != nil {
			return diag.FromErr(err)
		} else if err = d.Set("pam_settings", pamSettingsJSON); err != nil {
			return diag.FromErr(err)
		}
	}
	useSSL := getFieldResourceDataWithLabel("checkbox", "fields", secret, "useSSL")
	if err = d.Set("use_ssl", useSSL); err != nil {
		return diag.FromErr(err)
	}
	rotationScripts := getFieldResourceDataWithLabel("script", "fields", secret, "Rotation Scripts")
	if err = d.Set("rotation_scripts", rotationScripts); err != nil {
		return diag.FromErr(err)
	}
	connectDatabase := getFieldResourceDataWithLabel("text", "fields", secret, "Connect Database")
	if err = d.Set("connect_database", connectDatabase); err != nil {
		return diag.FromErr(err)
	}
	databaseId := getFieldResourceDataWithLabel("text", "fields", secret, "Database Id")
	if err = d.Set("database_id", databaseId); err != nil {
		return diag.FromErr(err)
	}
	// Read database_type as a simple string value
	if databaseTypeFields := secret.GetFieldsByType("databaseType"); len(databaseTypeFields) > 0 {
		fieldMap := databaseTypeFields[0]
		if valueInterface, exists := fieldMap["value"]; exists {
			if valueList, ok := valueInterface.([]interface{}); ok && len(valueList) > 0 {
				if dbType, ok := valueList[0].(string); ok && dbType != "" {
					if err = d.Set("database_type", dbType); err != nil {
						return diag.FromErr(err)
					}
				}
			}
		}
	}
	providerGroup := getFieldResourceDataWithLabel("text", "fields", secret, "Provider Group")
	if err = d.Set("provider_group", providerGroup); err != nil {
		return diag.FromErr(err)
	}
	providerRegion := getFieldResourceDataWithLabel("text", "fields", secret, "Provider Region")
	if err = d.Set("provider_region", providerRegion); err != nil {
		return diag.FromErr(err)
	}

	fileItems := getFileItemsResourceData(secret)
	if err := d.Set("file_ref", fileItems); err != nil {
		return diag.FromErr(err)
	}


	// Read custom fields
	customItems := getFieldItemsResourceData("custom", secret)
	if err := d.Set("custom", customItems); err != nil {
		return diag.FromErr(err)
	}

	// Warn if duplicate custom field labels detected
	if warnDiags := warnDuplicateCustomFieldLabels(secret); len(warnDiags) > 0 {
		diags = append(diags, warnDiags...)
	}

	d.SetId(uid)
	return diags
}

func resourcePamDatabaseUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client

	// Validate custom field labels are unique
	if validateDiags := validateUniqueCustomFieldLabels(d); len(validateDiags) > 0 {
		return validateDiags
	}

	uid := strings.TrimSpace(d.Get("uid").(string))
	if uid == "" {
		return diag.Errorf("'uid' is required to update existing resource")
	}

	hasRestrictedChanges := d.HasChange("folder_uid") || d.HasChange("uid") || d.HasChange("type")
	if hasRestrictedChanges {
		return diag.Errorf("changes to folder_uid, uid, and type are not allowed")
	}

	title := strings.TrimSpace(d.Get("title").(string))
	secret, err := getRecord(uid, title, client)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("title") {
		secret.SetTitle(d.Get("title").(string))
	}
	if d.HasChange("notes") {
		secret.SetNotes(d.Get("notes").(string))
	}

	if d.HasChange("login") {
		if _, err := ApplyFieldChange("fields", "login", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("password") {
		if _, err := ApplyFieldChange("fields", "password", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("pam_hostname") {
		// Handle pam_hostname using SetStandardFieldValue which calls update() to sync RecordDict to RawJson
		pamHostnameData := d.Get("pam_hostname")
		pamHostnameList, ok := pamHostnameData.([]interface{})
		if !ok || len(pamHostnameList) == 0 {
			return diag.Errorf("pam_hostname is not a valid list")
		}
		pamHostnameMap, ok := pamHostnameList[0].(map[string]interface{})
		if !ok {
			return diag.Errorf("pam_hostname[0] is not a valid map")
		}
		valueList, ok := pamHostnameMap["value"].([]interface{})
		if !ok || len(valueList) == 0 {
			return diag.Errorf("pam_hostname value is not a valid list")
		}
		valueMap, ok := valueList[0].(map[string]interface{})
		if !ok {
			return diag.Errorf("pam_hostname value[0] is not a valid map")
		}

		// Construct the pamHostname value - SDK expects []interface{} with Host map
		hostValue := []interface{}{
			map[string]interface{}{
				"hostName": valueMap["hostname"],
				"port":     valueMap["port"],
			},
		}

		// Use SetStandardFieldValue which calls update() to sync RecordDict to RawJson
		if err := secret.SetStandardFieldValue("pamHostname", hostValue); err != nil {
			return diag.FromErr(fmt.Errorf("failed to update pam_hostname: %w", err))
		}
	}
	if d.HasChange("pam_settings") {
		// Handle pam_settings JSON string field - parse and use SetStandardFieldValue to sync to RawJson
		pamSettingsJSON := d.Get("pam_settings").(string)
		var pamSettingsValue interface{}
		if err := json.Unmarshal([]byte(pamSettingsJSON), &pamSettingsValue); err != nil {
			return diag.FromErr(fmt.Errorf("failed to parse pam_settings JSON: %w", err))
		}

		// Use SetStandardFieldValue which calls update() to sync RecordDict to RawJson
		if err := secret.SetStandardFieldValue("pamSettings", pamSettingsValue); err != nil {
			return diag.FromErr(fmt.Errorf("failed to update pam_settings: %w", err))
		}
	}
	if d.HasChange("use_ssl") {
		if _, err := ApplyFieldChange("fields", "use_ssl", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("rotation_scripts") {
		if _, err := ApplyFieldChange("fields", "rotation_scripts", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("connect_database") {
		if _, err := ApplyFieldChange("fields", "connect_database", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("database_id") {
		if _, err := ApplyFieldChange("fields", "database_id", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("database_type") {
		databaseType := d.Get("database_type").(string)
		if databaseType != "" {
			// Use SetStandardFieldValue which calls update() to sync RecordDict to RawJson
			if err := secret.SetStandardFieldValue("databaseType", []string{databaseType}); err != nil {
				return diag.FromErr(fmt.Errorf("failed to update database_type: %w", err))
			}
		}
	}
	if d.HasChange("provider_group") {
		if _, err := ApplyFieldChange("fields", "provider_group", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("provider_region") {
		if _, err := ApplyFieldChange("fields", "provider_region", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("totp") {
		if _, err := ApplyFieldChange("fields", "totp", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("file_ref") {
		if _, err := ApplyFieldChange("fields", "file_ref", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("custom") {
		// Clear existing custom fields
		customFields := []interface{}{}
		// Add updated custom fields
		if customData := d.Get("custom"); customData != nil && len(customData.([]interface{})) > 0 {
			for _, customItem := range customData.([]interface{}) {
				if customMap, ok := customItem.(map[string]interface{}); ok {
					fieldType := "text"
					if ft, ok := customMap["type"].(string); ok && ft != "" {
						fieldType = ft
					}
					switch fieldType {
					case "text", "multiline", "secret", "url", "email":
						var field interface{}
						switch fieldType {
						case "text":
							field = &core.Text{KeeperRecordField: core.KeeperRecordField{Type: "text"}}
						case "multiline":
							field = &core.Multiline{KeeperRecordField: core.KeeperRecordField{Type: "multiline"}}
						case "secret":
							field = &core.Secret{KeeperRecordField: core.KeeperRecordField{Type: "secret"}}
						case "url":
							field = &core.Url{KeeperRecordField: core.KeeperRecordField{Type: "url"}}
						case "email":
							field = &core.Email{KeeperRecordField: core.KeeperRecordField{Type: "email"}}
						}

						// Set common properties using type assertion
						var fieldMap map[string]interface{}
						switch f := field.(type) {
						case *core.Text:
							if label, ok := customMap["label"].(string); ok { f.Label = label }
							if required, ok := customMap["required"].(bool); ok { f.Required = required }
							if privacyScreen, ok := customMap["privacy_screen"].(bool); ok { f.PrivacyScreen = privacyScreen }
							if value, ok := customMap["value"].(string); ok && value != "" { f.Value = []string{value} }
							fieldMap = convertFieldToMap(f.Type, f.Label, f.Required, f.PrivacyScreen, f.Value)
						case *core.Multiline:
							if label, ok := customMap["label"].(string); ok { f.Label = label }
							if required, ok := customMap["required"].(bool); ok { f.Required = required }
							if privacyScreen, ok := customMap["privacy_screen"].(bool); ok { f.PrivacyScreen = privacyScreen }
							if value, ok := customMap["value"].(string); ok && value != "" { f.Value = []string{value} }
							fieldMap = convertFieldToMap(f.Type, f.Label, f.Required, f.PrivacyScreen, f.Value)
						case *core.Secret:
							if label, ok := customMap["label"].(string); ok { f.Label = label }
							if required, ok := customMap["required"].(bool); ok { f.Required = required }
							if privacyScreen, ok := customMap["privacy_screen"].(bool); ok { f.PrivacyScreen = privacyScreen }
							if value, ok := customMap["value"].(string); ok && value != "" { f.Value = []string{value} }
							fieldMap = convertFieldToMap(f.Type, f.Label, f.Required, f.PrivacyScreen, f.Value)
						case *core.Url:
							if label, ok := customMap["label"].(string); ok { f.Label = label }
							if required, ok := customMap["required"].(bool); ok { f.Required = required }
							if privacyScreen, ok := customMap["privacy_screen"].(bool); ok { f.PrivacyScreen = privacyScreen }
							if value, ok := customMap["value"].(string); ok && value != "" { f.Value = []string{value} }
							fieldMap = convertFieldToMap(f.Type, f.Label, f.Required, f.PrivacyScreen, f.Value)
						case *core.Email:
							if label, ok := customMap["label"].(string); ok { f.Label = label }
							if required, ok := customMap["required"].(bool); ok { f.Required = required }
							if privacyScreen, ok := customMap["privacy_screen"].(bool); ok { f.PrivacyScreen = privacyScreen }
							if value, ok := customMap["value"].(string); ok && value != "" { f.Value = []string{value} }
							fieldMap = convertFieldToMap(f.Type, f.Label, f.Required, f.PrivacyScreen, f.Value)
						}

						customFields = append(customFields, fieldMap)
					}
				}
			}
		}
		secret.RecordDict["custom"] = customFields
	}

	secret.RawJson = core.DictToJson(secret.RecordDict)
	if err := saveRecord(secret, client); err != nil {
		return diag.FromErr(err)
	}

	return resourcePamDatabaseRead(ctx, d, m)
}

func resourcePamDatabaseDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics

	uid := strings.TrimSpace(d.Get("uid").(string))
	if uid == "" {
		return diag.Errorf("'uid' is required to delete existing resource")
	}

	if err := deleteRecord(uid, client); err != nil {
		if strings.HasSuffix(err.Error(), "unexpected status: ''") {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Record UID: %s not found - probably already deleted (externally)", uid),
				Detail: fmt.Sprintf("Delete record UID: %s returned empty status."+
					" That usually means the record doesn't exist -"+
					" either already deleted (externally),"+
					" or no longer shared to the corresponding KSM Application.", uid),
			})
		} else {
			return diag.FromErr(err)
		}
	}

	d.SetId("")
	return diags
}
func resourcePamDatabaseImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	uid := d.Id()
	if strings.TrimSpace(uid) == "" {
		return nil, errors.New("'uid' is required to import resource")
	}

	if err := d.Set("uid", uid); err != nil {
		return nil, err
	}

	diags := resourcePamDatabaseRead(ctx, d, m)
	if diags.HasError() {
		for _, d := range diags {
			if d.Severity == diag.Error {
				return nil, fmt.Errorf("error reading PAM Database: %s", d.Summary)
			}
		}
	}

	return []*schema.ResourceData{d}, nil
}
