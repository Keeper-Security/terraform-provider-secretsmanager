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

ephemeral "secretsmanager_health_insurance" "my_insurance" {
  path = "<record UID>"
}

output "login" {
  value     = ephemeral.secretsmanager_health_insurance.my_insurance.login
  ephemeral = true
}

output "password" {
  value     = ephemeral.secretsmanager_health_insurance.my_insurance.password
  ephemeral = true
}

output "account_number" {
  value     = ephemeral.secretsmanager_health_insurance.my_insurance.account_number
  ephemeral = true
}

output "url" {
  value     = ephemeral.secretsmanager_health_insurance.my_insurance.url
  ephemeral = true
}
