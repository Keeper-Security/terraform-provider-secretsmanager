# Keeper Secrets Manager Provider

The [Keeper Secrets Manager](https://docs.keeper.io/secrets-manager/) provider is used to interact with the
resources supported by Secrets Manager. The provider needs to be configured with Keeper credentials before it can be used.

You can set environment variable `KEEPER_CREDENTIAL` or read it from disk using the `file()` function.

## Installation

### Terraform 0.13+ ([Terraform Registry](https://registry.terraform.io/))
To install this provider, copy and paste this code into your Terraform configuration. Then, run `terraform init`.
```hcl
terraform {
  required_providers {
    secretsmanager = {
      source = "keeper-security/secretsmanager"
      version = ">= 1.1.2"
    }
  }
}

provider "secretsmanager" {
  # Configuration options
}
```

### Manual Install

Download archive with the [latest release](https://github.com/keeper-security/terraform-provider-secretsmanager/releases/latest) for your platform and copy it to the corresponding plugin folder (_Linux and MacOS:_ `~/.terraform.d/plugins/github.com/keeper-security/secretsmanager` _Windows:_ `%APPDATA%/terraform.d/plugins/github.com/keeper-security/secretsmanager`)  
Use the same config from above just remember to initialize `source` with the full URL `source  = "github.com/keeper-security/secretsmanager"`

For help on manually installing Terraform Providers, please refer to the [official Terraform documentation](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins).

## Example Usage

```hcl
provider "secretsmanager" {
  credential = file("~/.keeper/credential")
}
```

## Argument Reference

The following arguments are supported:

* `credential` - (Required) Credential to use for Secrets Manager authentication. Can also be sourced from the `KEEPER_CREDENTIAL` environment variable.
