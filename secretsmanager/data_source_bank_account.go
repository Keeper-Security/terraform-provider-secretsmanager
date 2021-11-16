package secretsmanager

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBankAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBankAccountRead,
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
			"bank_account": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The secret bank account information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The account type.",
						},
						"other_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The other type.",
						},
						"routing_number": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "The routing number.",
						},
						"account_number": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "The account number.",
						},
					},
				},
			},
			"name": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The secret name.",
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
			"card_ref": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The secret card reference.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The card ref UID.",
						},
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
					},
				},
			},
			"file_ref": {
				Type:        schema.TypeList,
				Computed:    true,
				Sensitive:   true,
				Description: "The secret file reference.",
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

func dataSourceBankAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics

	path := strings.TrimSpace(d.Get("path").(string))
	title := strings.TrimSpace(d.Get("title").(string))
	secret, err := getRecord(path, title, client)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceType := "bankAccount"
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

	acctItems := getBankAccountItemData(secret)
	if err = d.Set("bank_account", acctItems); err != nil {
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

	// Missing external reference is not an error:
	// - cardRef is not a required field so empty external cardRef field is valid
	// - external cardRef UID present but its record may not be shared to the app or externally deleted
	if cardRef := strings.TrimSpace(secret.GetFieldValueByType("cardRef")); cardRef != "" {
		cardItems := []interface{}{map[string]interface{}{"uid": cardRef}}
		if secretCardRefs, err := client.GetSecrets([]string{cardRef}); err == nil && len(secretCardRefs) > 0 {
			cardItems = getCardRefItemData(secretCardRefs[0], cardRef)
		}
		if err = d.Set("card_ref", cardItems); err != nil {
			return diag.FromErr(err)
		}
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
	if err = d.Set("totp", totpItems); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(path)

	return diags
}
