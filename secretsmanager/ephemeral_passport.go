package secretsmanager

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ ephemeral.EphemeralResource              = &ephemeralPassport{}
	_ ephemeral.EphemeralResourceWithConfigure = &ephemeralPassport{}
)

type ephemeralPassport struct {
	meta providerMeta
}

type ephemeralPassportModel struct {
	Path           types.String `tfsdk:"path"`
	Type           types.String `tfsdk:"type"`
	Title          types.String `tfsdk:"title"`
	Notes          types.String `tfsdk:"notes"`
	PassportNumber types.String `tfsdk:"passport_number"`
	Name           types.List   `tfsdk:"name"`
	BirthDate      types.String `tfsdk:"birth_date"`
	ExpirationDate types.String `tfsdk:"expiration_date"`
	DateIssued     types.String `tfsdk:"date_issued"`
	Password       types.String `tfsdk:"password"`
	AddressRef     types.List   `tfsdk:"address_ref"`
	FileRef        types.List   `tfsdk:"file_ref"`
}

func NewEphemeralPassport() ephemeral.EphemeralResource {
	return &ephemeralPassport{}
}

func (e *ephemeralPassport) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_passport"
}

func (e *ephemeralPassport) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this ephemeral resource to read a passport record from Keeper Secrets Manager. Values are never stored in state.",
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
			"passport_number": schema.StringAttribute{
				Computed:    true,
				Description: "Passport Number.",
			},
			"name": nameEphemeralAttribute(),
			"birth_date": schema.StringAttribute{
				Computed:    true,
				Description: "Date of birth.",
			},
			"expiration_date": schema.StringAttribute{
				Computed:    true,
				Description: "Expiration date.",
			},
			"date_issued": schema.StringAttribute{
				Computed:    true,
				Description: "Date issued.",
			},
			"password": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The secret password.",
			},
			"address_ref": addressRefEphemeralAttribute(),
			"file_ref":    fileRefEphemeralAttribute(),
		},
	}
}

func (e *ephemeralPassport) Configure(_ context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
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

func (e *ephemeralPassport) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data ephemeralPassportModel
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
	if recordType != "passport" {
		resp.Diagnostics.AddError("Record Type Mismatch",
			"record type '"+recordType+"' is not the expected type 'passport' for this ephemeral resource")
		return
	}

	data.Type = types.StringValue(recordType)
	data.Title = types.StringValue(secret.Title())
	data.Notes = types.StringValue(secret.Notes())
	data.PassportNumber = types.StringValue(secret.GetFieldValueByType("accountNumber"))
	data.BirthDate = types.StringValue(dateFieldToString(secret.GetFieldValueByType("birthDate")))
	data.ExpirationDate = types.StringValue(dateFieldToString(secret.GetFieldValueByType("expirationDate")))
	data.DateIssued = types.StringValue(dateFieldToString(secret.GetFieldValueByType("date")))
	data.Password = types.StringValue(secret.GetFieldValueByType("password"))

	nameList, diags := nameToListValue(ctx, secret)
	resp.Diagnostics.Append(diags...)
	data.Name = nameList

	addressRefList, diags := addressRefToListValue(ctx, secret, client)
	resp.Diagnostics.Append(diags...)
	data.AddressRef = addressRefList

	fileRefList, diags := fileItemsToListValue(ctx, secret.Files)
	resp.Diagnostics.Append(diags...)
	data.FileRef = fileRefList

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
