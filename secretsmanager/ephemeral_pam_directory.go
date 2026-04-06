package secretsmanager

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ ephemeral.EphemeralResource              = &ephemeralPamDirectory{}
	_ ephemeral.EphemeralResourceWithConfigure = &ephemeralPamDirectory{}
)

type ephemeralPamDirectory struct {
	meta providerMeta
}

type ephemeralPamDirectoryModel struct {
	Path              types.String `tfsdk:"path"`
	Type              types.String `tfsdk:"type"`
	Title             types.String `tfsdk:"title"`
	Notes             types.String `tfsdk:"notes"`
	FolderUID         types.String `tfsdk:"folder_uid"`
	PamHostname       types.List   `tfsdk:"pam_hostname"`
	PamSettings       types.String `tfsdk:"pam_settings"`
	DirectoryType     types.String `tfsdk:"directory_type"`
	UseSSL            types.Bool   `tfsdk:"use_ssl"`
	DistinguishedName types.String `tfsdk:"distinguished_name"`
	DomainName        types.String `tfsdk:"domain_name"`
	DirectoryId       types.String `tfsdk:"directory_id"`
	UserMatch         types.String `tfsdk:"user_match"`
	ProviderGroup     types.String `tfsdk:"provider_group"`
	ProviderRegion    types.String `tfsdk:"provider_region"`
	AlternativeIPs    types.String `tfsdk:"alternative_ips"`
	FileRef           types.List   `tfsdk:"file_ref"`
	TOTP              types.List   `tfsdk:"totp"`
}

func NewEphemeralPamDirectory() ephemeral.EphemeralResource {
	return &ephemeralPamDirectory{}
}

func (e *ephemeralPamDirectory) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pam_directory"
}

func (e *ephemeralPamDirectory) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this ephemeral resource to read a PAM Directory record from Keeper Secrets Manager. Values are never stored in state.",
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
			"pam_hostname": hostEphemeralAttribute(),
			"pam_settings": schema.StringAttribute{
				Computed:    true,
				Description: "PAM connection settings as a JSON string.",
			},
			"directory_type": schema.StringAttribute{
				Computed:    true,
				Description: "Directory type.",
			},
			"use_ssl": schema.BoolAttribute{
				Computed:    true,
				Description: "Use SSL.",
			},
			"distinguished_name": schema.StringAttribute{
				Computed:    true,
				Description: "Distinguished Name.",
			},
			"domain_name": schema.StringAttribute{
				Computed:    true,
				Description: "Domain Name.",
			},
			"directory_id": schema.StringAttribute{
				Computed:    true,
				Description: "Directory Id.",
			},
			"user_match": schema.StringAttribute{
				Computed:    true,
				Description: "User Match.",
			},
			"provider_group": schema.StringAttribute{
				Computed:    true,
				Description: "Provider Group.",
			},
			"provider_region": schema.StringAttribute{
				Computed:    true,
				Description: "Provider Region.",
			},
			"alternative_ips": schema.StringAttribute{
				Computed:    true,
				Description: "Alternative IPs (multiline).",
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
		},
	}
}

func (e *ephemeralPamDirectory) Configure(_ context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
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

func (e *ephemeralPamDirectory) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data ephemeralPamDirectoryModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
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
	if recordType != "pamDirectory" {
		resp.Diagnostics.AddError("Record Type Mismatch",
			"record type '"+recordType+"' is not the expected type 'pamDirectory' for this ephemeral resource")
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

	pamHostnameList, diags := pamHostnameToListValue(ctx, secret)
	resp.Diagnostics.Append(diags...)
	data.PamHostname = pamHostnameList

	data.PamSettings = types.StringValue(pamSettingsToString(secret))
	data.DirectoryType = types.StringValue(pamFieldStringByType("directoryType", secret))
	data.UseSSL = types.BoolValue(pamFieldBoolWithLabel("checkbox", "fields", secret, "useSSL"))
	data.DistinguishedName = types.StringValue(pamFieldStringWithLabel("text", "fields", secret, "Distinguished Name"))
	data.DomainName = types.StringValue(pamFieldStringWithLabel("text", "fields", secret, "domainName"))
	data.DirectoryId = types.StringValue(pamFieldStringWithLabel("text", "fields", secret, "directoryId"))
	data.UserMatch = types.StringValue(pamFieldStringWithLabel("text", "fields", secret, "userMatch"))
	data.ProviderGroup = types.StringValue(pamFieldStringWithLabel("text", "fields", secret, "providerGroup"))
	data.ProviderRegion = types.StringValue(pamFieldStringWithLabel("text", "fields", secret, "providerRegion"))
	data.AlternativeIPs = types.StringValue(pamFieldStringWithLabel("multiline", "fields", secret, "alternativeIPs"))

	totpUrl := strings.TrimSpace(pamFieldString("oneTimeCode", "fields", secret))
	totpList, diags := totpToListValue(ctx, totpUrl)
	resp.Diagnostics.Append(diags...)
	data.TOTP = totpList

	fileRefList, diags := fileItemsToListValue(ctx, secret.Files)
	resp.Diagnostics.Append(diags...)
	data.FileRef = fileRefList

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
