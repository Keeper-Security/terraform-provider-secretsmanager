package secretsmanager

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSshKeys() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSshKeysRead,
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
			"key_pair": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Public and private key pair.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"public_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The public key.",
						},
						"private_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "The private key.",
						},
					},
				},
			},
			"passphrase": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The passphrase to unlock the key.",
			},
			"host": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Hostname and port",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The hostname.",
						},
						"port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The port.",
						},
					},
				},
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

func dataSourceSshKeysRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics

	path := strings.TrimSpace(d.Get("path").(string))
	title := strings.TrimSpace(d.Get("title").(string))
	secret, err := getRecord(path, title, client)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceType := "sshKeys"
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

	pairItems := getKeyPairItemData(secret)
	if err := d.Set("key_pair", pairItems); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("passphrase", secret.GetFieldValueByType("password")); err != nil {
		return diag.FromErr(err)
	}

	hostItems := getHostItemData(secret)
	if err := d.Set("host", hostItems); err != nil {
		return diag.FromErr(err)
	}

	fileItems := getFileItemsData(secret.Files)
	if err := d.Set("file_ref", fileItems); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(path)

	return diags
}
