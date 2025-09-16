package secretsmanager

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBankCard() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBankCardRead,
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
			"payment_card": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The payment card information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"card_number": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "The card number.",
						},
						"card_expiration_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "The card expiration date.",
						},
						"card_security_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "The card security code.",
						},
					},
				},
			},
			"cardholder_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The cardholder name.",
			},
			"pin_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The PIN code.",
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

func dataSourceBankCardRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics

	path := strings.TrimSpace(d.Get("path").(string))
	title := strings.TrimSpace(d.Get("title").(string))
	secret, err := getRecord(path, title, client)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceType := "bankCard"
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

	cardItems := getPaymentCardItemData(secret)
	if err = d.Set("payment_card", cardItems); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("cardholder_name", secret.GetFieldValueByType("text")); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("pin_code", secret.GetFieldValueByType("pinCode")); err != nil {
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
