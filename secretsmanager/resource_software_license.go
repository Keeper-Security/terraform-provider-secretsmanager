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

func resourceSoftwareLicense() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSoftwareLicenseCreate,
		ReadContext:   resourceSoftwareLicenseRead,
		UpdateContext: resourceSoftwareLicenseUpdate,
		DeleteContext: resourceSoftwareLicenseDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSoftwareLicenseImport,
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
			"license_number":  schemaLicenseNumberField(),
			"activation_date": schemaDateField(),
			"expiration_date": schemaExpirationDateField(),
			"file_ref":        schemaFileRefField(),
			"custom": schemaCustomField(),
		},
	}
}

func resourceSoftwareLicenseCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics

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

	nrc := core.NewRecordCreate("softwareLicense", "")
	if title := d.Get("title"); title != nil && title.(string) != "" {
		nrc.Title = title.(string)
	}
	if notes := d.Get("notes"); notes != nil && notes.(string) != "" {
		nrc.Notes = notes.(string)
	}

	if fieldData := d.Get("license_number"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("licenseNumber", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			nrc.Fields = append(nrc.Fields, field)
			if err := SetFieldTypeInSchema(d, "license_number", "licenseNumber"); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if fieldData := d.Get("activation_date"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("date", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			nrc.Fields = append(nrc.Fields, field)
			if err := SetFieldTypeInSchema(d, "activation_date", "date"); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if fieldData := d.Get("expiration_date"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("expirationDate", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			nrc.Fields = append(nrc.Fields, field)
			if err := SetFieldTypeInSchema(d, "expiration_date", "expirationDate"); err != nil {
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
	if err = d.Set("type", "softwareLicense"); err != nil {
		return diag.FromErr(err)
	}


	d.SetId(uid)
	return diags
}

func resourceSoftwareLicenseRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	resourceType := "softwareLicense"
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

	licenseNumber := getFieldResourceData("licenseNumber", "fields", secret)
	if err = d.Set("license_number", licenseNumber); err != nil {
		return diag.FromErr(err)
	}
	activationDate := getFieldResourceData("date", "fields", secret)
	if err = d.Set("activation_date", activationDate); err != nil {
		return diag.FromErr(err)
	}
	expirationDate := getFieldResourceData("expirationDate", "fields", secret)
	if err = d.Set("expiration_date", expirationDate); err != nil {
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

	d.SetId(uid)
	return diags
}

func resourceSoftwareLicenseUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics

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

	if d.HasChange("license_number") {
		if _, err := ApplyFieldChange("fields", "license_number", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("activation_date") {
		if _, err := ApplyFieldChange("fields", "activation_date", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("expiration_date") {
		if _, err := ApplyFieldChange("fields", "expiration_date", d, secret); err != nil {
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

func resourceSoftwareLicenseDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func resourceSoftwareLicenseImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	uid := d.Id()

	err := d.Set("uid", uid)
	if err != nil {
		return nil, err
	}

	diags := resourceSoftwareLicenseRead(ctx, d, m)
	if diags.HasError() {
		for i := range diags {
			if diags[i].Severity == diag.Error {
				return nil, errors.New(diags[i].Summary + " *** Details: " + diags[i].Detail)
			}
		}
	}

	return []*schema.ResourceData{d}, nil
}
