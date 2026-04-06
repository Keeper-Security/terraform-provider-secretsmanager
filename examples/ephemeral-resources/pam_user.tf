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

# Example 1: Read PAM User by UID (recommended - always unique)
ephemeral "secretsmanager_pam_user" "db_admin_by_uid" {
  path = "<record UID>" # Replace with your record UID
}

# Example 2: Read PAM User by title (errors if multiple records have same title)
ephemeral "secretsmanager_pam_user" "db_admin_by_title" {
  title = "GatewayTest - RDP User" # Replace with your record title
}

# Output the PAM User data
output "db_login" {
  value     = ephemeral.secretsmanager_pam_user.db_admin_by_uid.login[0].value
  ephemeral = true
}

output "db_password" {
  value     = ephemeral.secretsmanager_pam_user.db_admin_by_uid.password[0].value
  ephemeral = true
}

output "db_folder_uid" {
  value     = ephemeral.secretsmanager_pam_user.db_admin_by_uid.folder_uid
  ephemeral = true
}

output "db_distinguished_name" {
  value     = try(ephemeral.secretsmanager_pam_user.db_admin_by_uid.distinguished_name[0].value, "")
  ephemeral = true
}

output "db_connect_database" {
  value     = try(ephemeral.secretsmanager_pam_user.db_admin_by_uid.connect_database[0].value, "")
  ephemeral = true
}

output "db_managed" {
  value     = try(ephemeral.secretsmanager_pam_user.db_admin_by_uid.managed[0].value, false)
  ephemeral = true
}

# Example: Check if TOTP is configured
output "has_2fa" {
  value     = length(ephemeral.secretsmanager_pam_user.db_admin_by_uid.totp) > 0
  ephemeral = true
}

# Example: Build a connection string for database access
locals {
  db_user = ephemeral.secretsmanager_pam_user.db_admin_by_uid
  connection_info = {
    username = try(local.db_user.login[0].value, "")
    password = try(local.db_user.password[0].value, "")
    database = try(local.db_user.connect_database[0].value, "")
    managed  = try(local.db_user.managed[0].value, false)
  }
}
