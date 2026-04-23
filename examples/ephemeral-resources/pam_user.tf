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
  value     = ephemeral.secretsmanager_pam_user.db_admin_by_uid.login
  ephemeral = true
}

output "db_password" {
  value     = ephemeral.secretsmanager_pam_user.db_admin_by_uid.password
  ephemeral = true
}

output "db_folder_uid" {
  value     = ephemeral.secretsmanager_pam_user.db_admin_by_uid.folder_uid
  ephemeral = true
}

output "db_distinguished_name" {
  value     = ephemeral.secretsmanager_pam_user.db_admin_by_uid.distinguished_name
  ephemeral = true
}

output "db_connect_database" {
  value     = ephemeral.secretsmanager_pam_user.db_admin_by_uid.connect_database
  ephemeral = true
}

output "db_managed" {
  value     = ephemeral.secretsmanager_pam_user.db_admin_by_uid.managed
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
    username = local.db_user.login
    password = local.db_user.password
    database = local.db_user.connect_database
    managed  = local.db_user.managed
  }
}
