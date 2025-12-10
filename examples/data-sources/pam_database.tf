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

# Example 1: Read PAM Database by UID (recommended - always unique)
data "secretsmanager_pam_database" "mysql_by_uid" {
  path = "cMS_-lYDeLs07A--rKMlNw"  # Replace with your record UID
}

# Example 2: Read PAM Database by title (errors if multiple records have same title)
data "secretsmanager_pam_database" "mysql_by_title" {
  title = "GatewayTest - MySQL Database"  # Replace with your record title
}

# Output the PAM Database data
output "db_hostname" {
  value = data.secretsmanager_pam_database.mysql_by_uid.pam_hostname[0].value[0].hostname
}

output "db_port" {
  value = data.secretsmanager_pam_database.mysql_by_uid.pam_hostname[0].value[0].port
}

output "db_type" {
  value = data.secretsmanager_pam_database.mysql_by_uid.database_type
}

output "db_use_ssl" {
  value = try(data.secretsmanager_pam_database.mysql_by_uid.use_ssl[0].value[0], false)
}

# Access pamSettings as JSON
output "db_pam_settings" {
  value     = jsondecode(data.secretsmanager_pam_database.mysql_by_uid.pam_settings)
  sensitive = true
}

# Example: Extract specific settings from pamSettings
locals {
  db_settings       = jsondecode(data.secretsmanager_pam_database.mysql_by_uid.pam_settings)
  db_protocol       = try(local.db_settings[0].connection[0].protocol, "unknown")
  db_name           = try(local.db_settings[0].connection[0].database, "")
  allow_supply_user = try(local.db_settings[0].connection[0].allowSupplyUser, false)
}

output "db_protocol" {
  value = local.db_protocol
}

output "db_name" {
  value = local.db_name
}

output "db_allow_supply_user" {
  value = local.allow_supply_user
}

# Example: Access cloud database metadata
output "cloud_db_info" {
  value = {
    database_id = try(data.secretsmanager_pam_database.mysql_by_uid.database_id[0].value, "")
    provider    = try(data.secretsmanager_pam_database.mysql_by_uid.provider_group[0].value, "")
    region      = try(data.secretsmanager_pam_database.mysql_by_uid.provider_region[0].value, "")
  }
}
