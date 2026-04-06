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

# Example 1: Read PAM Machine by UID (recommended - always unique)
ephemeral "secretsmanager_pam_machine" "ssh_server_by_uid" {
  path = "<record UID>" # Replace with your record UID
}

# Example 2: Read PAM Machine by title (errors if multiple records have same title)
ephemeral "secretsmanager_pam_machine" "ssh_server_by_title" {
  title = "Production SSH Server" # Replace with your record title
}

# Output the PAM Machine data
output "ssh_hostname" {
  value     = ephemeral.secretsmanager_pam_machine.ssh_server_by_uid.pam_hostname[0].value[0].hostname
  ephemeral = true
}

output "ssh_port" {
  value     = ephemeral.secretsmanager_pam_machine.ssh_server_by_uid.pam_hostname[0].value[0].port
  ephemeral = true
}

output "machine_folder_uid" {
  value     = ephemeral.secretsmanager_pam_machine.ssh_server_by_uid.folder_uid
  ephemeral = true
}

output "machine_login" {
  value     = try(ephemeral.secretsmanager_pam_machine.ssh_server_by_uid.login[0].value, "")
  ephemeral = true
}

output "machine_ssl_verification" {
  value     = try(ephemeral.secretsmanager_pam_machine.ssh_server_by_uid.ssl_verification[0].value, false)
  ephemeral = true
}

output "machine_private_pem_key" {
  value     = try(ephemeral.secretsmanager_pam_machine.ssh_server_by_uid.private_pem_key[0].value, "")
  ephemeral = true
}

# Example: Access cloud instance metadata
output "instance_info" {
  value = {
    name     = try(ephemeral.secretsmanager_pam_machine.ssh_server_by_uid.instance_name[0].value, "")
    id       = try(ephemeral.secretsmanager_pam_machine.ssh_server_by_uid.instance_id[0].value, "")
    provider = try(ephemeral.secretsmanager_pam_machine.ssh_server_by_uid.provider_group[0].value, "")
    region   = try(ephemeral.secretsmanager_pam_machine.ssh_server_by_uid.provider_region[0].value, "")
  }
  ephemeral = true
}
