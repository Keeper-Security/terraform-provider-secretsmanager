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
  value     = ephemeral.secretsmanager_pam_machine.ssh_server_by_uid.pam_hostname[0].host_name
  ephemeral = true
}

output "ssh_port" {
  value     = ephemeral.secretsmanager_pam_machine.ssh_server_by_uid.pam_hostname[0].port
  ephemeral = true
}

output "machine_folder_uid" {
  value     = ephemeral.secretsmanager_pam_machine.ssh_server_by_uid.folder_uid
  ephemeral = true
}

output "machine_login" {
  value     = ephemeral.secretsmanager_pam_machine.ssh_server_by_uid.login
  ephemeral = true
}

output "machine_ssl_verification" {
  value     = ephemeral.secretsmanager_pam_machine.ssh_server_by_uid.ssl_verification
  ephemeral = true
}

output "machine_private_pem_key" {
  value     = ephemeral.secretsmanager_pam_machine.ssh_server_by_uid.private_pem_key
  ephemeral = true
}

# Example: Access cloud instance metadata
output "instance_info" {
  value = {
    name     = ephemeral.secretsmanager_pam_machine.ssh_server_by_uid.instance_name
    id       = ephemeral.secretsmanager_pam_machine.ssh_server_by_uid.instance_id
    provider = ephemeral.secretsmanager_pam_machine.ssh_server_by_uid.provider_group
    region   = ephemeral.secretsmanager_pam_machine.ssh_server_by_uid.provider_region
  }
  ephemeral = true
}
