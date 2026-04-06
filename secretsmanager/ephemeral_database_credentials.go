package secretsmanager

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ ephemeral.EphemeralResource              = &ephemeralDatabaseCredentials{}
	_ ephemeral.EphemeralResourceWithConfigure = &ephemeralDatabaseCredentials{}
)

type ephemeralDatabaseCredentials struct {
	meta providerMeta
}

type ephemeralDatabaseCredentialsModel struct {
	Path     types.String `tfsdk:"path"`
	Type     types.String `tfsdk:"type"`
	Title    types.String `tfsdk:"title"`
	Notes    types.String `tfsdk:"notes"`
	DbType   types.String `tfsdk:"db_type"`
	Login    types.String `tfsdk:"login"`
	Password types.String `tfsdk:"password"`
	Host     types.List   `tfsdk:"host"`
	FileRef  types.List   `tfsdk:"file_ref"`
}

func NewEphemeralDatabaseCredentials() ephemeral.EphemeralResource {
	return &ephemeralDatabaseCredentials{}
}

func (e *ephemeralDatabaseCredentials) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_credentials"
}

func (e *ephemeralDatabaseCredentials) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this ephemeral resource to read database credentials from Keeper Secrets Manager. Values are never stored in state.",
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
			"db_type": schema.StringAttribute{
				Computed:    true,
				Description: "The database type.",
			},
			"login": schema.StringAttribute{
				Computed:    true,
				Description: "The secret login.",
			},
			"password": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The password.",
			},
			"host":     hostEphemeralAttribute(),
			"file_ref": fileRefEphemeralAttribute(),
		},
	}
}

func (e *ephemeralDatabaseCredentials) Configure(_ context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
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

func (e *ephemeralDatabaseCredentials) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data ephemeralDatabaseCredentialsModel
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
	if recordType != "databaseCredentials" {
		resp.Diagnostics.AddError("Record Type Mismatch",
			"record type '"+recordType+"' is not the expected type 'databaseCredentials' for this ephemeral resource")
		return
	}

	data.Type = types.StringValue(recordType)
	data.Title = types.StringValue(secret.Title())
	data.Notes = types.StringValue(secret.Notes())
	data.DbType = types.StringValue(secret.GetFieldValueByType("text"))
	data.Login = types.StringValue(secret.GetFieldValueByType("login"))
	data.Password = types.StringValue(secret.GetFieldValueByType("password"))

	hostList, diags := hostToListValue(ctx, secret)
	resp.Diagnostics.Append(diags...)
	data.Host = hostList

	fileRefList, diags := fileItemsToListValue(ctx, secret.Files)
	resp.Diagnostics.Append(diags...)
	data.FileRef = fileRefList

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
