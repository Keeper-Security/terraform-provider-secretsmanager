package secretsmanager

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/keeper-security/secrets-manager-go/core"
)

func resourceBirthCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBirthCertificateCreate,
		ReadContext:   resourceBirthCertificateRead,
		UpdateContext: resourceBirthCertificateUpdate,
		DeleteContext: resourceBirthCertificateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceBirthCertificateImport,
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
			// fields[]
			"name":       schemaNameField(),
			"birth_date": schemaBirthDateField(),
			"file_ref":   schemaFileRefField(),
			"custom": schemaCustomField(),
		},
	}
}

func resourceBirthCertificateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	nrc := core.NewRecordCreate("birthCertificate", "")
	if title := d.Get("title"); title != nil && title.(string) != "" {
		nrc.Title = title.(string)
	}
	if notes := d.Get("notes"); notes != nil && notes.(string) != "" {
		nrc.Notes = notes.(string)
	}

	if fieldData := d.Get("name"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("name", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			nrc.Fields = append(nrc.Fields, field)
			if err := SetFieldTypeInSchema(d, "name", "name"); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if fieldData := d.Get("birth_date"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("birthDate", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			nrc.Fields = append(nrc.Fields, field)
			if err := SetFieldTypeInSchema(d, "birth_date", "birthDate"); err != nil {
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
	if err = d.Set("type", "birthCertificate"); err != nil {
		return diag.FromErr(err)
	}


	d.SetId(uid)
	return diags
}

func resourceBirthCertificateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
			// resource does not exist in the vault
			d.SetId("")  // mark for removal
			return diags // no error
		}
		return diag.FromErr(err)
	}

	resourceType := "birthCertificate"
	recordType := secret.Type()
	if recordType != resourceType {
		return diag.Errorf("record type '%s' is not the expected type '%s' for this data source", recordType, resourceType)
	}
	if uid == "" { // found by title: uid="", path="*"
		if err = d.Set("uid", secret.Uid); err != nil {
			return diag.FromErr(err)
		}
	}
	fuid := secret.InnerFolderUid() // in subfolder
	if fuid == "" {                 // directly in shared folder
		fuid = secret.FolderUid()
	}
	if fuid != "" {
		if err = d.Set("folder_uid", fuid); err != nil {
			return diag.FromErr(err)
		}
	} // else - directly shared to the KSM App (not through shared folder)
	if err = d.Set("type", recordType); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("title", secret.Title()); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("notes", secret.Notes()); err != nil {
		return diag.FromErr(err)
	}

	name := getFieldResourceData("name", "fields", secret)
	if err = d.Set("name", name); err != nil {
		return diag.FromErr(err)
	}
	birthDate := getFieldResourceData("birthDate", "fields", secret)
	if err = d.Set("birth_date", birthDate); err != nil {
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

func resourceBirthCertificateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics
	// Validate custom field labels are unique
	if validateDiags := validateUniqueCustomFieldLabels(d); len(validateDiags) > 0 {
		return validateDiags
	}


	// folderUid := strings.TrimSpace(d.Get("folder_uid").(string))
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

	if d.HasChange("name") {
		if _, err := ApplyFieldChange("fields", "name", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("birth_date") {
		if _, err := ApplyFieldChange("fields", "birth_date", d, secret); err != nil {
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

	d.SetId(uid)
	return diags
}

func resourceBirthCertificateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics

	uid := strings.TrimSpace(d.Get("uid").(string))
	if uid == "" {
		return diag.Errorf("'uid' is required to delete existing resource")
	}

	if err := deleteRecord(uid, client); err != nil {
		if strings.HasSuffix(err.Error(), "unexpected status: ''") {
			// record UID no longer exists - probably deleted externally
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
	// NB! Do not return an error if resource already deleted by the vault/app
	// This allows users to manually delete resources without breaking Terraform.
	d.SetId("")
	return diags
}

func resourceBirthCertificateImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	uid := d.Id()

	err := d.Set("uid", uid)
	if err != nil {
		return nil, err
	}

	diags := resourceBirthCertificateRead(ctx, d, m)
	if diags.HasError() {
		for i := range diags {
			if diags[i].Severity == diag.Error {
				return nil, errors.New(diags[i].Summary + " *** Details: " + diags[i].Detail)
			}
		}
	}

	return []*schema.ResourceData{d}, nil
}
