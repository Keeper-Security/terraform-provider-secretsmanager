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

ephemeral "secretsmanager_login" "db_server" {
  path = "<record UID>"
}

output "login" {
  value     = ephemeral.secretsmanager_login.db_server.login
  ephemeral = true
}

output "password" {
  value     = ephemeral.secretsmanager_login.db_server.password
  ephemeral = true
}

output "url" {
  value     = ephemeral.secretsmanager_login.db_server.url
  ephemeral = true
}

output "title" {
  value     = ephemeral.secretsmanager_login.db_server.title
  ephemeral = true
}

output "totp_token" {
  value     = length(ephemeral.secretsmanager_login.db_server.totp) < 1 ? "" : ephemeral.secretsmanager_login.db_server.totp.0.token
  ephemeral = true
}
