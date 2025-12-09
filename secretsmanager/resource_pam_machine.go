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

func resourcePamMachine() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePamMachineCreate,
		ReadContext:   resourcePamMachineRead,
		UpdateContext: resourcePamMachineUpdate,
		DeleteContext: resourcePamMachineDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePamMachineImport,
		},
		Schema: map[string]*schema.Schema{
			"folder_uid": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				AtLeastOneOf: []string{"folder_uid", "uid"},
				Description:  "The UID of the folder where the secret is stored. The folder or its parent shared folder must be accessible to your KSM application with 'Can Edit' permissions.",
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
			// PAM Machine specific fields
			"pam_hostname":     schemaPamHostnameField(),
			"pam_settings":     schemaPamSettingsField(),
			"login":            schemaLoginField(),
			"password":         schemaPasswordField(""),
			"rotation_scripts": schemaScriptField(),
			"operating_system": schemaTextField(),
			"instance_name":    schemaTextField(),
			"instance_id":      schemaTextField(),
			"provider_group":   schemaTextField(),
			"provider_region":  schemaTextField(),
			"file_ref":         schemaFileRefField(),
			"totp":             schemaOneTimeCodeField(),
		},
	}
}

func resourcePamMachineCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	nrc := core.NewRecordCreate("pamMachine", "")
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
	if fieldData := d.Get("operating_system"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("text", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			field.(*core.Text).Label = "Operating System"
			nrc.Fields = append(nrc.Fields, field)
			if err := SetFieldTypeInSchema(d, "operating_system", "text"); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if fieldData := d.Get("instance_name"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("text", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			field.(*core.Text).Label = "Instance Name"
			nrc.Fields = append(nrc.Fields, field)
			if err := SetFieldTypeInSchema(d, "instance_name", "text"); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if fieldData := d.Get("instance_id"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("text", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			field.(*core.Text).Label = "Instance Id"
			nrc.Fields = append(nrc.Fields, field)
			if err := SetFieldTypeInSchema(d, "instance_id", "text"); err != nil {
				return diag.FromErr(err)
			}
		}
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
	if err = d.Set("type", "pamMachine"); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(uid)
	return diags
}

func resourcePamMachineRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	resourceType := "pamMachine"
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

	// PAM Machine specific fields
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
	rotationScripts := getFieldResourceDataWithLabel("script", "fields", secret, "Rotation Scripts")
	if err = d.Set("rotation_scripts", rotationScripts); err != nil {
		return diag.FromErr(err)
	}
	operatingSystem := getFieldResourceDataWithLabel("text", "fields", secret, "Operating System")
	if err = d.Set("operating_system", operatingSystem); err != nil {
		return diag.FromErr(err)
	}
	instanceName := getFieldResourceDataWithLabel("text", "fields", secret, "Instance Name")
	if err = d.Set("instance_name", instanceName); err != nil {
		return diag.FromErr(err)
	}
	instanceId := getFieldResourceDataWithLabel("text", "fields", secret, "Instance Id")
	if err = d.Set("instance_id", instanceId); err != nil {
		return diag.FromErr(err)
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

	d.SetId(uid)
	return diags
}

func resourcePamMachineUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	if d.HasChange("rotation_scripts") {
		if _, err := ApplyFieldChange("fields", "rotation_scripts", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("operating_system") {
		if _, err := ApplyFieldChange("fields", "operating_system", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("instance_name") {
		if _, err := ApplyFieldChange("fields", "instance_name", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("instance_id") {
		if _, err := ApplyFieldChange("fields", "instance_id", d, secret); err != nil {
			return diag.FromErr(err)
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

	if err := saveRecord(secret, client); err != nil {
		return diag.FromErr(err)
	}

	return resourcePamMachineRead(ctx, d, m)
}

func resourcePamMachineDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
func resourcePamMachineImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	uid := d.Id()
	if strings.TrimSpace(uid) == "" {
		return nil, errors.New("'uid' is required to import resource")
	}

	if err := d.Set("uid", uid); err != nil {
		return nil, err
	}

	diags := resourcePamMachineRead(ctx, d, m)
	if diags.HasError() {
		for _, d := range diags {
			if d.Severity == diag.Error {
				return nil, fmt.Errorf("error reading PAM Machine: %s", d.Summary)
			}
		}
	}

	return []*schema.ResourceData{d}, nil
}
