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

# Example 1: Read PAM Directory by UID (recommended - always unique)
ephemeral "secretsmanager_pam_directory" "ad_by_uid" {
  path = "<record UID>" # Replace with your record UID
}

# Example 2: Read PAM Directory by title (errors if multiple records have same title)
ephemeral "secretsmanager_pam_directory" "ad_by_title" {
  title = "Production Active Directory" # Replace with your record title
}

# Output the PAM Directory data
output "ad_hostname" {
  value     = ephemeral.secretsmanager_pam_directory.ad_by_uid.pam_hostname[0].host_name
  ephemeral = true
}

output "ad_port" {
  value     = ephemeral.secretsmanager_pam_directory.ad_by_uid.pam_hostname[0].port
  ephemeral = true
}

output "ad_directory_type" {
  value     = ephemeral.secretsmanager_pam_directory.ad_by_uid.directory_type
  ephemeral = true
}

output "ad_distinguished_name" {
  value     = ephemeral.secretsmanager_pam_directory.ad_by_uid.distinguished_name
  ephemeral = true
}

output "ad_folder_uid" {
  value     = ephemeral.secretsmanager_pam_directory.ad_by_uid.folder_uid
  ephemeral = true
}

output "ad_use_ssl" {
  value     = ephemeral.secretsmanager_pam_directory.ad_by_uid.use_ssl
  ephemeral = true
}

# Example: Build connection info
output "ad_connection_info" {
  value = {
    type    = ephemeral.secretsmanager_pam_directory.ad_by_uid.directory_type
    host    = ephemeral.secretsmanager_pam_directory.ad_by_uid.pam_hostname[0].host_name
    port    = ephemeral.secretsmanager_pam_directory.ad_by_uid.pam_hostname[0].port
    ssl     = ephemeral.secretsmanager_pam_directory.ad_by_uid.use_ssl
    base_dn = ephemeral.secretsmanager_pam_directory.ad_by_uid.distinguished_name
  }
  ephemeral = true
}
