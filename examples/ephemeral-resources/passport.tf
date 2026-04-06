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

ephemeral "secretsmanager_passport" "my_passport" {
  path = "<record UID>"
}

output "passport_number" {
  value     = ephemeral.secretsmanager_passport.my_passport.passport_number
  ephemeral = true
}

output "name" {
  value     = ephemeral.secretsmanager_passport.my_passport.name
  ephemeral = true
}

output "expiration_date" {
  value     = ephemeral.secretsmanager_passport.my_passport.expiration_date
  ephemeral = true
}
