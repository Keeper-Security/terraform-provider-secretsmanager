package secretsmanager

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePamDatabase() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePamDatabaseRead,
		Schema: map[string]*schema.Schema{
			"path": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The UID or KSM notation path to the PAM Database secret (e.g., record UID or UID/field/password).",
				ExactlyOneOf: []string{"path", "title"},
			},
			"title": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The title of PAM Database secret to search.",
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
			// PAM Database specific fields
			"pam_hostname":     schemaPamHostnameField(),
			"pam_settings":     schemaPamSettingsField(),
			"use_ssl":          schemaCheckboxField(),
			"rotation_scripts": schemaScriptField(),
			"database_id":      schemaTextField(),
			"database_type":    schemaDatabaseTypeField(),
			"provider_group":   schemaTextField(),
			"provider_region":  schemaTextField(),
			"file_ref":         schemaFileRefField(),
			"totp":             schemaOneTimeCodeField(),
		},
	}
}

func dataSourcePamDatabaseRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics

	path := strings.TrimSpace(d.Get("path").(string))
	title := strings.TrimSpace(d.Get("title").(string))
	secret, err := getRecord(path, title, client)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceType := "pamDatabase"
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
