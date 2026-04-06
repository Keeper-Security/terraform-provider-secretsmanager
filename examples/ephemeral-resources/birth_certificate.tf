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

ephemeral "secretsmanager_birth_certificate" "my_birth_cert" {
  path = "<record UID>"
}

output "birth_date" {
  value     = ephemeral.secretsmanager_birth_certificate.my_birth_cert.birth_date
  ephemeral = true
}

output "name" {
  value     = ephemeral.secretsmanager_birth_certificate.my_birth_cert.name
  ephemeral = true
}
