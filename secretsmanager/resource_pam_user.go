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

func resourcePamUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePamUserCreate,
		ReadContext:   resourcePamUserRead,
		UpdateContext: resourcePamUserUpdate,
		DeleteContext: resourcePamUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePamUserImport,
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
			// PAM User specific fields
			"login":                  schemaLoginField(),
			"password":               schemaPasswordField(""),
			"rotation_scripts":       schemaScriptField(),
			"private_pem_key":        schemaPrivatePemKeyField(),
			"private_key_passphrase": schemaPrivateKeyPassphraseField(),
			"distinguished_name":     schemaTextField(),
			"connect_database":       schemaTextField(),
			"managed":                schemaCheckboxField(),
			"file_ref":               schemaFileRefField(),
			"totp":                   schemaOneTimeCodeField(),
			// custom[]
			"custom": schemaCustomField(),
		},
	}
}

func resourcePamUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	nrc := core.NewRecordCreate("pamUser", "")
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

	// Process private key passphrase — stored as custom field (type: secret, label: "Private Key Passphrase")
	var passphraseValue string
	if fieldData := d.Get("private_key_passphrase"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		// Use password generation infrastructure for passphrase
		if field, err := NewFieldFromSchema("password", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			if generated, err := applyGeneratePassword(fieldData, field); err != nil {
				return diag.FromErr(err)
			} else if generated {
				if err := d.Set("private_key_passphrase", fieldData); err != nil {
					return diag.FromErr(err)
				}
			}
			if pwField, ok := field.(*core.Password); ok && len(pwField.Value) > 0 {
				passphraseValue = pwField.Value[0]
			}
			// Store as custom field
			secretField := core.NewSecret(passphraseValue)
			secretField.Label = "Private Key Passphrase"
			nrc.Custom = append(nrc.Custom, secretField)
		}
	}

	// Process private PEM key — standard secret field, may generate using passphrase
	if fieldData := d.Get("private_pem_key"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		privatePEM, _, err := applyGeneratePamKey(fieldData, passphraseValue)
		if err != nil {
			return diag.FromErr(err)
		}
		if privatePEM != "" {
			if err := d.Set("private_pem_key", fieldData); err != nil {
				return diag.FromErr(err)
			}
		}
		// Get value (generated or manual)
		pemValue := privatePEM
		if pemValue == "" {
			if s, ok := fieldData.([]interface{}); ok && len(s) > 0 {
				if m, ok := s[0].(map[string]interface{}); ok {
					if v, ok := m["value"].(string); ok {
						pemValue = v
					}
				}
			}
		}
		if pemValue != "" {
			secretField := core.NewSecret(pemValue)
			secretField.Label = "Private PEM Key"
			nrc.Fields = append(nrc.Fields, secretField)
		}
	}

	if fieldData := d.Get("rotation_scripts"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("script", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			field.(*core.Scripts).Label = "Rotation Scripts"
			nrc.Fields = append(nrc.Fields, field)
			if err := SetFieldTypeInSchema(d, "rotation_scripts", "script"); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if fieldData := d.Get("distinguished_name"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("text", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			field.(*core.Text).Label = "Distinguished Name"
			nrc.Fields = append(nrc.Fields, field)
			if err := SetFieldTypeInSchema(d, "distinguished_name", "text"); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if fieldData := d.Get("connect_database"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("text", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			field.(*core.Text).Label = "Connect Database"
			nrc.Fields = append(nrc.Fields, field)
			if err := SetFieldTypeInSchema(d, "connect_database", "text"); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if fieldData := d.Get("managed"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("checkbox", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			field.(*core.Checkbox).Label = "Managed"
			nrc.Fields = append(nrc.Fields, field)
			if err := SetFieldTypeInSchema(d, "managed", "checkbox"); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if fieldData := d.Get("totp"); fieldData != nil && len(fieldData.([]interface{})) > 0 {
		if field, err := NewFieldFromSchema("oneTimeCode", fieldData); err != nil {
			return diag.FromErr(err)
		} else if field != nil {
			nrc.Fields = append(nrc.Fields, field)
			if err := SetFieldTypeInSchema(d, "totp", "oneTimeCode"); err != nil {
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

	// User-defined custom fields — appended after the platform-managed
	// "Private Key Passphrase" entry already placed in nrc.Custom above.
	if customData := d.Get("custom"); customData != nil {
		if fields, err := customFieldsFromSchema(customData.([]interface{})); err != nil {
			return diag.FromErr(err)
		} else {
			nrc.Custom = append(nrc.Custom, fields...)
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
	if err = d.Set("type", "pamUser"); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(uid)
	return diags
}

func resourcePamUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}

	resourceType := "pamUser"
	recordType := secret.Type()
	if recordType != resourceType {
		return diag.Errorf("record type '%s' is not the expected type '%s' for this resource", recordType, resourceType)
	}
	if uid == "" {
		if err = d.Set("uid", secret.Uid); err != nil {
			return diag.FromErr(err)
		}
	}
	fuid := secret.InnerFolderUid()
	if fuid == "" {
		fuid = secret.FolderUid()
	}
	if fuid != "" {
		if err = d.Set("folder_uid", fuid); err != nil {
			return diag.FromErr(err)
		}
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
	mergePassword(d.Get("password"), password)
	if err = d.Set("password", password); err != nil {
		return diag.FromErr(err)
	}
	oneTimeCode := getFieldResourceData("oneTimeCode", "fields", secret)
	if err = d.Set("totp", oneTimeCode); err != nil {
		return diag.FromErr(err)
	}

	// Private PEM Key (standard secret field)
	privatePemKey := getFieldResourceDataWithLabel("secret", "fields", secret, "Private PEM Key")
	mergePamKeyField(d.Get("private_pem_key"), privatePemKey)
	if err = d.Set("private_pem_key", privatePemKey); err != nil {
		return diag.FromErr(err)
	}

	// Private Key Passphrase (custom secret field)
	passphrase := getFieldResourceDataWithLabel("secret", "custom", secret, "Private Key Passphrase")
	mergePamPassphrase(d.Get("private_key_passphrase"), passphrase)
	if err = d.Set("private_key_passphrase", passphrase); err != nil {
		return diag.FromErr(err)
	}

	// PAM-specific fields
	rotationScripts := getFieldResourceDataWithLabel("script", "fields", secret, "Rotation Scripts")
	if err = d.Set("rotation_scripts", rotationScripts); err != nil {
		return diag.FromErr(err)
	}
	distinguishedName := getFieldResourceDataWithLabel("text", "fields", secret, "Distinguished Name")
	if err = d.Set("distinguished_name", distinguishedName); err != nil {
		return diag.FromErr(err)
	}
	connectDatabase := getFieldResourceDataWithLabel("text", "fields", secret, "Connect Database")
	if err = d.Set("connect_database", connectDatabase); err != nil {
		return diag.FromErr(err)
	}
	managed := getFieldResourceDataWithLabel("checkbox", "fields", secret, "Managed")
	if err = d.Set("managed", managed); err != nil {
		return diag.FromErr(err)
	}

	fileItems := getFileItemsResourceData(secret)
	if err := d.Set("file_ref", fileItems); err != nil {
		return diag.FromErr(err)
	}

	// "Private Key Passphrase" is a platform-managed entry in the custom section.
	// Exclude it from user-visible state to prevent a perpetual diff.
	allCustom := getFieldItemsData(secret.RecordDict, "custom")
	userCustom := make([]interface{}, 0, len(allCustom))
	for _, item := range allCustom {
		if m, ok := item.(map[string]interface{}); ok {
			if label, _ := m["label"].(string); label == "Private Key Passphrase" {
				continue
			}
		}
		userCustom = append(userCustom, item)
	}
	if err := d.Set("custom", userCustom); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(uid)
	return diags
}

func resourcePamUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics

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
	if d.HasChange("password") {
		if _, err := ApplyFieldChange("fields", "password", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}
	// Handle private key passphrase (custom field)
	if d.HasChange("private_key_passphrase") {
		if _, err := ApplyFieldChange("custom", "private_key_passphrase", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}
	// Handle private PEM key (standard field with potential regeneration)
	if d.HasChange("private_pem_key") {
		fieldData := d.Get("private_pem_key")
		// Extract passphrase for potential key encryption
		var passphrase string
		if ppData := d.Get("private_key_passphrase"); ppData != nil {
			if s, ok := ppData.([]interface{}); ok && len(s) > 0 {
				if m, ok := s[0].(map[string]interface{}); ok {
					if v, ok := m["value"].(string); ok {
						passphrase = v
					}
				}
			}
		}
		privatePEM, _, err := applyGeneratePamKey(fieldData, passphrase)
		if err != nil {
			return diag.FromErr(err)
		}
		if privatePEM != "" {
			if err := d.Set("private_pem_key", fieldData); err != nil {
				return diag.FromErr(err)
			}
		}
		pemValue := privatePEM
		if pemValue == "" {
			if s, ok := fieldData.([]interface{}); ok && len(s) > 0 {
				if m, ok := s[0].(map[string]interface{}); ok {
					if v, ok := m["value"].(string); ok {
						pemValue = v
					}
				}
			}
		}
		secretField := core.NewSecret(pemValue)
		secretField.Label = "Private PEM Key"
		if secret.FieldExists("fields", "secret") {
			if err := secret.UpdateField("fields", secretField); err != nil {
				return diag.FromErr(err)
			}
		} else {
			if err := secret.InsertField("fields", secretField); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if d.HasChange("rotation_scripts") {
		if _, err := ApplyFieldChange("fields", "rotation_scripts", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("distinguished_name") {
		if _, err := ApplyFieldChange("fields", "distinguished_name", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("connect_database") {
		if _, err := ApplyFieldChange("fields", "connect_database", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("managed") {
		if _, err := ApplyFieldChange("fields", "managed", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("totp") {
		if _, err := ApplyFieldChange("fields", "totp", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("file_ref") {
		if _, err := ApplyFieldChange("fields", "file_ref", d, secret); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("custom") {
		customData := d.Get("custom").([]interface{})
		userFields, err := customFieldsFromSchema(customData)
		if err != nil {
			return diag.FromErr(err)
		}
		// Preserve "Private Key Passphrase" from the current vault state;
		// replace all other custom entries with the user-defined fields.
		var preserved []interface{}
		if current, ok := secret.RecordDict["custom"].([]interface{}); ok {
			for _, item := range current {
				if m, ok := item.(map[string]interface{}); ok {
					if label, _ := m["label"].(string); label == "Private Key Passphrase" {
						preserved = append(preserved, item)
					}
				}
			}
		}
		secret.RecordDict["custom"] = append(preserved, customFieldsToDict(userFields)...)
	}

	secret.RawJson = core.DictToJson(secret.RecordDict)
	if err := saveRecord(secret, client); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(uid)
	return diags
}

func resourcePamUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics

	uid := strings.TrimSpace(d.Get("uid").(string))
	if uid == "" {
		return diag.Errorf("'uid' is required to delete existing resource")
	}

	if err := deleteRecord(uid, client); err != nil {
		if strings.HasSuffix(err.Error(), "unexpected status: ''") {
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

	d.SetId("")
	return diags
}

func resourcePamUserImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	uid := d.Id()
	if strings.TrimSpace(uid) == "" {
		return nil, errors.New("'uid' is required to import resource")
	}

	if err := d.Set("uid", uid); err != nil {
		return nil, err
	}

	diags := resourcePamUserRead(ctx, d, m)
	if diags.HasError() {
		for _, d := range diags {
			if d.Severity == diag.Error {
				return nil, fmt.Errorf("error reading PAM User: %s", d.Summary)
			}
		}
	}

	return []*schema.ResourceData{d}, nil
}
