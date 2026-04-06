terraform {
  required_providers {
    secretsmanager = {
      source  = "keeper-security/secretsmanager"
      version = ">= 1.3.0"
    }
  }
}

provider "secretsmanager" {
  credential = "<CREDENTIAL>"
  # credential = file("~/.keeper/credential")
}

# Ephemeral resources do not store secret values in the Terraform state file.
# This makes them a more secure option for accessing sensitive credentials.

ephemeral "secretsmanager_address" "my_address" {
  path = "<record UID>"
}

output "zip_code" {
  value     = length(ephemeral.secretsmanager_address.my_address.address) < 1 ? "" : ephemeral.secretsmanager_address.my_address.address.0.zip
  ephemeral = true
}

output "city" {
  value     = length(ephemeral.secretsmanager_address.my_address.address) < 1 ? "" : ephemeral.secretsmanager_address.my_address.address.0.city
  ephemeral = true
}

output "state" {
  value     = length(ephemeral.secretsmanager_address.my_address.address) < 1 ? "" : ephemeral.secretsmanager_address.my_address.address.0.state
  ephemeral = true
}
