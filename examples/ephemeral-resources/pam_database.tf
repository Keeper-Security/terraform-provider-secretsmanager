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

# Example 1: Read PAM Database by UID (recommended - always unique)
ephemeral "secretsmanager_pam_database" "mysql_by_uid" {
  path = "<record UID>" # Replace with your record UID
}

# Example 2: Read PAM Database by title (errors if multiple records have same title)
ephemeral "secretsmanager_pam_database" "mysql_by_title" {
  title = "Production MySQL Database" # Replace with your record title
}

# Output the PAM Database data
output "db_hostname" {
  value     = ephemeral.secretsmanager_pam_database.mysql_by_uid.pam_hostname[0].value[0].hostname
  ephemeral = true
}

output "db_port" {
  value     = ephemeral.secretsmanager_pam_database.mysql_by_uid.pam_hostname[0].value[0].port
  ephemeral = true
}

output "db_type" {
  value     = ephemeral.secretsmanager_pam_database.mysql_by_uid.database_type
  ephemeral = true
}

output "db_use_ssl" {
  value     = try(ephemeral.secretsmanager_pam_database.mysql_by_uid.use_ssl[0].value, false)
  ephemeral = true
}

output "db_folder_uid" {
  value     = ephemeral.secretsmanager_pam_database.mysql_by_uid.folder_uid
  ephemeral = true
}

# Example: Access cloud database metadata
output "cloud_db_info" {
  value = {
    database_id = try(ephemeral.secretsmanager_pam_database.mysql_by_uid.database_id[0].value, "")
    provider    = try(ephemeral.secretsmanager_pam_database.mysql_by_uid.provider_group[0].value, "")
    region      = try(ephemeral.secretsmanager_pam_database.mysql_by_uid.provider_region[0].value, "")
  }
  ephemeral = true
}
