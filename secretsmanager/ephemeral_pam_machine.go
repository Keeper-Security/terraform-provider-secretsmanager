package secretsmanager

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ ephemeral.EphemeralResource              = &ephemeralPamMachine{}
	_ ephemeral.EphemeralResourceWithConfigure = &ephemeralPamMachine{}
)

type ephemeralPamMachine struct {
	meta providerMeta
}

type ephemeralPamMachineModel struct {
	Path                 types.String `tfsdk:"path"`
	Type                 types.String `tfsdk:"type"`
	Title                types.String `tfsdk:"title"`
	Notes                types.String `tfsdk:"notes"`
	FolderUID            types.String `tfsdk:"folder_uid"`
	PamHostname          types.List   `tfsdk:"pam_hostname"`
	PamSettings          types.String `tfsdk:"pam_settings"`
	Login                types.String `tfsdk:"login"`
	Password             types.String `tfsdk:"password"`
	PrivatePemKey        types.String `tfsdk:"private_pem_key"`
	PrivateKeyPassphrase types.String `tfsdk:"private_key_passphrase"`
	OperatingSystem      types.String `tfsdk:"operating_system"`
	SslVerification      types.Bool   `tfsdk:"ssl_verification"`
	InstanceName         types.String `tfsdk:"instance_name"`
	InstanceId           types.String `tfsdk:"instance_id"`
	ProviderGroup        types.String `tfsdk:"provider_group"`
	ProviderRegion       types.String `tfsdk:"provider_region"`
	FileRef              types.List   `tfsdk:"file_ref"`
	TOTP                 types.List   `tfsdk:"totp"`
	Custom  types.List   `tfsdk:"custom"`
}

func NewEphemeralPamMachine() ephemeral.EphemeralResource {
	return &ephemeralPamMachine{}
}

func (e *ephemeralPamMachine) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pam_machine"
}

func (e *ephemeralPamMachine) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this ephemeral resource to read a PAM Machine record from Keeper Secrets Manager. Values are never stored in state.",
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
			"login": schema.StringAttribute{
				Computed:    true,
				Description: "The login username.",
			},
			"password": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The password.",
			},
			"private_pem_key": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The private PEM key.",
			},
			"private_key_passphrase": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The private key passphrase.",
			},
			"operating_system": schema.StringAttribute{
				Computed:    true,
				Description: "Operating System.",
			},
			"ssl_verification": schema.BoolAttribute{
				Computed:    true,
				Description: "SSL Verification.",
			},
			"instance_name": schema.StringAttribute{
				Computed:    true,
				Description: "Instance Name.",
			},
			"instance_id": schema.StringAttribute{
				Computed:    true,
				Description: "Instance Id.",
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
			"custom": genericFieldEphemeralAttribute("Custom fields of the record."),
		},
	}
}

func (e *ephemeralPamMachine) Configure(_ context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
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

func (e *ephemeralPamMachine) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data ephemeralPamMachineModel
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
	if recordType != "pamMachine" {
		resp.Diagnostics.AddError("Record Type Mismatch",
			"record type '"+recordType+"' is not the expected type 'pamMachine' for this ephemeral resource")
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
	data.Login = types.StringValue(pamFieldString("login", "fields", secret))
	data.Password = types.StringValue(pamFieldString("password", "fields", secret))
	data.PrivatePemKey = types.StringValue(pamFieldStringWithLabel("secret", "fields", secret, "Private PEM Key"))
	data.PrivateKeyPassphrase = types.StringValue(pamFieldStringWithLabel("secret", "custom", secret, "Private Key Passphrase"))
	data.OperatingSystem = types.StringValue(pamFieldStringWithLabel("text", "fields", secret, "Operating System"))
	data.SslVerification = types.BoolValue(pamFieldBoolWithLabel("checkbox", "fields", secret, "SSL Verification"))
	data.InstanceName = types.StringValue(pamFieldStringWithLabel("text", "fields", secret, "Instance Name"))
	data.InstanceId = types.StringValue(pamFieldStringWithLabel("text", "fields", secret, "Instance Id"))
	data.ProviderGroup = types.StringValue(pamFieldStringWithLabel("text", "fields", secret, "Provider Group"))
	data.ProviderRegion = types.StringValue(pamFieldStringWithLabel("text", "fields", secret, "Provider Region"))

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
