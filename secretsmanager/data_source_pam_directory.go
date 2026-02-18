package secretsmanager

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePamDirectory() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePamDirectoryRead,
		Schema: map[string]*schema.Schema{
			"path": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The UID or KSM notation path to the PAM Directory secret (e.g., record UID or UID/field/password).",
				ExactlyOneOf: []string{"path", "title"},
			},
			"title": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The title of PAM Directory secret to search.",
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
			// PAM Directory specific fields
			"pam_hostname":       schemaPamHostnameField(),
			"pam_settings":       schemaPamSettingsField(),
			"directory_type":     schemaDirectoryTypeField(),
			"rotation_scripts":   schemaScriptField(),
			"use_ssl":            schemaCheckboxField(),
			"distinguished_name": schemaTextField(),
			"domain_name":        schemaTextField(),
			"directory_id":       schemaTextField(),
			"user_match":         schemaTextField(),
			"provider_group":     schemaTextField(),
			"provider_region":    schemaTextField(),
			"alternative_ips":    schemaMultilineField(),
			"file_ref":           schemaFileRefField(),
			"totp":               schemaOneTimeCodeField(),
		},
	}
}

func dataSourcePamDirectoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics

	path := strings.TrimSpace(d.Get("path").(string))
	title := strings.TrimSpace(d.Get("title").(string))
	secret, err := getRecord(path, title, client)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceType := "pamDirectory"
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

	// PAM Directory specific fields
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
	if directoryTypeFields := secret.GetFieldsByType("directoryType"); len(directoryTypeFields) > 0 {
		fieldMap := directoryTypeFields[0]
		if valueInterface, exists := fieldMap["value"]; exists {
			if valueList, ok := valueInterface.([]interface{}); ok && len(valueList) > 0 {
				if dirType, ok := valueList[0].(string); ok {
					if err = d.Set("directory_type", dirType); err != nil {
						return diag.FromErr(err)
					}
				}
			}
		}
	}
	rotationScripts := getFieldResourceDataWithLabel("script", "fields", secret, "Rotation Scripts")
	if err = d.Set("rotation_scripts", rotationScripts); err != nil {
		return diag.FromErr(err)
	}
	useSSL := getFieldResourceDataWithLabel("checkbox", "fields", secret, "useSSL")
	if err = d.Set("use_ssl", useSSL); err != nil {
		return diag.FromErr(err)
	}
	distinguishedName := getFieldResourceDataWithLabel("text", "fields", secret, "Distinguished Name")
	if err = d.Set("distinguished_name", distinguishedName); err != nil {
		return diag.FromErr(err)
	}
	domainName := getFieldResourceDataWithLabel("text", "fields", secret, "domainName")
	if err = d.Set("domain_name", domainName); err != nil {
		return diag.FromErr(err)
	}
	directoryId := getFieldResourceDataWithLabel("text", "fields", secret, "directoryId")
	if err = d.Set("directory_id", directoryId); err != nil {
		return diag.FromErr(err)
	}
	userMatch := getFieldResourceDataWithLabel("text", "fields", secret, "userMatch")
	if err = d.Set("user_match", userMatch); err != nil {
		return diag.FromErr(err)
	}
	providerGroup := getFieldResourceDataWithLabel("text", "fields", secret, "providerGroup")
	if err = d.Set("provider_group", providerGroup); err != nil {
		return diag.FromErr(err)
	}
	providerRegion := getFieldResourceDataWithLabel("text", "fields", secret, "providerRegion")
	if err = d.Set("provider_region", providerRegion); err != nil {
		return diag.FromErr(err)
	}
	alternativeIPs := getFieldResourceDataWithLabel("multiline", "fields", secret, "alternativeIPs")
	if err = d.Set("alternative_ips", alternativeIPs); err != nil {
		return diag.FromErr(err)
	}

	fileItems := getFileItemsData(secret.Files)
	if err := d.Set("file_ref", fileItems); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(path)

	return diags
}
