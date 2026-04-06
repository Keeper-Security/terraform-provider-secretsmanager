package secretsmanager

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePamRemoteBrowser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePamRemoteBrowserRead,
		Schema: map[string]*schema.Schema{
			"path": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The UID or KSM notation path to the PAM Remote Browser secret (e.g., record UID or UID/field/rbiUrl).",
				ExactlyOneOf: []string{"path", "title"},
			},
			"title": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The title of PAM Remote Browser secret to search.",
				ExactlyOneOf: []string{"path", "title"},
			},
			"folder_uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The folder UID where the secret is stored.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The secret type.",
			},
			"notes": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The secret notes.",
			},
			// PAM Remote Browser specific fields
			"rbi_url": schemaTextField(),
			"pam_remote_browser_settings": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "PAM Remote Browser connection settings as JSON string.",
			},
			"traffic_encryption_seed": schemaTextField(),
			"file_ref":                schemaFileRefField(),
			"totp":                    schemaOneTimeCodeField(),
		},
	}
}

func dataSourcePamRemoteBrowserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics

	path := strings.TrimSpace(d.Get("path").(string))
	title := strings.TrimSpace(d.Get("title").(string))
	secret, err := getRecord(path, title, client)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceType := "pamRemoteBrowser"
	recordType := secret.Type()
	if recordType != dataSourceType {
		return diag.Errorf("record type '%s' is not the expected type '%s' for this data source", recordType, dataSourceType)
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
	fuid := secret.InnerFolderUid()
	if fuid == "" {
		fuid = secret.FolderUid()
	}
	if fuid != "" {
		if err = d.Set("folder_uid", fuid); err != nil {
			return diag.FromErr(err)
		}
	}

	// PAM Remote Browser specific fields
	// rbi_url: custom field type stored in "rbiUrl" — read via GetFieldsByType
	if rbiUrlFields := secret.GetFieldsByType("rbiUrl"); len(rbiUrlFields) > 0 {
		fieldMap := rbiUrlFields[0]
		if valueInterface, exists := fieldMap["value"]; exists {
			if valueList, ok := valueInterface.([]interface{}); ok && len(valueList) > 0 {
				if rbiUrl, ok := valueList[0].(string); ok && rbiUrl != "" {
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
	// pam_remote_browser_settings: JSON dict field stored in "pamRemoteBrowserSettings"
	if pamRbsFields := secret.GetFieldsByType("pamRemoteBrowserSettings"); len(pamRbsFields) > 0 {
		if pamRbsJSON, err := pamSettingsFieldToJSON(pamRbsFields[0]); err != nil {
			return diag.FromErr(err)
		} else if err = d.Set("pam_remote_browser_settings", pamRbsJSON); err != nil {
			return diag.FromErr(err)
		}
	}
	// traffic_encryption_seed: custom field type stored in "trafficEncryptionSeed"
	if trafficFields := secret.GetFieldsByType("trafficEncryptionSeed"); len(trafficFields) > 0 {
		fieldMap := trafficFields[0]
		if valueInterface, exists := fieldMap["value"]; exists {
			if valueList, ok := valueInterface.([]interface{}); ok && len(valueList) > 0 {
				if seed, ok := valueList[0].(string); ok && seed != "" {
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
	fileItems := getFileItemsData(secret.Files)
	if err := d.Set("file_ref", fileItems); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(path)

	return diags
}
