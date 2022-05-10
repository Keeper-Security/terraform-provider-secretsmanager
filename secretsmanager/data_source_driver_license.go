package secretsmanager

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDriverLicense() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDriverLicenseRead,
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
			"driver_license_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Driver's License Number.",
			},
			"name": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The name.",
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
			"expiration_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date of expiration.",
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

func dataSourceDriverLicenseRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics

	path := strings.TrimSpace(d.Get("path").(string))
	title := strings.TrimSpace(d.Get("title").(string))
	secret, err := getRecord(path, title, client)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceType := "driverLicense"
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

	if err = d.Set("driver_license_number", secret.GetFieldValueByType("accountNumber")); err != nil {
		return diag.FromErr(err)
	}

	nameItems := getNameItemData(secret)
	if err = d.Set("name", nameItems); err != nil {
		return diag.FromErr(err)
	}

	// TF timestamp() uses RFC3339
	bdate := secret.GetFieldValueByType("birthDate")
	if unixTime, err := strconv.Atoi(bdate); err == nil {
		birthDate := time.Unix(int64(unixTime/1000), 0).Format(time.RFC3339)
		if err = d.Set("birth_date", birthDate); err != nil {
			return diag.FromErr(err)
		}
	}

	edate := secret.GetFieldValueByType("expirationDate")
	if unixTime, err := strconv.Atoi(edate); err == nil {
		expirationDate := time.Unix(int64(unixTime/1000), 0).Format(time.RFC3339)
		if err = d.Set("expiration_date", expirationDate); err != nil {
			return diag.FromErr(err)
		}
	}

	// Missing external reference is not an error:
	// - addresRef is not a required field so empty external addresRef field is valid
	// - external addredsRef UID present but its record may not be shared to the app or externally deleted
	if addressRef := strings.TrimSpace(secret.GetFieldValueByType("addressRef")); addressRef != "" {
		addrItems := []interface{}{map[string]interface{}{"uid": addressRef}}
		if secretAddrRefs, err := client.GetSecrets([]string{addressRef}); err == nil && len(secretAddrRefs) > 0 {
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
