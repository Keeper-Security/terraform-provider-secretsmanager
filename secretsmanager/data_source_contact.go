package secretsmanager

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceContact() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContactRead,
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
				Description: "The contact name.",
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
			"company": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Company name.",
			},
			"email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Contact's e-mail.",
			},
			"phone": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Phone number.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region.",
						},
						"number": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Phone number.",
						},
						"ext": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Extension.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type - Mobile, Home or Work",
						},
					},
				},
			},
			"address_ref": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The address information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The address ref UID.",
						},
						"street1": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Street line one.",
						},
						"street2": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Street line one.",
						},
						"city": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "City.",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "State.",
						},
						"zip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ZIP code.",
						},
						"country": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Country.",
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

func dataSourceContactRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics

	path := strings.TrimSpace(d.Get("path").(string))
	title := strings.TrimSpace(d.Get("title").(string))
	secret, err := getRecord(path, title, client)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceType := "contact"
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

	if err = d.Set("company", secret.GetFieldValueByType("text")); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("email", secret.GetFieldValueByType("email")); err != nil {
		return diag.FromErr(err)
	}

	phoneItems := getPhoneItemData(secret)
	if err = d.Set("phone", phoneItems); err != nil {
		return diag.FromErr(err)
	}
	// Missing external reference is not an error:
	// - addresRef is not a required field so empty external addresRef field is valid
	// - external addredsRef UID present but its record may not be shared to the app or externally deleted
	if addressRef := strings.TrimSpace(secret.GetFieldValueByType("addressRef")); addressRef != "" {
		addrItems := []interface{}{map[string]interface{}{"uid": addressRef}}
		if secretAddrRefs, err := getSecrets(client, []string{addressRef}); err == nil && len(secretAddrRefs) > 0 {
			addrItems = getAddressRefItemData(secretAddrRefs[0], addressRef)
		}
		if err = d.Set("address_ref", addrItems); err != nil {
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
