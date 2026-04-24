package secretsmanager

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ ephemeral.EphemeralResource              = &ephemeralField{}
	_ ephemeral.EphemeralResourceWithConfigure = &ephemeralField{}
)

type ephemeralField struct {
	meta providerMeta
}

type ephemeralFieldModel struct {
	Path  types.String `tfsdk:"path"`
	Title types.String `tfsdk:"title"`
	Value types.String `tfsdk:"value"`
}

func NewEphemeralField() ephemeral.EphemeralResource {
	return &ephemeralField{}
}

func (e *ephemeralField) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_field"
}

func (e *ephemeralField) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this ephemeral resource to read a single field from a Keeper record using notation. Values are never stored in state.",
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				Required:    true,
				Description: "The path where the secret is stored.",
			},
			"title": schema.StringAttribute{
				Optional:    true,
				Description: "The secret title. (To find record by title - replace UID in path with '*')",
			},
			"value": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The value of the secret field.",
			},
		},
	}
}

func (e *ephemeralField) Configure(_ context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
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

func (e *ephemeralField) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data ephemeralFieldModel
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

	// find by title requested
	if title != "" && strings.Contains(path, "*") {
		records, err := getSecrets(client, []string{})
		if err != nil {
			resp.Diagnostics.AddError("Error fetching records", err.Error())
			return
		}
		uids := []string{}
		for _, r := range records {
			if r.Title() == title {
				uids = append(uids, r.Uid)
			}
		}
		if len(uids) != 1 {
			resp.Diagnostics.AddError("Record Not Found",
				fmt.Sprintf("expected 1 record - found %d records with title: %s", len(uids), title))
			return
		}
		path = strings.Replace(path, "*", uids[0], 1)
	}

	value, err := getNotation(client, path)
	if err != nil {
		resp.Diagnostics.AddError("Error reading field", err.Error())
		return
	}

	strValue := ""
	if len(value) == 1 {
		strValue = fmt.Sprintf("%v", value[0])
	} else {
		strValue = fmt.Sprintf("%v", value)
	}

	data.Value = types.StringValue(strValue)

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
