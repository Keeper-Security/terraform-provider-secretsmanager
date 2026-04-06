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

ephemeral "secretsmanager_record" "my_record" {
  path = "<record UID>"
}

output "record_type" {
  value     = ephemeral.secretsmanager_record.my_record.type
  ephemeral = true
}

output "record_title" {
  value     = ephemeral.secretsmanager_record.my_record.title
  ephemeral = true
}

output "first_field_type" {
  value     = length(ephemeral.secretsmanager_record.my_record.fields) < 1 ? "" : ephemeral.secretsmanager_record.my_record.fields.0.type
  ephemeral = true
}

output "first_field_value" {
  value     = length(ephemeral.secretsmanager_record.my_record.fields) < 1 ? "" : ephemeral.secretsmanager_record.my_record.fields.0.value
  ephemeral = true
}
