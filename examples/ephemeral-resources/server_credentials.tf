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

ephemeral "secretsmanager_server_credentials" "my_server_creds" {
  path = "<record UID>"
}

output "login" {
  value     = ephemeral.secretsmanager_server_credentials.my_server_creds.login
  ephemeral = true
}

output "password" {
  value     = ephemeral.secretsmanager_server_credentials.my_server_creds.password
  ephemeral = true
}

output "host" {
  value     = length(ephemeral.secretsmanager_server_credentials.my_server_creds.host) < 1 ? "" : ephemeral.secretsmanager_server_credentials.my_server_creds.host.0.host_name
  ephemeral = true
}

output "port" {
  value     = length(ephemeral.secretsmanager_server_credentials.my_server_creds.host) < 1 ? "" : ephemeral.secretsmanager_server_credentials.my_server_creds.host.0.port
  ephemeral = true
}
