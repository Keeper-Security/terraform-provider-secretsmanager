package secretsmanager

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMetadata() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMetadataRead,
		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The record UID, or \"*\" to look up the record by title.",
			},
			"uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The record UID.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The record type (e.g. \"login\").",
			},
			"title": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The record title. Used to look up the record when path is \"*\".",
			},
			"notes": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The record notes.",
			},
			"revision": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: "Record revision counter. Increments on every modification (value, title, and notes changes all bump it). " +
					"Use as a version signal paired with write-only attributes on other providers — when the Keeper record rotates, " +
					"revision changes, triggering Terraform to re-apply the write-only value on the next plan.",
			},
			"folder_uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UID of the folder containing this record.",
			},
			"is_editable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the KSM application credential can edit this record.",
			},
		},
	}
}

func dataSourceMetadataRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics

	path := strings.TrimSpace(d.Get("path").(string))
	title := strings.TrimSpace(d.Get("title").(string))
	secret, err := getRecord(path, title, client)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("uid", secret.Uid); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("type", secret.Type()); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("title", secret.Title()); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("notes", secret.Notes()); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("revision", int(secret.Revision)); err != nil {
		return diag.FromErr(err)
	}
	fuid := secret.InnerFolderUid()
	if fuid == "" {
		fuid = secret.FolderUid()
	}
	if err = d.Set("folder_uid", fuid); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("is_editable", secret.IsEditable); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(path)
	return diags
}
