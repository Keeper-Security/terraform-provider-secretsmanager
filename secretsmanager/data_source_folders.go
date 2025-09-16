package secretsmanager

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceFolders() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFoldersRead,
		Schema: map[string]*schema.Schema{
			"folders": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of all folders shared to the KSM Application.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The folder UID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The folder name.",
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
					},
				},
			},
		},
	}
}

func dataSourceFoldersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics

	folders, err := getFolders(client)
	if err != nil {
		return diag.FromErr(err)
	}

	folderItems := []interface{}{}
	if len(folders) > 0 {
		fis := make([]interface{}, len(folders))
		for i, folderItem := range folders {
			fi := map[string]interface{}{}
			fi["uid"] = folderItem.FolderUid
			fi["name"] = folderItem.Name
			fi["parent_uid"] = folderItem.ParentUid
			fi["shared"] = strings.TrimSpace(folderItem.ParentUid) == ""
			fis[i] = fi
		}
		folderItems = fis
	}

	if err := d.Set("folders", folderItems); err != nil {
		return diag.FromErr(err)
	}

	// folder list could change any time so just use creation timestamp
	if d.Id() == "" {
		id := fmt.Sprintf("%x", time.Now().UTC().UnixNano())
		d.SetId(id)
	}

	return diags
}
