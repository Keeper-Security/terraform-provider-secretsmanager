package secretsmanager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/keeper-security/secrets-manager-go/core"
)

func resourcePamRemoteBrowser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePamRemoteBrowserCreate,
		ReadContext:   resourcePamRemoteBrowserRead,
		UpdateContext: resourcePamRemoteBrowserUpdate,
		DeleteContext: resourcePamRemoteBrowserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePamRemoteBrowserImport,
		},
		Schema: map[string]*schema.Schema{
			"folder_uid": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				AtLeastOneOf: []string{"folder_uid", "uid"},
				Description:  "The folder UID where the secret is stored. The parent shared folder must be non empty.",
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
			// PAM Remote Browser specific fields
			"rbi_url":                 schemaTextField(),
			"pam_remote_browser_settings": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: suppressEquivalentJSON,
				Description:      "PAM Remote Browser connection settings as JSON string.",
			},
			"traffic_encryption_seed": schemaTextSensitiveField(),
			"file_ref":                schemaFileRefField(),
			"totp":                    schemaOneTimeCodeField(),
			// custom[]
			"custom": schemaCustomField(),
		},
	}
}

// createRbiUrlField creates a raw rbiUrl field from a string value.
func createRbiUrlField(value string) interface{} {
	return &struct {
		core.KeeperRecordField
		Value []string `json:"value"`
	}{
		KeeperRecordField: core.KeeperRecordField{Type: "rbiUrl"},
		Value:             []string{value},
	}
}

// createTrafficEncryptionSeedField creates a raw trafficEncryptionSeed field from a string value.
func createTrafficEncryptionSeedField(value string) interface{} {
	return &struct {
		core.KeeperRecordField
		Value []string `json:"value"`
	}{
		KeeperRecordField: core.KeeperRecordField{Type: "trafficEncryptionSeed"},
		Value:             []string{value},
	}
}

// createPamRemoteBrowserSettingsFieldFromJSON creates a pamRemoteBrowserSettings field from JSON string.
func createPamRemoteBrowserSettingsFieldFromJSON(jsonStr string) (interface{}, error) {
	if jsonStr == "" {
		return nil, nil
	}

	var test interface{}
	if err := json.Unmarshal([]byte(jsonStr), &test); err != nil {
		return nil, fmt.Errorf("failed to parse pam_remote_browser_settings JSON: %w", err)
	}

	field := &struct {
		core.KeeperRecordField
		Value json.RawMessage `json:"value"`
	}{
		KeeperRecordField: core.KeeperRecordField{Type: "pamRemoteBrowserSettings"},
		Value:             json.RawMessage(jsonStr),
	}

	return field, nil
}

// getTextFieldValueFromSchema extracts the string value from a schemaTextField list.
func getTextFieldValueFromSchema(fieldData interface{}) string {
	if s, ok := fieldData.([]interface{}); ok && len(s) > 0 {
		if m, ok := s[0].(map[string]interface{}); ok {
			if v, ok := m["value"].(string); ok {
				return v
			}
		}
	}
	return ""
}

func resourcePamRemoteBrowserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics

	var err error
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

	nrc := core.NewRecordCreate("pamRemoteBrowser", "")
	if title := d.Get("title"); title != nil && title.(string) != "" {
		nrc.Title = title.(string)
	}
	if notes := d.Get("notes"); notes != nil && notes.(string) != "" {
		nrc.Notes = notes.(string)
	}

	// rbi_url: custom field type "rbiUrl" — not in NewFieldFromSchema switch
	if fieldData := d.Get("rbi_url"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if value := getTextFieldValueFromSchema(fieldData); value != "" {
			nrc.Fields = append(nrc.Fields, createRbiUrlField(value))
		}
	}
	// pam_remote_browser_settings: JSON dict field — same pattern as pam_settings
	if pamRbsJSON := d.Get("pam_remote_browser_settings").(string); pamRbsJSON != "" {
		if field, err := createPamRemoteBrowserSettingsFieldFromJSON(pamRbsJSON); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			nrc.Fields = append(nrc.Fields, field)
		}
	}
	// traffic_encryption_seed: custom field type "trafficEncryptionSeed"
	if fieldData := d.Get("traffic_encryption_seed"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if value := getTextFieldValueFromSchema(fieldData); value != "" {
			nrc.Fields = append(nrc.Fields, createTrafficEncryptionSeedField(value))
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

	if customData := d.Get("custom"); customData != nil {
		if fields, err := customFieldsFromSchema(customData.([]interface{})); err != nil {
			return diag.FromErr(err)
		} else {
			nrc.Custom = fields
		}
	}

	if folderUid == "*" {
		if fuid, err := getTemplateFolder(folderUid, client); err != nil {
			return diag.FromErr(err)
		} else {
			folderUid = fuid
		}
	}

	uid, err = createRecord(uid, folderUid, nrc, client)
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
	if err = d.Set("type", "pamRemoteBrowser"); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(uid)
	return diags
}

func resourcePamRemoteBrowserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	resourceType := "pamRemoteBrowser"
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

	// rbi_url: read from custom field type "rbiUrl"
	if rbiUrlFields := secret.GetFieldsByType("rbiUrl"); len(rbiUrlFields) > 0 {
		fieldMap := rbiUrlFields[0]
		if valueInterface, exists := fieldMap["value"]; exists {
			if valueList, ok := valueInterface.([]interface{}); ok && len(valueList) > 0 {
				if rbiUrl, ok := valueList[0].(string); ok {
					rbiUrlData := []interface{}{
						map[string]interface{}{
							"type":           "rbiUrl",
							"label":          "",
							"required":       false,
							"privacy_screen": false,
							"value":          rbiUrl,
						},
					}
					if err = d.Set("rbi_url", rbiUrlData); err != nil {
						return diag.FromErr(err)
					}
				}
			}
		}
	}
	// pam_remote_browser_settings: JSON dict field
	if pamRbsFields := secret.GetFieldsByType("pamRemoteBrowserSettings"); len(pamRbsFields) > 0 {
		if pamRbsJSON, err := pamSettingsFieldToJSON(pamRbsFields[0]); err != nil {
			return diag.FromErr(err)
		} else if err = d.Set("pam_remote_browser_settings", pamRbsJSON); err != nil {
			return diag.FromErr(err)
		}
	}
	// traffic_encryption_seed: read from custom field type "trafficEncryptionSeed"
	if trafficFields := secret.GetFieldsByType("trafficEncryptionSeed"); len(trafficFields) > 0 {
		fieldMap := trafficFields[0]
		if valueInterface, exists := fieldMap["value"]; exists {
			if valueList, ok := valueInterface.([]interface{}); ok && len(valueList) > 0 {
				if seed, ok := valueList[0].(string); ok {
					seedData := []interface{}{
						map[string]interface{}{
							"type":           "trafficEncryptionSeed",
							"label":          "",
							"required":       false,
							"privacy_screen": false,
							"value":          seed,
						},
					}
					if err = d.Set("traffic_encryption_seed", seedData); err != nil {
						return diag.FromErr(err)
					}
				}
			}
		}
	}

	oneTimeCode := getFieldResourceData("oneTimeCode", "fields", secret)
	if err = d.Set("totp", oneTimeCode); err != nil {
		return diag.FromErr(err)
	}

	fileItems := getFileItemsResourceData(secret)
	if err := d.Set("file_ref", fileItems); err != nil {
		return diag.FromErr(err)
	}

	customItems := getFieldItemsData(secret.RecordDict, "custom")
	if err := d.Set("custom", customItems); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(uid)
	return diags
}

func resourcePamRemoteBrowserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client

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

	if d.HasChange("rbi_url") {
		fieldData := d.Get("rbi_url")
		value := getTextFieldValueFromSchema(fieldData)
		if err := secret.SetStandardFieldValue("rbiUrl", []interface{}{value}); err != nil {
			return diag.FromErr(fmt.Errorf("failed to update rbi_url: %w", err))
		}
	}
	if d.HasChange("pam_remote_browser_settings") {
		pamRbsJSON := d.Get("pam_remote_browser_settings").(string)
		var pamRbsValue interface{}
		if err := json.Unmarshal([]byte(pamRbsJSON), &pamRbsValue); err != nil {
			return diag.FromErr(fmt.Errorf("failed to parse pam_remote_browser_settings JSON: %w", err))
		}
		if err := secret.SetStandardFieldValue("pamRemoteBrowserSettings", pamRbsValue); err != nil {
			return diag.FromErr(fmt.Errorf("failed to update pam_remote_browser_settings: %w", err))
		}
	}
	if d.HasChange("traffic_encryption_seed") {
		fieldData := d.Get("traffic_encryption_seed")
		value := getTextFieldValueFromSchema(fieldData)
		if err := secret.SetStandardFieldValue("trafficEncryptionSeed", []interface{}{value}); err != nil {
			return diag.FromErr(fmt.Errorf("failed to update traffic_encryption_seed: %w", err))
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
		customData := d.Get("custom").([]interface{})
		fields, err := customFieldsFromSchema(customData)
		if err != nil {
			return diag.FromErr(err)
		}
		secret.RecordDict["custom"] = customFieldsToDict(fields)
	}

	secret.RawJson = core.DictToJson(secret.RecordDict)
	if err := saveRecord(secret, client); err != nil {
		return diag.FromErr(err)
	}

	return resourcePamRemoteBrowserRead(ctx, d, m)
}

func resourcePamRemoteBrowserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func resourcePamRemoteBrowserImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	uid := d.Id()
	if strings.TrimSpace(uid) == "" {
		return nil, errors.New("'uid' is required to import resource")
	}

	if err := d.Set("uid", uid); err != nil {
		return nil, err
	}

	diags := resourcePamRemoteBrowserRead(ctx, d, m)
	if diags.HasError() {
		for _, d := range diags {
			if d.Severity == diag.Error {
				return nil, fmt.Errorf("error reading PAM Remote Browser: %s", d.Summary)
			}
		}
	}

	return []*schema.ResourceData{d}, nil
}
