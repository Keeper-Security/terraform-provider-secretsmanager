terraform {
  required_providers {
    secretsmanager = {
      source  = "keeper-security/secretsmanager"
      version = ">= 1.1.8"
    }
  }
}

provider "secretsmanager" {
  credential = "<CREDENTIAL>"
  # credential = file("~/.keeper/credential")
}

# Example 1: Read PAM Machine by UID (recommended - always unique)
data "secretsmanager_pam_machine" "ssh_server_by_uid" {
  path = "EYgeFFfgXRl_DQFIQD71uw"  # Replace with your record UID
}

# Example 2: Read PAM Machine by title (errors if multiple records have same title)
data "secretsmanager_pam_machine" "ssh_server_by_title" {
  title = "tf_acc_test_pam_machine_create"  # Replace with your record title
}

# Output the PAM Machine data
output "ssh_hostname" {
  value = data.secretsmanager_pam_machine.ssh_server_by_uid.pam_hostname[0].value[0].hostname
}

output "ssh_port" {
  value = data.secretsmanager_pam_machine.ssh_server_by_uid.pam_hostname[0].value[0].port
}

output "ssh_login" {
  value = data.secretsmanager_pam_machine.ssh_server_by_uid.login[0].value
}

output "ssh_password" {
  value     = data.secretsmanager_pam_machine.ssh_server_by_uid.password[0].value
  sensitive = true
}

# Access pamSettings as JSON
output "ssh_pam_settings" {
  value = jsondecode(data.secretsmanager_pam_machine.ssh_server_by_uid.pam_settings)
  sensitive = true
}

# Example: Extract specific settings from pamSettings
locals {
  ssh_settings = jsondecode(data.secretsmanager_pam_machine.ssh_server_by_uid.pam_settings)
  protocol = try(local.ssh_settings[0].connection[0].protocol, "unknown")
  recording_enabled = try(local.ssh_settings[0].connection[0].recordingIncludeKeys, false)
}

output "ssh_protocol" {
  value = local.protocol
}

output "ssh_recording_enabled" {
  value = local.recording_enabled
}

# Example: Access cloud instance metadata
output "instance_info" {
  value = {
    name = try(data.secretsmanager_pam_machine.ssh_server_by_uid.instance_name[0].value, "")
    id = try(data.secretsmanager_pam_machine.ssh_server_by_uid.instance_id[0].value, "")
    provider = try(data.secretsmanager_pam_machine.ssh_server_by_uid.provider_group[0].value, "")
    region = try(data.secretsmanager_pam_machine.ssh_server_by_uid.provider_region[0].value, "")
  }
}

# Example: Use in another resource
resource "null_resource" "ssh_connection" {
  triggers = {
    host = data.secretsmanager_pam_machine.ssh_server_by_uid.pam_hostname[0].value[0].hostname
    port = data.secretsmanager_pam_machine.ssh_server_by_uid.pam_hostname[0].value[0].port
  }

  provisioner "local-exec" {
    command = "echo Connecting to ${self.triggers.host}:${self.triggers.port}"
  }
}
