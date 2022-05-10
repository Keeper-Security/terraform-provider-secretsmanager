package secretsmanager

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceHealthInsurance() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHealthInsuranceRead,
		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The path where the secret is stored.",
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
				Computed:    true,
				Description: "The secret notes.",
			},
			// fields[]
			"account_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Account Number.",
			},
			"name": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Full name.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"first": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "First name.",
						},
						"middle": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Middle name.",
						},
						"last": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last name.",
						},
					},
				},
			},
			"login": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The secret login.",
			},
			"password": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The secret password.",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The secret url.",
			},
			"file_ref": {
				Type:        schema.TypeList,
				Computed:    true,
				Sensitive:   true,
				Description: "The secret file references",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The file ref UID.",
						},
						"title": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The file title.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The file name.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The file type.",
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The file size.",
						},
						"last_modified": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The file last modified date.",
						},
						"content_base64": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The file content (base64).",
						},
					},
				},
			},
		},
	}
}

func dataSourceHealthInsuranceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics

	path := strings.TrimSpace(d.Get("path").(string))
	title := strings.TrimSpace(d.Get("title").(string))
	secret, err := getRecord(path, title, client)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceType := "healthInsurance"
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

	if err = d.Set("account_number", secret.GetFieldValueByType("accountNumber")); err != nil {
		return diag.FromErr(err)
	}

	nameItems := getNameItemData(secret)
	if err = d.Set("name", nameItems); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("login", secret.GetFieldValueByType("login")); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("password", secret.GetFieldValueByType("password")); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("url", secret.GetFieldValueByType("url")); err != nil {
		return diag.FromErr(err)
	}

	fileItems := getFileItemsData(secret.Files)
	if err := d.Set("file_ref", fileItems); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(path)

	return diags
}
