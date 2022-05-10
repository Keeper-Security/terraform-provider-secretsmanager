package secretsmanager

import (
	"context"
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/keeper-security/secrets-manager-go/core"
)

func resourceContact() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceContactCreate,
		ReadContext:   resourceContactRead,
		UpdateContext: resourceContactUpdate,
		DeleteContext: resourceContactDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceContactImport,
		},
		Schema: map[string]*schema.Schema{
			"folder_uid": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				AtLeastOneOf: []string{"folder_uid", "uid"},
				Description:  "The folder UID where the secret is stored. The shared folder must be non empty.",
			},
			"uid": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"folder_uid", "uid"},
				Description:  "The UID of the new secret (using RFC4648 URL and Filename Safe Alphabet).",
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
				Optional:    true,
				Description: "The secret notes.",
			},
			// fields[]
			"name":        schemaNameField(),
			"company":     schemaTextField(),
			"email":       schemaEmailField(),
			"phone":       schemaPhoneField(),
			"address_ref": schemaAddressRefField(),
			"file_ref":    schemaFileRefField(),
		},
	}
}

func resourceContactCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics

	uid := strings.TrimSpace(d.Get("uid").(string))
	if uid == "" {
		uid = core.GenerateUid()
	}
	if validUid := validateUid(uid); !validUid {
		return diag.Errorf("invalid UID format - use unpadded base64url encoded value (RFC 4648)")
	}

	folderUid := strings.TrimSpace(d.Get("folder_uid").(string))
	if folderUid == "" {
		return diag.Errorf("'folder_uid' is required to create new resource")
	}

	nrc := core.NewRecordCreate("contact", "")
	if title := d.Get("title"); title != nil && title.(string) != "" {
		nrc.Title = title.(string)
	}
	if notes := d.Get("notes"); notes != nil && notes.(string) != "" {
		nrc.Notes = notes.(string)
	}

	if fieldData := d.Get("name"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("name", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			nrc.Fields = append(nrc.Fields, field)
			if err := SetFieldTypeInSchema(d, "name", "name"); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if fieldData := d.Get("company"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("text", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			nrc.Fields = append(nrc.Fields, field)
			if err := SetFieldTypeInSchema(d, "company", "text"); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if fieldData := d.Get("email"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("email", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			nrc.Fields = append(nrc.Fields, field)
			if err := SetFieldTypeInSchema(d, "email", "email"); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if fieldData := d.Get("phone"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("phone", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			nrc.Fields = append(nrc.Fields, field)
			if err := SetFieldTypeInSchema(d, "phone", "phone"); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if fieldData := d.Get("address_ref"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("addressRef", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			nrc.Fields = append(nrc.Fields, field)
			if err := SetFieldTypeInSchema(d, "address_ref", "addressRef"); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if fieldData := d.Get("file_ref"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("fileRef", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			nrc.Fields = append(nrc.Fields, field)
			if err := SetFieldTypeInSchema(d, "file_ref", "fileRef"); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if folderUid == "*" {
		if fuid, err := getTemplateFolder(folderUid, client); err != nil {
			return diag.FromErr(err)
		} else {
			folderUid = fuid
		}
	}
	uid, err := client.CreateSecretWithRecordData(uid, folderUid, nrc)
	if err != nil {
		return diag.FromErr(err)
	}

	if fuid := strings.TrimSpace(d.Get("folder_uid").(string)); fuid == "*" {
		if err = d.Set("folder_uid", folderUid); err != nil {
			return diag.FromErr(err)
		}
	}
	if err = d.Set("uid", uid); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("type", "contact"); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(uid)
	return diags
}

func resourceContactRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics

	uid := strings.TrimSpace(d.Get("uid").(string))
	title := strings.TrimSpace(d.Get("title").(string))
	if uid == "" && title == "" {
		return diag.Errorf("record UID and/or title required to locate the record")
	}

	secret, err := getRecord(uid, title, client)
	if err != nil {
		if strings.HasPrefix(err.Error(), "record not found") {
			// resource does not exist in the vault
			d.SetId("")  // mark for removal
			return diags // no error
		}
		return diag.FromErr(err)
	}

	resourceType := "contact"
	recordType := secret.Type()
	if recordType != resourceType {
		return diag.Errorf("record type '%s' is not the expected type '%s' for this data source", recordType, resourceType)
	}
	if uid == "" { // found by title: uid="", path="*"
		if err = d.Set("uid", secret.Uid); err != nil {
			return diag.FromErr(err)
		}
	}
	if err = d.Set("folder_uid", secret.FolderUid()); err != nil {
		return diag.FromErr(err)
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

	name := getFieldResourceData("name", "fields", secret)
	if err = d.Set("name", name); err != nil {
		return diag.FromErr(err)
	}
	company := getFieldResourceData("text", "fields", secret)
	if err = d.Set("company", company); err != nil {
		return diag.FromErr(err)
	}
	email := getFieldResourceData("email", "fields", secret)
	if err = d.Set("email", email); err != nil {
		return diag.FromErr(err)
	}
	phone := getFieldResourceData("phone", "fields", secret)
	if err = d.Set("phone", phone); err != nil {
		return diag.FromErr(err)
	}
	addressRef := getFieldResourceData("addressRef", "fields", secret)
	if err = d.Set("address_ref", addressRef); err != nil {
		return diag.FromErr(err)
	}

	fileItems := getFileItemsResourceData(secret)
	if err := d.Set("file_ref", fileItems); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(uid)
	return diags
}

func resourceContactUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics

	// folderUid := strings.TrimSpace(d.Get("folder_uid").(string))
	uid := strings.TrimSpace(d.Get("uid").(string))
	if uid == "" {
		return diag.Errorf("'uid' is required to update existing resource")
	}

	hasRestrictedChanges := d.HasChange("folder_uid") || d.HasChange("uid") || d.HasChange("type")
	if hasRestrictedChanges {
		return diag.Errorf("changes to folder_uid, uid, and type are not allowed")
	}

	title := strings.TrimSpace(d.Get("title").(string))
	secret, err := getRecord(uid, title, client)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("title") {
		secret.SetTitle(d.Get("title").(string))
	}
	if d.HasChange("notes") {
		secret.SetNotes(d.Get("notes").(string))
	}

	if d.HasChange("name") {
		if _, err := ApplyFieldChange("fields", "name", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("company") {
		if _, err := ApplyFieldChange("fields", "company", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("email") {
		if _, err := ApplyFieldChange("fields", "email", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("phone") {
		if _, err := ApplyFieldChange("fields", "phone", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("address_ref") {
		if _, err := ApplyFieldChange("fields", "address_ref", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("file_ref") {
		if _, err := ApplyFieldChange("fields", "file_ref", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}

	secret.RawJson = core.DictToJson(secret.RecordDict)
	if err := client.Save(secret); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(uid)
	return diags
}

func resourceContactDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	// provider := m.(providerMeta)
	// client := *provider.client

	// uid := strings.TrimSpace(d.Get("uid").(string))
	// if uid == "" {
	// 	return diag.Errorf("'uid' is required to delete existing resource")
	// }

	// if err := client.DeleteSecret(uid); err != nil && !strings.HasPrefix(err.Error(), "record not found") {
	// 	return diag.FromErr(err)
	// }
	// // NB! Do not return an error if resource already deleted by the vault/app
	// // This allows users to manually delete resources without breaking Terraform.
	d.SetId("")
	return diags
}

func resourceContactImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	uid := d.Id()

	err := d.Set("uid", uid)
	if err != nil {
		return nil, err
	}

	diags := resourceContactRead(ctx, d, m)
	if diags.HasError() {
		for i := range diags {
			if diags[i].Severity == diag.Error {
				return nil, errors.New(diags[i].Summary + " *** Details: " + diags[i].Detail)
			}
		}
	}

	return []*schema.ResourceData{d}, nil
}
