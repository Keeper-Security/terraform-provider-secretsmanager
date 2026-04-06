package secretsmanager

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ ephemeral.EphemeralResource              = &ephemeralPamDatabase{}
	_ ephemeral.EphemeralResourceWithConfigure = &ephemeralPamDatabase{}
)

type ephemeralPamDatabase struct {
	meta providerMeta
}

type ephemeralPamDatabaseModel struct {
	Path           types.String `tfsdk:"path"`
	Type           types.String `tfsdk:"type"`
	Title          types.String `tfsdk:"title"`
	Notes          types.String `tfsdk:"notes"`
	FolderUID      types.String `tfsdk:"folder_uid"`
	PamHostname    types.List   `tfsdk:"pam_hostname"`
	PamSettings    types.String `tfsdk:"pam_settings"`
	UseSSL         types.Bool   `tfsdk:"use_ssl"`
	DatabaseId     types.String `tfsdk:"database_id"`
	DatabaseType   types.String `tfsdk:"database_type"`
	ProviderGroup  types.String `tfsdk:"provider_group"`
	ProviderRegion types.String `tfsdk:"provider_region"`
	FileRef        types.List   `tfsdk:"file_ref"`
	TOTP           types.List   `tfsdk:"totp"`
}

func NewEphemeralPamDatabase() ephemeral.EphemeralResource {
	return &ephemeralPamDatabase{}
}

func (e *ephemeralPamDatabase) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pam_database"
}

func (e *ephemeralPamDatabase) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this ephemeral resource to read a PAM Database record from Keeper Secrets Manager. Values are never stored in state.",
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
			"use_ssl": schema.BoolAttribute{
				Computed:    true,
				Description: "Use SSL.",
			},
			"database_id": schema.StringAttribute{
				Computed:    true,
				Description: "Database Id.",
			},
			"database_type": schema.StringAttribute{
				Computed:    true,
				Description: "Database type (e.g. postgresql, mysql, mongodb).",
			},
			"provider_group": schema.StringAttribute{
				Computed:    true,
				Description: "Provider Group.",
			},
			"provider_region": schema.StringAttribute{
				Computed:    true,
				Description: "Provider Region.",
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

func (e *ephemeralPamDatabase) Configure(_ context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
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

func (e *ephemeralPamDatabase) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data ephemeralPamDatabaseModel
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
	if recordType != "pamDatabase" {
		resp.Diagnostics.AddError("Record Type Mismatch",
			"record type '"+recordType+"' is not the expected type 'pamDatabase' for this ephemeral resource")
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
	data.UseSSL = types.BoolValue(pamFieldBoolWithLabel("checkbox", "fields", secret, "useSSL"))
	data.DatabaseId = types.StringValue(pamFieldStringWithLabel("text", "fields", secret, "Database Id"))
	data.DatabaseType = types.StringValue(pamFieldStringByType("databaseType", secret))
	data.ProviderGroup = types.StringValue(pamFieldStringWithLabel("text", "fields", secret, "Provider Group"))
	data.ProviderRegion = types.StringValue(pamFieldStringWithLabel("text", "fields", secret, "Provider Region"))

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
