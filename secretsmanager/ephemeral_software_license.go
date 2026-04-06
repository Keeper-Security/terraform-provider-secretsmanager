package secretsmanager

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ ephemeral.EphemeralResource              = &ephemeralSoftwareLicense{}
	_ ephemeral.EphemeralResourceWithConfigure = &ephemeralSoftwareLicense{}
)

type ephemeralSoftwareLicense struct {
	meta providerMeta
}

type ephemeralSoftwareLicenseModel struct {
	Path           types.String `tfsdk:"path"`
	Type           types.String `tfsdk:"type"`
	Title          types.String `tfsdk:"title"`
	Notes          types.String `tfsdk:"notes"`
	LicenseNumber  types.String `tfsdk:"license_number"`
	ActivationDate types.String `tfsdk:"activation_date"`
	ExpirationDate types.String `tfsdk:"expiration_date"`
	FileRef        types.List   `tfsdk:"file_ref"`
}

func NewEphemeralSoftwareLicense() ephemeral.EphemeralResource {
	return &ephemeralSoftwareLicense{}
}

func (e *ephemeralSoftwareLicense) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_software_license"
}

func (e *ephemeralSoftwareLicense) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this ephemeral resource to read a software license record from Keeper Secrets Manager. Values are never stored in state.",
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
			"license_number": schema.StringAttribute{
				Computed:    true,
				Description: "License Number.",
			},
			"activation_date": schema.StringAttribute{
				Computed:    true,
				Description: "Date of activation.",
			},
			"expiration_date": schema.StringAttribute{
				Computed:    true,
				Description: "Date of expiration.",
			},
			"file_ref": fileRefEphemeralAttribute(),
		},
	}
}

func (e *ephemeralSoftwareLicense) Configure(_ context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
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

func (e *ephemeralSoftwareLicense) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data ephemeralSoftwareLicenseModel
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
	if recordType != "softwareLicense" {
		resp.Diagnostics.AddError("Record Type Mismatch",
			"record type '"+recordType+"' is not the expected type 'softwareLicense' for this ephemeral resource")
		return
	}

	data.Type = types.StringValue(recordType)
	data.Title = types.StringValue(secret.Title())
	data.Notes = types.StringValue(secret.Notes())
	data.LicenseNumber = types.StringValue(secret.GetFieldValueByType("licenseNumber"))
	// data source uses "date" for activation_date and "expirationDate" for expiration_date
	data.ActivationDate = types.StringValue(dateFieldToString(secret.GetFieldValueByType("date")))
	data.ExpirationDate = types.StringValue(dateFieldToString(secret.GetFieldValueByType("expirationDate")))

	fileRefList, diags := fileItemsToListValue(ctx, secret.Files)
	resp.Diagnostics.Append(diags...)
	data.FileRef = fileRefList

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
