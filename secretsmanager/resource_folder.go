package secretsmanager

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFolder() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFolderCreate,
		ReadContext:   resourceFolderRead,
		UpdateContext: resourceFolderUpdate,
		DeleteContext: resourceFolderDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceFolderImport,
		},
		Schema: map[string]*schema.Schema{
			"parent_uid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The parent folder UID where the folder is created.",
			},
			"uid": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Computed:    true,
				Description: "The folder UID (using RFC4648 URL and Filename Safe Alphabet).",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The folder name.",
			},
			"force_delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Force deletion of non empty folders.",
			},
		},
	}
}

func resourceFolderCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics

	parentFolderUid := strings.TrimSpace(d.Get("parent_uid").(string))
	if parentFolderUid == "" {
		return diag.Errorf("'parent_uid' is required to create new resource")
	}
	if validUid := validateUid(parentFolderUid); !validUid {
		return diag.Errorf("invalid UID format - use unpadded base64url encoded value (RFC 4648)")
	}

	folderName := d.Get("name").(string)
	if folderName == "" {
		return diag.Errorf("'name' is required to create new resource")
	}

	// folderUid := strings.TrimSpace(d.Get("uid").(string))
	folderUid, err := createFolder(parentFolderUid, folderName, client)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("uid", folderUid); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(folderUid)
	return diags
}

func resourceFolderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics

	parentFolderUid := strings.TrimSpace(d.Get("parent_uid").(string))
	if parentFolderUid == "" {
		return diag.Errorf("'parent_uid' is required to locate the sub-folder")
	}
	if validUid := validateUid(parentFolderUid); !validUid {
		return diag.Errorf("invalid UID format - use unpadded base64url encoded value (RFC 4648)")
	}

	folderUid := strings.TrimSpace(d.Get("uid").(string))
	folderName := d.Get("name").(string)
	if folderUid == "" && folderName == "" {
		return diag.Errorf("folder UID and/or name required to locate the folder")
	}

	folders, err := findSubFolder(parentFolderUid, folderUid, folderName, client)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(folders) == 0 {
		// resource does not exist in the vault
		d.SetId("")  // mark for removal
		return diags // no error
	}
	if len(folders) > 1 {
		return diag.Errorf("multilpe subfolders with same name '%s' found in parent folder UID '%s'", folderName, parentFolderUid)
	}

	folder := folders[0]
	if folder.ParentUid != parentFolderUid {
		// since folder UID is unique - we found the right folder but it changed parents
		return diag.Errorf("folder UID '%s' found but in different parent folder UID '%s'", folder.FolderUid, folder.ParentUid)
	}
	if err = d.Set("uid", folder.FolderUid); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("name", folder.Name); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(folder.FolderUid)
	return diags
}

func resourceFolderUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics

	parentFolderUid := strings.TrimSpace(d.Get("parent_uid").(string))
	if parentFolderUid == "" {
		return diag.Errorf("'parent_uid' is required to update existing resource")
	}
	folderUid := strings.TrimSpace(d.Get("uid").(string))
	if folderUid == "" {
		return diag.Errorf("'uid' is required to update existing resource")
	}

	hasRestrictedChanges := d.HasChange("parent_uid") || d.HasChange("uid")
	if hasRestrictedChanges {
		return diag.Errorf("changes to parent_uid and uid are not allowed")
	}

	folderName := strings.TrimSpace(d.Get("name").(string))
	if d.HasChange("name") {
		if err := client.UpdateFolder(folderUid, folderName, nil); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(folderUid)
	return diags
}

func resourceFolderDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics

	folderUid := strings.TrimSpace(d.Get("uid").(string))
	if folderUid == "" {
		return diag.Errorf("'uid' is required to delete existing resource")
	}

	forceDelete, ok := d.Get("force_delete").(bool)
	if !ok {
		forceDelete = false
	}
	if err := deleteFolder(folderUid, forceDelete, client); err != nil {
		if strings.HasSuffix(err.Error(), "unexpected status: ''") {
			// record UID no longer exists - probably deleted externally
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Folder UID: %s not found - probably already deleted (externally)", folderUid),
				Detail: fmt.Sprintf("Delete Folder UID: %s returned empty status."+
					" That usually means the folder doesn't exist -"+
					" either already deleted (externally),"+
					" or no longer shared to the corresponding KSM Application.", folderUid),
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

func resourceFolderImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	provider := m.(providerMeta)
	client := *provider.client

	uid := d.Id()
	err := d.Set("uid", uid)
	if err != nil {
		return nil, err
	}

	folders, err := findSubFolder("", uid, "", client)
	if err != nil {
		return nil, err
	}
	if len(folders) == 0 {
		return nil, errors.New("failed to import folder UID='" + uid + "' - folder not found.")
	}
	// if len(folders) > 1 {
	// 	return nil, errors.New("failed to import folder UID='" + uid + "' - multiple matching folders found.")
	// }

	folder := folders[0]
	parentFolderUid := folder.ParentUid
	if strings.TrimSpace(parentFolderUid) == "" {
		return nil, errors.New("cannot import root shared folder UID='" + uid + "' - directly shared to KSM App.")
	}
	err = d.Set("parent_uid", parentFolderUid)
	if err != nil {
		return nil, err
	}

	diags := resourceFolderRead(ctx, d, m)
	if diags.HasError() {
		for i := range diags {
			if diags[i].Severity == diag.Error {
				return nil, errors.New(diags[i].Summary + " *** Details: " + diags[i].Detail)
			}
		}
	}

	return []*schema.ResourceData{d}, nil
}
