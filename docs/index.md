# Keeper Secrets Manager Provider

The [Keeper Secrets Manager](https://docs.keeper.io/secrets-manager/) provider is used to interact with the
resources supported by Secrets Manager. The provider needs to be configured with Keeper credentials before it can be used.

You can set environment variable `KEEPER_CREDENTIAL` or read it from disk using the `file()` function.

## Installation

### Manual Install

Get the latest version of the Terraform Provider from [GitHub](https://github.com/keeper-security/terraform-provider-keeper) as a single [zip](https://github.com/Keeper-Security/terraform-provider-keeper/releases) archive or clone with git
```git
git clone https://github.com/keeper-security/terraform-provider-keeper
```
Build
```
go build
```
Copy plugin to Terraform plugin folder
```bash
cp terraform-provider-keeper ~/.terraform.d/plugins/terraform-provider-keeper_v0.1.0
```
Note: Default plugin path is %APPDATA%\terraform.d\plugins for Windows and ~/.terraform.d/plugins for all other operating systems.

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
