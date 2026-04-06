package secretsmanager

import (
	"context"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	fwschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/secrets-manager-go/core"
)

var (
	_ provider.Provider                        = &fwProvider{}
	_ provider.ProviderWithEphemeralResources  = &fwProvider{}
)

// fwProvider is the Plugin Framework provider that serves ephemeral resources.
// It runs alongside the SDKv2 provider via the mux server.
type fwProvider struct {
	meta providerMeta
}

type fwProviderModel struct {
	Credential types.String `tfsdk:"credential"`
}

func NewFWProvider() provider.Provider {
	return &fwProvider{}
}

func (p *fwProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "secretsmanager"
}

// Schema must mirror the SDKv2 provider schema for mux compatibility.
func (p *fwProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = fwschema.Schema{
		Attributes: map[string]fwschema.Attribute{
			"credential": fwschema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Credential to use for Secrets Manager authentication. Can also be sourced from the `KEEPER_CREDENTIAL` environment variable.",
			},
		},
	}
}

func (p *fwProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config fwProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	creds := config.Credential.ValueString()
	if creds == "" {
		// Fall back to environment variable
		creds = envDefault("KEEPER_CREDENTIAL")
	}
	if strings.TrimSpace(creds) == "" {
		resp.Diagnostics.AddError("Missing Credential", "empty credential")
		return
	}

	ksmConfig := core.NewMemoryKeyValueStorage(creds)
	if ksmConfig.Get(core.KEY_APP_KEY) == "" || ksmConfig.Get(core.KEY_CLIENT_ID) == "" || ksmConfig.Get(core.KEY_PRIVATE_KEY) == "" {
		resp.Diagnostics.AddError(
			"Invalid Credentials",
			"Invalid credentials - please provide a valid base64 encoded KSM config. One-time tokens are not allowed.",
		)
		return
	}

	client := core.NewSecretsManager(&core.ClientOptions{Config: ksmConfig})
	p.meta = providerMeta{client: client}

	resp.EphemeralResourceData = p.meta
}

// Resources returns empty — all managed resources are served by the SDKv2 provider.
func (p *fwProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}

// DataSources returns empty — all data sources are served by the SDKv2 provider.
func (p *fwProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// Functions returns empty — no provider functions.
func (p *fwProvider) Functions(_ context.Context) []func() function.Function {
	return nil
}

// EphemeralResources registers all ephemeral resources served by the Framework provider.
func (p *fwProvider) EphemeralResources(_ context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		NewEphemeralLogin,
		NewEphemeralField,
		NewEphemeralRecord,
		NewEphemeralDatabaseCredentials,
		NewEphemeralServerCredentials,
		NewEphemeralSshKeys,
		NewEphemeralEncryptedNotes,
		NewEphemeralAddress,
		NewEphemeralBankAccount,
		NewEphemeralBankCard,
		NewEphemeralBirthCertificate,
		NewEphemeralContact,
		NewEphemeralDriverLicense,
		NewEphemeralHealthInsurance,
		NewEphemeralMembership,
		NewEphemeralPassport,
		NewEphemeralPhoto,
		NewEphemeralSoftwareLicense,
		NewEphemeralSsnCard,
		NewEphemeralFile,
		NewEphemeralPamUser,
		NewEphemeralPamMachine,
		NewEphemeralPamDatabase,
		NewEphemeralPamDirectory,
		NewEphemeralPamRemoteBrowser,
	}
}

func envDefault(key string) string {
	return os.Getenv(key)
}
