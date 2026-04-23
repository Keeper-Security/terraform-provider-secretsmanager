package secretsmanager

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ ephemeral.EphemeralResource              = &ephemeralRecord{}
	_ ephemeral.EphemeralResourceWithConfigure = &ephemeralRecord{}
)

type ephemeralRecord struct {
	meta providerMeta
}

type ephemeralRecordModel struct {
	Path    types.String `tfsdk:"path"`
	Type    types.String `tfsdk:"type"`
	Title   types.String `tfsdk:"title"`
	Notes   types.String `tfsdk:"notes"`
	Fields  types.List   `tfsdk:"fields"`
	Custom  types.List   `tfsdk:"custom"`
	FileRef types.List   `tfsdk:"file_ref"`
}

func NewEphemeralRecord() ephemeral.EphemeralResource {
	return &ephemeralRecord{}
}

func (e *ephemeralRecord) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_record"
}

func (e *ephemeralRecord) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this ephemeral resource to read a generic record from Keeper Secrets Manager. Values are never stored in state.",
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
			"fields": genericFieldEphemeralAttribute("Standard fields of the record."),
			"custom": genericFieldEphemeralAttribute("Custom fields of the record."),
			"file_ref": fileRefEphemeralAttribute(),
		},
	}
}

func (e *ephemeralRecord) Configure(_ context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
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

func (e *ephemeralRecord) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data ephemeralRecordModel
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

	data.Type = types.StringValue(secret.Type())
	data.Title = types.StringValue(secret.Title())
	data.Notes = types.StringValue(secret.Notes())

	fieldsItems := getFieldItemsData(secret.RecordDict, "fields")
	fieldsList, diags := genericFieldItemsToListValue(ctx, fieldsItems)
	resp.Diagnostics.Append(diags...)
	data.Fields = fieldsList

	customItems := getFieldItemsData(secret.RecordDict, "custom")
	customList, diags := genericFieldItemsToListValue(ctx, customItems)
	resp.Diagnostics.Append(diags...)
	data.Custom = customList

	fileRefList, diags := fileItemsToListValue(ctx, secret.Files)
	resp.Diagnostics.Append(diags...)
	data.FileRef = fileRefList

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}

// genericFieldEphemeralAttribute returns a computed list nested attribute for generic fields.
func genericFieldEphemeralAttribute(description string) schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Computed:    true,
		Description: description,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Computed:    true,
					Description: "The field type.",
				},
				"label": schema.StringAttribute{
					Computed:    true,
					Description: "The field label.",
				},
				"value": schema.StringAttribute{
					Computed:    true,
					Sensitive:   true,
					Description: "The field value. Complex types (phone, name, address, paymentCard) are returned as JSON.",
				},
				"required": schema.BoolAttribute{
					Computed:    true,
					Description: "Whether this field is required.",
				},
				"privacy_screen": schema.BoolAttribute{
					Computed:    true,
					Description: "Whether this field is hidden behind a privacy screen in the Keeper UI.",
				},
			},
		},
	}
}

var genericFieldObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"type":           types.StringType,
		"label":          types.StringType,
		"value":          types.StringType,
		"required":       types.BoolType,
		"privacy_screen": types.BoolType,
	},
}

// genericFieldItemsToListValue converts the SDKv2 getFieldItemsData result to a Framework types.List.
func genericFieldItemsToListValue(ctx context.Context, items []interface{}) (types.List, diag.Diagnostics) {
	if len(items) == 0 {
		return types.ListValueMust(genericFieldObjectType, []attr.Value{}), nil
	}

	objects := make([]attr.Value, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		fieldType := ""
		if v, ok := m["type"].(string); ok {
			fieldType = v
		}
		label := ""
		if v, ok := m["label"].(string); ok {
			label = v
		}

		value, _ := m["value"].(string)

		required, _ := m["required"].(bool)
		privacyScreen, _ := m["privacy_screen"].(bool)

		obj, d := types.ObjectValue(genericFieldObjectType.AttrTypes, map[string]attr.Value{
			"type":           types.StringValue(fieldType),
			"label":          types.StringValue(label),
			"value":          types.StringValue(value),
			"required":       types.BoolValue(required),
			"privacy_screen": types.BoolValue(privacyScreen),
		})
		if d.HasError() {
			return types.ListNull(genericFieldObjectType), d
		}
		objects = append(objects, obj)
	}

	return types.ListValue(genericFieldObjectType, objects)
}
