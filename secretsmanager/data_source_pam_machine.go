package secretsmanager

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePamMachine() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePamMachineRead,
		Schema: map[string]*schema.Schema{
			"path": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The UID or KSM notation path to the PAM Machine secret (e.g., record UID or UID/field/password).",
				ExactlyOneOf: []string{"path", "title"},
			},
			"title": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The title of PAM Machine secret to search.",
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
			// PAM Machine specific fields
			"pam_hostname":     schemaPamHostnameField(),
			"pam_settings":     schemaPamSettingsField(),
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

func dataSourcePamMachineRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics

	path := strings.TrimSpace(d.Get("path").(string))
	title := strings.TrimSpace(d.Get("title").(string))
	secret, err := getRecord(path, title, client)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceType := "pamMachine"
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

	fileItems := getFileItemsData(secret.Files)
	if err := d.Set("file_ref", fileItems); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(path)

	return diags
}
