# Keeper Secrets Manager Provider

The [Keeper Secrets Manager](https://docs.keeper.io/secrets-manager/) provider is used to interact with the
resources supported by Secrets Manager. The provider needs to be configured with Keeper credentials before it can be used.

You can set environment variable `KEEPER_CREDENTIAL` or read it from disk using the `file()` function.

## Installation

### Manual Install

Download archive from [latest release](https://github.com/keeper-security/terraform-provider-keeper/releases/latest) for your platform and copy it to the corresponding plugin folder (_Linux and MacOS:_ `~/.terraform.d/plugins/github.com/keeper-security/keeper` _Windows:_ `%APPDATA%/terraform.d/plugins/github.com/keeper-security/keeper`) 

For help on manually installing Terraform Providers, please refer to the [official Terraform documentation](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins).

## Example Usage

```hcl
provider "keeper" {
  credential = file("~/.keeper/credential")
}
```

## Argument Reference

The following arguments are supported:

* `credential` - (Required) Credential to use for Secrets Manager authentication. Can also be sourced from the `KEEPER_CREDENTIAL` environment variable.
