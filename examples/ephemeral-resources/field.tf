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

ephemeral "secretsmanager_field" "my_field" {
  path = "<record UID>/field/login"
}

output "field_value" {
  value     = ephemeral.secretsmanager_field.my_field.value
  ephemeral = true
}
