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

ephemeral "secretsmanager_ssn_card" "my_ssn" {
  path = "<record UID>"
}

output "identity_number" {
  value     = ephemeral.secretsmanager_ssn_card.my_ssn.identity_number
  ephemeral = true
}

output "name" {
  value     = ephemeral.secretsmanager_ssn_card.my_ssn.name
  ephemeral = true
}
