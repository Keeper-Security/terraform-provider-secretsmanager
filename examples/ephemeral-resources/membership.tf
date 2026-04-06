terraform {
  required_providers {
    secretsmanager = {
      source  = "keeper-security/secretsmanager"
      version = ">= 1.2.0"
    }
  }
}

provider "secretsmanager" {
  credential = "<CREDENTIAL>"
  # credential = file("~/.keeper/credential")
}

# Ephemeral resources do not store secret values in the Terraform state file.
# This makes them a more secure option for accessing sensitive credentials.

ephemeral "secretsmanager_membership" "my_membership" {
  path = "<record UID>"
}

output "account_number" {
  value     = ephemeral.secretsmanager_membership.my_membership.account_number
  ephemeral = true
}

output "password" {
  value     = ephemeral.secretsmanager_membership.my_membership.password
  ephemeral = true
}

output "name" {
  value     = ephemeral.secretsmanager_membership.my_membership.name
  ephemeral = true
}
