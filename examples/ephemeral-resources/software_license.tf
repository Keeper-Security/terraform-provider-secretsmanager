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

ephemeral "secretsmanager_software_license" "my_license" {
  path = "<record UID>"
}

output "license_number" {
  value     = ephemeral.secretsmanager_software_license.my_license.license_number
  ephemeral = true
}

output "expiration_date" {
  value     = ephemeral.secretsmanager_software_license.my_license.expiration_date
  ephemeral = true
}

output "activation_date" {
  value     = ephemeral.secretsmanager_software_license.my_license.activation_date
  ephemeral = true
}
