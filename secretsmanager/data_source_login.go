package secretsmanager

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLogin() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLoginRead,
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
			"totp": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The one time password.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "TOTP URL.",
						},
						"token": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "Generated TOTP token.",
						},
						"ttl": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Time to live for TOTP token in seconds.",
						},
					},
				},
			},
		},
	}
}

func dataSourceLoginRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics

	path := strings.TrimSpace(d.Get("path").(string))
	title := strings.TrimSpace(d.Get("title").(string))
	secret, err := getRecord(path, title, client)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceType := "login"
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
	if err = d.Set("url", secret.GetFieldValueByType("url")); err != nil {
		return diag.FromErr(err)
	}

	fileItems := getFileItemsData(secret.Files)
	if err := d.Set("file_ref", fileItems); err != nil {
		return diag.FromErr(err)
	}

	totpItems := []interface{}{}
	if totp := strings.TrimSpace(secret.GetFieldValueByType("oneTimeCode")); totp != "" {
		if code, seconds, err := getTotpCode(totp); err != nil {
			return diag.FromErr(err)
		} else {
			totpItems = []interface{}{
				map[string]interface{}{
					"url":   totp,
					"token": code,
					"ttl":   seconds,
				},
			}
		}
	}
	if err := d.Set("totp", totpItems); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(path)

	return diags
}
