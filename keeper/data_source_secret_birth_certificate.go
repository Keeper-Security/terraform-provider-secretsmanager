package keeper

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSecretBirthCertificate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecretBirthCertificateRead,
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
			"birth_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date of birth.",
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
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The file URL.",
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

func dataSourceSecretBirthCertificateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics

	path := strings.TrimSpace(d.Get("path").(string))
	title := strings.TrimSpace(d.Get("title").(string))
	secret, err := getRecord(path, title, client)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceType := "birthCertificate"
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

	nameItems := getNameItemData(secret)
	if err = d.Set("name", nameItems); err != nil {
		return diag.FromErr(err)
	}

	bdate := secret.GetFieldValueByType("birthDate")
	if unixTime, err := strconv.Atoi(bdate); err == nil {
		// TF timestamp() uses RFC3339
		birthDate := time.Unix(int64(unixTime/1000), 0).Format(time.RFC3339)
		if err = d.Set("birth_date", birthDate); err != nil {
			return diag.FromErr(err)
		}
	}

	fileItems := getFileItemsData(secret.Files)
	if err := d.Set("file_ref", fileItems); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(path)

	return diags
}
