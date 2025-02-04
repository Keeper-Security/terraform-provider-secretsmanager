package secretsmanager

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceFolder() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFolderRead,
		Schema: map[string]*schema.Schema{
			"uid": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"name"},
				Description:  "The folder uid.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"uid"},
				Description:  "The folder name.",
			},
			"parent_uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The parent folder uid.",
			},
			"shared": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Shared folder flag.",
			},
			// Following official documentation:
			// Plural data sources should return zero, one, or multiple results without error.
			// Singular data sources should still return an error when the remote component is not found.
			// --
			// Terraform cannot differentiate between an legitimate error situation and a missing resource.
			// A data source is meant to represent a single unmanaged remote resource,
			// and it should normally return an error, because the necessary resource is expected to be found.
			// To handle this in the configuration, the data source could be expanded:
			// "allow_missing_resource": {Type: schema.TypeBool, Optional: true, Description: "When not found populate exists flag instead of error."},
			// "exists":                 {Type: schema.TypeBool, Computed: true, Description: "The found/exists flag."},
		},
	}
}

func dataSourceFolderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics

	parentUid := strings.TrimSpace(d.Get("parent_uid").(string))
	uid := strings.TrimSpace(d.Get("uid").(string))
	name := strings.TrimSpace(d.Get("name").(string))
	folders, err := findFolder(parentUid, uid, name, client)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(folders) == 0 {
		return diag.Errorf("folder UID: '%s', Name: '%s' not found", uid, name)
	} else if len(folders) > 1 {
		return diag.Errorf("multiple folders (%d) match folder UID: '%s', Name: '%s'", len(folders), uid, name)
	}

	folder := folders[0]
	if err = d.Set("parent_uid", folder.ParentUid); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("shared", strings.TrimSpace(folder.ParentUid) == ""); err != nil {
		return diag.FromErr(err)
	}
	if uid == "" {
		if err = d.Set("uid", folder.FolderUid); err != nil {
			return diag.FromErr(err)
		}
	}
	if name == "" {
		if err = d.Set("name", folder.Name); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(folder.FolderUid)

	return diags
}
