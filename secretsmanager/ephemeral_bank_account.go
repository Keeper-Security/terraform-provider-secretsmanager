package secretsmanager

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ ephemeral.EphemeralResource              = &ephemeralBankAccount{}
	_ ephemeral.EphemeralResourceWithConfigure = &ephemeralBankAccount{}
)

type ephemeralBankAccount struct {
	meta providerMeta
}

type ephemeralBankAccountModel struct {
	Path        types.String `tfsdk:"path"`
	Type        types.String `tfsdk:"type"`
	Title       types.String `tfsdk:"title"`
	Notes       types.String `tfsdk:"notes"`
	BankAccount types.List   `tfsdk:"bank_account"`
	Name        types.List   `tfsdk:"name"`
	Login       types.String `tfsdk:"login"`
	Password    types.String `tfsdk:"password"`
	URL         types.String `tfsdk:"url"`
	CardRef     types.List   `tfsdk:"card_ref"`
	FileRef     types.List   `tfsdk:"file_ref"`
	TOTP        types.List   `tfsdk:"totp"`
}

func NewEphemeralBankAccount() ephemeral.EphemeralResource {
	return &ephemeralBankAccount{}
}

func (e *ephemeralBankAccount) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bank_account"
}

func (e *ephemeralBankAccount) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this ephemeral resource to read a bank account record from Keeper Secrets Manager. Values are never stored in state.",
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
			"bank_account": bankAccountEphemeralAttribute(),
			"name":         nameEphemeralAttribute(),
			"login": schema.StringAttribute{
				Computed:    true,
				Description: "The secret login.",
			},
			"password": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The secret password.",
			},
			"url": schema.StringAttribute{
				Computed:    true,
				Description: "The secret url.",
			},
			"card_ref": cardRefEphemeralAttribute(),
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

func (e *ephemeralBankAccount) Configure(_ context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
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

func (e *ephemeralBankAccount) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data ephemeralBankAccountModel
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
	if recordType != "bankAccount" {
		resp.Diagnostics.AddError("Record Type Mismatch",
			"record type '"+recordType+"' is not the expected type 'bankAccount' for this ephemeral resource")
		return
	}

	data.Type = types.StringValue(recordType)
	data.Title = types.StringValue(secret.Title())
	data.Notes = types.StringValue(secret.Notes())
	data.Login = types.StringValue(secret.GetFieldValueByType("login"))
	data.Password = types.StringValue(secret.GetFieldValueByType("password"))
	data.URL = types.StringValue(secret.GetFieldValueByType("url"))

	bankAccountList, diags := bankAccountToListValue(ctx, secret)
	resp.Diagnostics.Append(diags...)
	data.BankAccount = bankAccountList

	nameList, diags := nameToListValue(ctx, secret)
	resp.Diagnostics.Append(diags...)
	data.Name = nameList

	cardRefList, diags := cardRefToListValue(ctx, secret, client)
	resp.Diagnostics.Append(diags...)
	data.CardRef = cardRefList

	fileRefList, diags := fileItemsToListValue(ctx, secret.Files)
	resp.Diagnostics.Append(diags...)
	data.FileRef = fileRefList

	totpUrl := strings.TrimSpace(secret.GetFieldValueByType("oneTimeCode"))
	totpList, diags := totpToListValue(ctx, totpUrl)
	resp.Diagnostics.Append(diags...)
	data.TOTP = totpList

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
