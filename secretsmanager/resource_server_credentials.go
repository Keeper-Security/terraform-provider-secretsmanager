package secretsmanager

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/keeper-security/secrets-manager-go/core"
)

func resourceServerCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServerCredentialsCreate,
		ReadContext:   resourceServerCredentialsRead,
		UpdateContext: resourceServerCredentialsUpdate,
		DeleteContext: resourceServerCredentialsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceServerCredentialsImport,
		},
		Schema: map[string]*schema.Schema{
			"folder_uid": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				AtLeastOneOf: []string{"folder_uid", "uid"},
				Description:  "The folder UID where the secret is stored. The parent shared folder must be non empty.",
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
			"host":     schemaHostField(),
			"login":    schemaLoginField(),
			"password": schemaPasswordField(""),
			"file_ref": schemaFileRefField(),
		},
	}
}

func resourceServerCredentialsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	nrc := core.NewRecordCreate("serverCredentials", "")
	if title := d.Get("title"); title != nil && title.(string) != "" {
		nrc.Title = title.(string)
	}
	if notes := d.Get("notes"); notes != nil && notes.(string) != "" {
		nrc.Notes = notes.(string)
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
	if fieldData := d.Get("password"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("password", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			if generated, err := applyGeneratePassword(fieldData, field); err != nil {
				return diag.FromErr(err)
			} else if generated {
				if err := d.Set("password", fieldData); err != nil {
					return diag.FromErr(err)
				}
			}
			nrc.Fields = append(nrc.Fields, field)
			if err := SetFieldTypeInSchema(d, "password", "password"); err != nil {
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
	uid, err := createRecord(uid, folderUid, nrc, client)
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
	if err = d.Set("type", "serverCredentials"); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(uid)
	return diags
}

func resourceServerCredentialsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	resourceType := "serverCredentials"
	recordType := secret.Type()
	if recordType != resourceType {
		return diag.Errorf("record type '%s' is not the expected type '%s' for this data source", recordType, resourceType)
	}
	if uid == "" { // found by title: uid="", path="*"
		if err = d.Set("uid", secret.Uid); err != nil {
			return diag.FromErr(err)
		}
	}
	fuid := secret.InnerFolderUid() // in subfolder
	if fuid == "" {                 // directly in shared folder
		fuid = secret.FolderUid()
	}
	if fuid != "" {
		if err = d.Set("folder_uid", fuid); err != nil {
			return diag.FromErr(err)
		}
	} // else - directly shared to the KSM App (not through shared folder)
	if err = d.Set("type", recordType); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("title", secret.Title()); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("notes", secret.Notes()); err != nil {
		return diag.FromErr(err)
	}

	host := getFieldResourceData("host", "fields", secret)
	if err = d.Set("host", host); err != nil {
		return diag.FromErr(err)
	}
	login := getFieldResourceData("login", "fields", secret)
	if err = d.Set("login", login); err != nil {
		return diag.FromErr(err)
	}
	password := getFieldResourceData("password", "fields", secret)
	mergePassword(d.Get("password"), password)
	if err = d.Set("password", password); err != nil {
		return diag.FromErr(err)
	}

	fileItems := getFileItemsResourceData(secret)
	if err := d.Set("file_ref", fileItems); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(uid)
	return diags
}

func resourceServerCredentialsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	if d.HasChange("host") {
		if _, err := ApplyFieldChange("fields", "host", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("login") {
		if _, err := ApplyFieldChange("fields", "login", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("password") {
		if _, err := ApplyFieldChange("fields", "password", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("file_ref") {
		if _, err := ApplyFieldChange("fields", "file_ref", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}

	secret.RawJson = core.DictToJson(secret.RecordDict)
	if err := saveRecord(secret, client); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(uid)
	return diags
}

func resourceServerCredentialsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics

	uid := strings.TrimSpace(d.Get("uid").(string))
	if uid == "" {
		return diag.Errorf("'uid' is required to delete existing resource")
	}

	if err := deleteRecord(uid, client); err != nil {
		if strings.HasSuffix(err.Error(), "unexpected status: ''") {
			// record UID no longer exists - probably deleted externally
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Record UID: %s not found - probably already deleted (externally)", uid),
				Detail: fmt.Sprintf("Delete record UID: %s returned empty status."+
					" That usually means the record doesn't exist -"+
					" either already deleted (externally),"+
					" or no longer shared to the corresponding KSM Application.", uid),
			})
		} else {
			return diag.FromErr(err)
		}
	}
	// NB! Do not return an error if resource already deleted by the vault/app
	// This allows users to manually delete resources without breaking Terraform.
	d.SetId("")
	return diags
}

func resourceServerCredentialsImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	uid := d.Id()

	err := d.Set("uid", uid)
	if err != nil {
		return nil, err
	}

	diags := resourceServerCredentialsRead(ctx, d, m)
	if diags.HasError() {
		for i := range diags {
			if diags[i].Severity == diag.Error {
				return nil, errors.New(diags[i].Summary + " *** Details: " + diags[i].Detail)
			}
		}
	}

	return []*schema.ResourceData{d}, nil
}
