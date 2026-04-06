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

ephemeral "secretsmanager_bank_account" "my_account" {
  path = "<record UID>"
}

output "login" {
  value     = ephemeral.secretsmanager_bank_account.my_account.login
  ephemeral = true
}

output "password" {
  value     = ephemeral.secretsmanager_bank_account.my_account.password
  ephemeral = true
}

output "name" {
  value     = ephemeral.secretsmanager_bank_account.my_account.name
  ephemeral = true
}

output "totp_token" {
  value     = length(ephemeral.secretsmanager_bank_account.my_account.totp) < 1 ? "" : ephemeral.secretsmanager_bank_account.my_account.totp.0.token
  ephemeral = true
}
