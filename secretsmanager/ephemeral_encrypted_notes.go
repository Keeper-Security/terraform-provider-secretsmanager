package secretsmanager

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ ephemeral.EphemeralResource              = &ephemeralEncryptedNotes{}
	_ ephemeral.EphemeralResourceWithConfigure = &ephemeralEncryptedNotes{}
)

type ephemeralEncryptedNotes struct {
	meta providerMeta
}

type ephemeralEncryptedNotesModel struct {
	Path    types.String `tfsdk:"path"`
	Type    types.String `tfsdk:"type"`
	Title   types.String `tfsdk:"title"`
	Notes   types.String `tfsdk:"notes"`
	Note    types.String `tfsdk:"note"`
	Date    types.String `tfsdk:"date"`
	FileRef types.List   `tfsdk:"file_ref"`
	Custom  types.List   `tfsdk:"custom"`
}

func NewEphemeralEncryptedNotes() ephemeral.EphemeralResource {
	return &ephemeralEncryptedNotes{}
}

func (e *ephemeralEncryptedNotes) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_encrypted_notes"
}

func (e *ephemeralEncryptedNotes) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this ephemeral resource to read encrypted notes from Keeper Secrets Manager. Values are never stored in state.",
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
			"note": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "Encrypted note.",
			},
			"date": schema.StringAttribute{
				Computed:    true,
				Description: "Date.",
			},
			"file_ref": fileRefEphemeralAttribute(),
			"custom": genericFieldEphemeralAttribute("Custom fields of the record."),
		},
	}
}

func (e *ephemeralEncryptedNotes) Configure(_ context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
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

func (e *ephemeralEncryptedNotes) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data ephemeralEncryptedNotesModel
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
	if recordType != "encryptedNotes" {
		resp.Diagnostics.AddError("Record Type Mismatch",
			"record type '"+recordType+"' is not the expected type 'encryptedNotes' for this ephemeral resource")
		return
	}

	data.Type = types.StringValue(recordType)
	data.Title = types.StringValue(secret.Title())
	data.Notes = types.StringValue(secret.Notes())
	data.Note = types.StringValue(secret.GetFieldValueByType("note"))
	data.Date = types.StringValue(dateFieldToString(secret.GetFieldValueByType("date")))

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
