package secretsmanager

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ ephemeral.EphemeralResource              = &ephemeralBankCard{}
	_ ephemeral.EphemeralResourceWithConfigure = &ephemeralBankCard{}
)

type ephemeralBankCard struct {
	meta providerMeta
}

type ephemeralBankCardModel struct {
	Path            types.String `tfsdk:"path"`
	Type            types.String `tfsdk:"type"`
	Title           types.String `tfsdk:"title"`
	Notes           types.String `tfsdk:"notes"`
	PaymentCard     types.List   `tfsdk:"payment_card"`
	CardholderName  types.String `tfsdk:"cardholder_name"`
	PinCode         types.String `tfsdk:"pin_code"`
	AddressRef      types.List   `tfsdk:"address_ref"`
	FileRef         types.List   `tfsdk:"file_ref"`
	Custom  types.List   `tfsdk:"custom"`
}

func NewEphemeralBankCard() ephemeral.EphemeralResource {
	return &ephemeralBankCard{}
}

func (e *ephemeralBankCard) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bank_card"
}

func (e *ephemeralBankCard) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this ephemeral resource to read a bank card record from Keeper Secrets Manager. Values are never stored in state.",
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
			"payment_card": paymentCardEphemeralAttribute(),
			"cardholder_name": schema.StringAttribute{
				Computed:    true,
				Description: "The cardholder name.",
			},
			"pin_code": schema.StringAttribute{
				Computed:    true,
				Description: "The PIN code.",
			},
			"address_ref": addressRefEphemeralAttribute(),
			"file_ref":    fileRefEphemeralAttribute(),
			"custom": genericFieldEphemeralAttribute("Custom fields of the record."),
		},
	}
}

func (e *ephemeralBankCard) Configure(_ context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
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

func (e *ephemeralBankCard) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data ephemeralBankCardModel
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
	if recordType != "bankCard" {
		resp.Diagnostics.AddError("Record Type Mismatch",
			"record type '"+recordType+"' is not the expected type 'bankCard' for this ephemeral resource")
		return
	}

	data.Type = types.StringValue(recordType)
	data.Title = types.StringValue(secret.Title())
	data.Notes = types.StringValue(secret.Notes())
	data.CardholderName = types.StringValue(secret.GetFieldValueByType("text"))
	data.PinCode = types.StringValue(secret.GetFieldValueByType("pinCode"))

	paymentCardList, diags := paymentCardToListValue(ctx, secret)
	resp.Diagnostics.Append(diags...)
	data.PaymentCard = paymentCardList

	addressRefList, diags := addressRefToListValue(ctx, secret, client)
	resp.Diagnostics.Append(diags...)
	data.AddressRef = addressRefList

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
