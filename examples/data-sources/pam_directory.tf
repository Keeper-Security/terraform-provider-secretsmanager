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

# Example 1: Read PAM Directory by UID (recommended - always unique)
data "secretsmanager_pam_directory" "ad_by_uid" {
  path = "AptSy2tZsUPhtaXjUlrxiQ" # Replace with your record UID
}

# Example 2: Read PAM Directory by title (errors if multiple records have same title)
data "secretsmanager_pam_directory" "ad_by_title" {
  title = "Test PAM Directory - Alternative IPs" # Replace with your record title
}

# Output the PAM Directory data
output "ad_hostname" {
  value = data.secretsmanager_pam_directory.ad_by_uid.pam_hostname[0].value[0].hostname
}

output "ad_port" {
  value = data.secretsmanager_pam_directory.ad_by_uid.pam_hostname[0].value[0].port
}

output "ad_directory_type" {
  value = data.secretsmanager_pam_directory.ad_by_uid.directory_type
}

output "ad_distinguished_name" {
  value = try(data.secretsmanager_pam_directory.ad_by_uid.distinguished_name[0].value, "")
}

output "ad_folder_uid" {
  value = data.secretsmanager_pam_directory.ad_by_uid.folder_uid
}

output "ad_totp_uri" {
  value     = try(data.secretsmanager_pam_directory.ad_by_uid.totp[0].value, "")
  sensitive = true
}

# Access pamSettings as JSON
output "ad_pam_settings" {
  value     = jsondecode(data.secretsmanager_pam_directory.ad_by_uid.pam_settings)
  sensitive = true
}

# Example: Extract specific settings from pamSettings
locals {
  ad_settings = jsondecode(data.secretsmanager_pam_directory.ad_by_uid.pam_settings)
  protocol    = try(local.ad_settings[0].connection[0].protocol, "unknown")
  port        = try(local.ad_settings[0].connection[0].port, "389")
  ssl_enabled = try(data.secretsmanager_pam_directory.ad_by_uid.use_ssl[0].value, false)
}

output "ad_protocol" {
  value = local.protocol
}

output "ad_connection_port" {
  value = local.port
}

output "ad_ssl_enabled" {
  value = local.ssl_enabled
}

# Example: Build connection string
output "ad_connection_info" {
  value = {
    type     = data.secretsmanager_pam_directory.ad_by_uid.directory_type
    host     = data.secretsmanager_pam_directory.ad_by_uid.pam_hostname[0].value[0].hostname
    port     = data.secretsmanager_pam_directory.ad_by_uid.pam_hostname[0].value[0].port
    protocol = local.protocol
    ssl      = local.ssl_enabled
    base_dn  = try(data.secretsmanager_pam_directory.ad_by_uid.distinguished_name[0].value, "")
  }
}

# Example: Use in another resource (e.g., LDAP client configuration)
resource "null_resource" "ldap_connection_test" {
  triggers = {
    host = data.secretsmanager_pam_directory.ad_by_uid.pam_hostname[0].value[0].hostname
    port = data.secretsmanager_pam_directory.ad_by_uid.pam_hostname[0].value[0].port
    type = data.secretsmanager_pam_directory.ad_by_uid.directory_type
  }

  provisioner "local-exec" {
    command = "echo Testing connection to ${self.triggers.type} at ${self.triggers.host}:${self.triggers.port}"
  }
}

# Example: Conditional output based on directory type
output "directory_specific_notes" {
  value = data.secretsmanager_pam_directory.ad_by_uid.directory_type == "Active Directory" ? "Using Active Directory - ensure LDAPS is configured" : "Using OpenLDAP - verify SSL configuration"
}
