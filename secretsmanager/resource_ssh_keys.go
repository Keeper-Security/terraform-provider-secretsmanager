package secretsmanager

import (
	"context"
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/keeper-security/secrets-manager-go/core"
)

func resourceSshKeys() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSshKeysCreate,
		ReadContext:   resourceSshKeysRead,
		UpdateContext: resourceSshKeysUpdate,
		DeleteContext: resourceSshKeysDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSshKeysImport,
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
			"login":      schemaLoginField(),
			"key_pair":   schemaKeyPairField(),
			"passphrase": schemaPasswordField("passphrase"),
			"host":       schemaHostField(),
			"file_ref":   schemaFileRefField(),
		},
	}
}

func resourceSshKeysCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	nrc := core.NewRecordCreate("sshKeys", "")
	if title := d.Get("title"); title != nil && title.(string) != "" {
		nrc.Title = title.(string)
	}
	if notes := d.Get("notes"); notes != nil && notes.(string) != "" {
		nrc.Notes = notes.(string)
	}

	if fieldData := d.Get("login"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("login", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			nrc.Fields = append(nrc.Fields, field)
			if err := SetFieldTypeInSchema(d, "login", "login"); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if fieldData := d.Get("passphrase"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("password", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			if generated, err := applyGeneratePassword(fieldData, field); err != nil {
				return diag.FromErr(err)
			} else if generated {
				if err := d.Set("passphrase", fieldData); err != nil {
					return diag.FromErr(err)
				}
			}
			nrc.Fields = append(nrc.Fields, field)
			if err := SetFieldTypeInSchema(d, "passphrase", "password"); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if fieldData := d.Get("host"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("host", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			nrc.Fields = append(nrc.Fields, field)
			if err := SetFieldTypeInSchema(d, "host", "host"); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if fieldData := d.Get("key_pair"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("keyPair", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			nrc.Fields = append(nrc.Fields, field)
			if err := SetFieldTypeInSchema(d, "key_pair", "keyPair"); err != nil {
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
	if err = d.Set("type", "sshKeys"); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(uid)
	return diags
}

func resourceSshKeysRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	resourceType := "sshKeys"
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

	login := getFieldResourceData("login", "fields", secret)
	if err = d.Set("login", login); err != nil {
		return diag.FromErr(err)
	}
	password := getFieldResourceData("password", "fields", secret)
	mergePassword(d.Get("passphrase"), password)
	if err = d.Set("passphrase", password); err != nil {
		return diag.FromErr(err)
	}
	host := getFieldResourceData("host", "fields", secret)
	if err = d.Set("host", host); err != nil {
		return diag.FromErr(err)
	}
	keyPair := getFieldResourceData("keyPair", "fields", secret)
	if err = d.Set("key_pair", keyPair); err != nil {
		return diag.FromErr(err)
	}

	fileItems := getFileItemsResourceData(secret)
	if err := d.Set("file_ref", fileItems); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(uid)
	return diags
}

func resourceSshKeysUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	if d.HasChange("login") {
		if _, err := ApplyFieldChange("fields", "login", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("passphrase") {
		if _, err := ApplyFieldChange("fields", "passphrase", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("host") {
		if _, err := ApplyFieldChange("fields", "host", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("key_pair") {
		if _, err := ApplyFieldChange("fields", "key_pair", d, secret); err != nil {
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

func resourceSshKeysDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func resourceSshKeysImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	uid := d.Id()

	err := d.Set("uid", uid)
	if err != nil {
		return nil, err
	}

	diags := resourceSshKeysRead(ctx, d, m)
	if diags.HasError() {
		for i := range diags {
			if diags[i].Severity == diag.Error {
				return nil, errors.New(diags[i].Summary + " *** Details: " + diags[i].Detail)
			}
		}
	}

	return []*schema.ResourceData{d}, nil
}
