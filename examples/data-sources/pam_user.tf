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

# Example 1: Read PAM User by UID (recommended - always unique)
data "secretsmanager_pam_user" "db_admin_by_uid" {
  path = "EARSF1XFshSbkFmc84BBHA" # Replace with your record UID
}

# Example 2: Read PAM User by title (errors if multiple records have same title)
data "secretsmanager_pam_user" "db_admin_by_title" {
  title = "GatewayTest - RDP User" # Replace with your record title
}

# Output the PAM User data
output "db_login" {
  value = data.secretsmanager_pam_user.db_admin_by_uid.login[0].value
}

output "db_password" {
  value     = data.secretsmanager_pam_user.db_admin_by_uid.password[0].value
  sensitive = true
}

output "db_distinguished_name" {
  value = try(data.secretsmanager_pam_user.db_admin_by_uid.distinguished_name[0].value, "")
}

output "db_connect_database" {
  value = try(data.secretsmanager_pam_user.db_admin_by_uid.connect_database[0].value, "")
}

output "db_managed" {
  value = try(data.secretsmanager_pam_user.db_admin_by_uid.managed[0].value, false)
}

# Example: Check if TOTP is configured
output "has_2fa" {
  value = length(data.secretsmanager_pam_user.db_admin_by_uid.totp) > 0
}

# Example: Access rotation scripts
output "rotation_configured" {
  value = length(data.secretsmanager_pam_user.db_admin_by_uid.rotation_scripts) > 0
}

# Example: Build a connection string for database access
locals {
  db_user = data.secretsmanager_pam_user.db_admin_by_uid
  connection_info = {
    username = try(local.db_user.login[0].value, "")
    password = try(local.db_user.password[0].value, "")
    database = try(local.db_user.connect_database[0].value, "")
    managed  = try(local.db_user.managed[0].value, false)
  }
}

output "connection_info" {
  value = {
    username = local.connection_info.username
    database = local.connection_info.database
    managed  = local.connection_info.managed
  }
  sensitive = false
}

# Example: Use in another resource to configure database access
resource "null_resource" "db_access_setup" {
  triggers = {
    username = data.secretsmanager_pam_user.db_admin_by_uid.login[0].value
    database = try(data.secretsmanager_pam_user.db_admin_by_uid.connect_database[0].value, "")
  }

  provisioner "local-exec" {
    command = "echo Setting up database access for ${self.triggers.username} on ${self.triggers.database}"
  }
}

# Example: Extract LDAP distinguished name components
locals {
  dn       = try(data.secretsmanager_pam_user.db_admin_by_uid.distinguished_name[0].value, "")
  dn_parts = split(",", local.dn)
}

output "ldap_common_name" {
  value       = try(split("=", local.dn_parts[0])[1], "")
  description = "Extracted CN from distinguished name"
}

output "ldap_ou" {
  value       = try(split("=", local.dn_parts[1])[1], "")
  description = "Extracted OU from distinguished name"
}
