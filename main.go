package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
	"github.com/keeper-security/terraform-provider-secretsmanager/secretsmanager"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	ctx := context.Background()
	serverFactory, err := secretsmanager.ProtoV6ProviderServerFactory(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf6server.ServeOpt
	if debug {
		serveOpts = append(serveOpts, tf6server.WithManagedDebug())
	}

	err = tf6server.Serve(
		"registry.terraform.io/keeper-security/secretsmanager",
		serverFactory,
		serveOpts...,
	)
	if err != nil {
		log.Fatal(err)
	}
}
