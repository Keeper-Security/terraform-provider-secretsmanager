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

ephemeral "secretsmanager_database_credentials" "my_db_creds" {
  path = "<record UID>"
}

output "db_type" {
  value     = ephemeral.secretsmanager_database_credentials.my_db_creds.db_type
  ephemeral = true
}

output "login" {
  value     = ephemeral.secretsmanager_database_credentials.my_db_creds.login
  ephemeral = true
}

output "password" {
  value     = ephemeral.secretsmanager_database_credentials.my_db_creds.password
  ephemeral = true
}

output "host" {
  value     = length(ephemeral.secretsmanager_database_credentials.my_db_creds.host) < 1 ? "" : ephemeral.secretsmanager_database_credentials.my_db_creds.host.0.host_name
  ephemeral = true
}

output "port" {
  value     = length(ephemeral.secretsmanager_database_credentials.my_db_creds.host) < 1 ? "" : ephemeral.secretsmanager_database_credentials.my_db_creds.host.0.port
  ephemeral = true
}
