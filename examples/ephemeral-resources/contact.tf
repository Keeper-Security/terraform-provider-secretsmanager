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

ephemeral "secretsmanager_contact" "my_contact" {
  path = "<record UID>"
}

output "name" {
  value     = ephemeral.secretsmanager_contact.my_contact.name
  ephemeral = true
}

output "company" {
  value     = ephemeral.secretsmanager_contact.my_contact.company
  ephemeral = true
}

output "email" {
  value     = ephemeral.secretsmanager_contact.my_contact.email
  ephemeral = true
}
