package secretsmanager

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePamUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePamUserRead,
		Schema: map[string]*schema.Schema{
			"path": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The path to PAM User secret.",
				ExactlyOneOf: []string{"path", "title"},
			},
			"title": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The title of PAM User secret to search.",
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
			// PAM User specific fields
			"login":              schemaLoginField(),
			"password":           schemaPasswordField(""),
			"rotation_scripts":   schemaScriptField(),
			"distinguished_name": schemaTextField(),
			"connect_database":   schemaTextField(),
			"managed":            schemaCheckboxField(),
			"file_ref":           schemaFileRefField(),
			"totp":               schemaOneTimeCodeField(),
		},
	}
}

func dataSourcePamUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics

	path := strings.TrimSpace(d.Get("path").(string))
	title := strings.TrimSpace(d.Get("title").(string))
	secret, err := getRecord(path, title, client)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceType := "pamUser"
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
	if err = d.Set("login", secret.GetFieldValueByType("login")); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("password", secret.GetFieldValueByType("password")); err != nil {
		return diag.FromErr(err)
	}

	// PAM User specific fields
	rotationScripts := getFieldResourceDataWithLabel("script", "fields", secret, "Rotation Scripts")
	if err = d.Set("rotation_scripts", rotationScripts); err != nil {
		return diag.FromErr(err)
	}
	distinguishedName := getFieldResourceDataWithLabel("text", "fields", secret, "Distinguished Name")
	if err = d.Set("distinguished_name", distinguishedName); err != nil {
		return diag.FromErr(err)
	}
	connectDatabase := getFieldResourceDataWithLabel("text", "fields", secret, "Connect Database")
	if err = d.Set("connect_database", connectDatabase); err != nil {
		return diag.FromErr(err)
	}
	managed := getFieldResourceDataWithLabel("checkbox", "fields", secret, "Managed")
	if err = d.Set("managed", managed); err != nil {
		return diag.FromErr(err)
	}

	fileItems := getFileItemsData(secret.Files)
	if err := d.Set("file_ref", fileItems); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(path)

	return diags
}
