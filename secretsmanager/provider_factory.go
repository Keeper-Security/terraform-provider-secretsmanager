package secretsmanager

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ProtoV6ProviderServerFactory returns a mux server that combines the SDKv2
// provider (resources + data sources) with the Framework provider (ephemeral resources).
// Uses protocol v6 to support nested attributes in ephemeral resource schemas.
func ProtoV6ProviderServerFactory(ctx context.Context) (func() tfprotov6.ProviderServer, error) {
	sdkv2Provider := Provider()

	// Upgrade SDKv2 provider from protocol v5 to v6
	upgradedSdkv2, err := tf5to6server.UpgradeServer(ctx, sdkv2Provider.GRPCProvider)
	if err != nil {
		return nil, err
	}

	servers := []func() tfprotov6.ProviderServer{
		func() tfprotov6.ProviderServer { return upgradedSdkv2 },
		providerserver.NewProtocol6(NewFWProvider()),
	}

	muxServer, err := tf6muxserver.NewMuxServer(ctx, servers...)
	if err != nil {
		return nil, err
	}

	return muxServer.ProviderServer, nil
}

// SDKv2Provider returns the SDKv2 provider for use in acceptance tests.
func SDKv2Provider() *schema.Provider {
	return Provider()
}
