package secretsmanager

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ ephemeral.EphemeralResource              = &ephemeralHealthInsurance{}
	_ ephemeral.EphemeralResourceWithConfigure = &ephemeralHealthInsurance{}
)

type ephemeralHealthInsurance struct {
	meta providerMeta
}

type ephemeralHealthInsuranceModel struct {
	Path          types.String `tfsdk:"path"`
	Type          types.String `tfsdk:"type"`
	Title         types.String `tfsdk:"title"`
	Notes         types.String `tfsdk:"notes"`
	AccountNumber types.String `tfsdk:"account_number"`
	Name          types.List   `tfsdk:"name"`
	Login         types.String `tfsdk:"login"`
	Password      types.String `tfsdk:"password"`
	URL           types.String `tfsdk:"url"`
	FileRef       types.List   `tfsdk:"file_ref"`
}

func NewEphemeralHealthInsurance() ephemeral.EphemeralResource {
	return &ephemeralHealthInsurance{}
}

func (e *ephemeralHealthInsurance) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_health_insurance"
}

func (e *ephemeralHealthInsurance) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this ephemeral resource to read a health insurance record from Keeper Secrets Manager. Values are never stored in state.",
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
			"account_number": schema.StringAttribute{
				Computed:    true,
				Description: "Account Number.",
			},
			"name": nameEphemeralAttribute(),
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
			"file_ref": fileRefEphemeralAttribute(),
		},
	}
}

func (e *ephemeralHealthInsurance) Configure(_ context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
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

func (e *ephemeralHealthInsurance) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data ephemeralHealthInsuranceModel
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
	if recordType != "healthInsurance" {
		resp.Diagnostics.AddError("Record Type Mismatch",
			"record type '"+recordType+"' is not the expected type 'healthInsurance' for this ephemeral resource")
		return
	}

	data.Type = types.StringValue(recordType)
	data.Title = types.StringValue(secret.Title())
	data.Notes = types.StringValue(secret.Notes())
	data.AccountNumber = types.StringValue(secret.GetFieldValueByType("accountNumber"))
	data.Login = types.StringValue(secret.GetFieldValueByType("login"))
	data.Password = types.StringValue(secret.GetFieldValueByType("password"))
	data.URL = types.StringValue(secret.GetFieldValueByType("url"))

	nameList, diags := nameToListValue(ctx, secret)
	resp.Diagnostics.Append(diags...)
	data.Name = nameList

	fileRefList, diags := fileItemsToListValue(ctx, secret.Files)
	resp.Diagnostics.Append(diags...)
	data.FileRef = fileRefList

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
