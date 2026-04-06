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

ephemeral "secretsmanager_ssh_keys" "my_ssh_keys" {
  path = "<record UID>"
}

output "login" {
  value     = ephemeral.secretsmanager_ssh_keys.my_ssh_keys.login
  ephemeral = true
}

output "passphrase" {
  value     = ephemeral.secretsmanager_ssh_keys.my_ssh_keys.passphrase
  ephemeral = true
}

output "public_key" {
  value     = length(ephemeral.secretsmanager_ssh_keys.my_ssh_keys.key_pair) < 1 ? "" : ephemeral.secretsmanager_ssh_keys.my_ssh_keys.key_pair.0.public_key
  ephemeral = true
}

output "private_key" {
  value     = length(ephemeral.secretsmanager_ssh_keys.my_ssh_keys.key_pair) < 1 ? "" : ephemeral.secretsmanager_ssh_keys.my_ssh_keys.key_pair.0.private_key
  ephemeral = true
}

output "host" {
  value     = length(ephemeral.secretsmanager_ssh_keys.my_ssh_keys.host) < 1 ? "" : ephemeral.secretsmanager_ssh_keys.my_ssh_keys.host.0.host_name
  ephemeral = true
}
