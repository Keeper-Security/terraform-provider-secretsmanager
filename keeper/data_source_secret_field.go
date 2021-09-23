package keeper

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSecretField() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecretFieldRead,
		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The path where the secret is stored.",
			},
			"title": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The secret title. (To find record by title - replace UID in path with '*')",
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The value of the secret field.",
			},
		},
	}
}

func dataSourceSecretFieldRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics

	path := strings.TrimSpace(d.Get("path").(string))
	title := strings.TrimSpace(d.Get("title").(string))

	// find by title requested
	if title != "" && strings.Contains(path, "*") {
		uids := []string{}
		records, err := client.GetSecrets([]string{})
		if err != nil {
			return diag.FromErr(err)
		}
		for _, r := range records {
			if r.Title() == title {
				uids = append(uids, r.Uid)
			}
		}
		if len(uids) != 1 {
			return diag.Errorf("expected 1 record - found %d records with title: %s", len(uids), title)
		}
		// replace placeholder with the record UID
		path = strings.Replace(path, "*", uids[0], 1)
		if err = d.Set("path", path); err != nil {
			return diag.FromErr(err)
		}
	}

	value, err := client.GetNotation(path)
	if err != nil {
		return diag.FromErr(err)
	}

	strValue := ""
	if len(value) == 1 {
		strValue = fmt.Sprintf("%v", value[0])
	} else {
		strValue = fmt.Sprintf("%v", value)
	}

	if err = d.Set("value", strValue); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(path)

	return diags
}
