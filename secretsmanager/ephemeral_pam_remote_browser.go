package secretsmanager

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ ephemeral.EphemeralResource              = &ephemeralPamRemoteBrowser{}
	_ ephemeral.EphemeralResourceWithConfigure = &ephemeralPamRemoteBrowser{}
)

type ephemeralPamRemoteBrowser struct {
	meta providerMeta
}

type ephemeralPamRemoteBrowserModel struct {
	Path                      types.String `tfsdk:"path"`
	Type                      types.String `tfsdk:"type"`
	Title                     types.String `tfsdk:"title"`
	Notes                     types.String `tfsdk:"notes"`
	FolderUID                 types.String `tfsdk:"folder_uid"`
	RbiUrl                    types.String `tfsdk:"rbi_url"`
	PamRemoteBrowserSettings  types.String `tfsdk:"pam_remote_browser_settings"`
	TrafficEncryptionSeed     types.String `tfsdk:"traffic_encryption_seed"`
	FileRef                   types.List   `tfsdk:"file_ref"`
	TOTP                      types.List   `tfsdk:"totp"`
	Custom  types.List   `tfsdk:"custom"`
}

func NewEphemeralPamRemoteBrowser() ephemeral.EphemeralResource {
	return &ephemeralPamRemoteBrowser{}
}

func (e *ephemeralPamRemoteBrowser) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pam_remote_browser"
}

func (e *ephemeralPamRemoteBrowser) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this ephemeral resource to read a PAM Remote Browser record from Keeper Secrets Manager. Values are never stored in state.",
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				Required:    true,
				Description: "The path where the secret is stored.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The secret type.",
			},
			"title": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The secret title.",
			},
			"notes": schema.StringAttribute{
				Computed:    true,
				Description: "The secret notes.",
			},
			"folder_uid": schema.StringAttribute{
				Computed:    true,
				Description: "The folder UID where the secret is stored.",
			},
			"rbi_url": schema.StringAttribute{
				Computed:    true,
				Description: "The Remote Browser Interface URL.",
			},
			"pam_remote_browser_settings": schema.StringAttribute{
				Computed:    true,
				Description: "PAM Remote Browser connection settings as a JSON string.",
			},
			"traffic_encryption_seed": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "Base64 encoded traffic encryption seed.",
			},
			"file_ref": fileRefEphemeralAttribute(),
			"totp": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The one time password.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"url": schema.StringAttribute{
							Computed:    true,
							Description: "TOTP URL.",
						},
						"token": schema.StringAttribute{
							Computed:    true,
							Sensitive:   true,
							Description: "Generated TOTP token.",
						},
						"ttl": schema.Int64Attribute{
							Computed:    true,
							Description: "Time to live for TOTP token in seconds.",
						},
					},
				},
			},
			"custom": genericFieldEphemeralAttribute("Custom fields of the record."),
		},
	}
}

func (e *ephemeralPamRemoteBrowser) Configure(_ context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta, ok := req.ProviderData.(providerMeta)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data", "Expected providerMeta")
		return
	}
	e.meta = meta
}

func (e *ephemeralPamRemoteBrowser) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data ephemeralPamRemoteBrowserModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if e.meta.client == nil {
		resp.Diagnostics.AddError("Provider Not Configured", "KSM client is not configured. Ensure the provider credential is set.")
		return
	}

	client := *e.meta.client
	path := strings.TrimSpace(data.Path.ValueString())
	title := ""
	if !data.Title.IsNull() && !data.Title.IsUnknown() {
		title = strings.TrimSpace(data.Title.ValueString())
	}

	secret, err := getRecord(path, title, client)
	if err != nil {
		resp.Diagnostics.AddError("Error reading secret", err.Error())
		return
	}

	recordType := secret.Type()
	if recordType != "pamRemoteBrowser" {
		resp.Diagnostics.AddError("Record Type Mismatch",
			"record type '"+recordType+"' is not the expected type 'pamRemoteBrowser' for this ephemeral resource")
		return
	}

	data.Type = types.StringValue(recordType)
	data.Title = types.StringValue(secret.Title())
	data.Notes = types.StringValue(secret.Notes())

	fuid := secret.InnerFolderUid()
	if fuid == "" {
		fuid = secret.FolderUid()
	}
	data.FolderUID = types.StringValue(fuid)

	// rbi_url: extract from custom field type "rbiUrl"
	data.RbiUrl = types.StringValue(pamFieldStringByType("rbiUrl", secret))

	// pam_remote_browser_settings: convert pamRemoteBrowserSettings field to JSON string
	pamRbsFields := secret.GetFieldsByType("pamRemoteBrowserSettings")
	if len(pamRbsFields) > 0 {
		pamRbsJSON, err := pamSettingsFieldToJSON(pamRbsFields[0])
		if err == nil {
			data.PamRemoteBrowserSettings = types.StringValue(pamRbsJSON)
		} else {
			data.PamRemoteBrowserSettings = types.StringValue("")
		}
	} else {
		data.PamRemoteBrowserSettings = types.StringValue("")
	}

	// traffic_encryption_seed: extract from custom field type "trafficEncryptionSeed"
	data.TrafficEncryptionSeed = types.StringValue(pamFieldStringByType("trafficEncryptionSeed", secret))

	totpUrl := strings.TrimSpace(pamFieldString("oneTimeCode", "fields", secret))
	totpList, diags := totpToListValue(ctx, totpUrl)
	resp.Diagnostics.Append(diags...)
	data.TOTP = totpList

	fileRefList, diags := fileItemsToListValue(ctx, secret.Files)
	resp.Diagnostics.Append(diags...)
	data.FileRef = fileRefList

	customItems := getFieldItemsData(secret.RecordDict, "custom")
	customList, diags := genericFieldItemsToListValue(ctx, customItems)
	resp.Diagnostics.Append(diags...)
	data.Custom = customList


	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
