terraform {
  required_providers {
    secretsmanager = {
      source  = "keeper-security/secretsmanager"
      version = ">= 1.1.7"
    }
  }
}

provider "secretsmanager" {
  credential = "<CREDENTIAL>"
  # credential = file("~/.keeper/credential")
}

# Example 1: Read PAM Database by path
data "secretsmanager_pam_database" "postgres_by_path" {
  path = "/Databases/Production/PostgreSQL"
}

# Example 2: Read PAM Database by title
data "secretsmanager_pam_database" "postgres_by_title" {
  title = "Production PostgreSQL"
}

# Output the PAM Database data
output "db_hostname" {
  value = data.secretsmanager_pam_database.postgres_by_path.pam_hostname[0].hostname
}

output "db_port" {
  value = data.secretsmanager_pam_database.postgres_by_path.pam_hostname[0].port
}

output "db_type" {
  value = data.secretsmanager_pam_database.postgres_by_path.database_type[0].value[0]
}

output "db_login" {
  value = data.secretsmanager_pam_database.postgres_by_path.login[0].value
}

output "db_password" {
  value     = data.secretsmanager_pam_database.postgres_by_path.password[0].value
  sensitive = true
}

output "db_use_ssl" {
  value = try(data.secretsmanager_pam_database.postgres_by_path.use_ssl[0].value[0], false)
}

# Access pamSettings as JSON
output "db_pam_settings" {
  value = jsondecode(data.secretsmanager_pam_database.postgres_by_path.pam_settings)
  sensitive = true
}

# Example: Extract specific settings from pamSettings
locals {
  db_settings = jsondecode(data.secretsmanager_pam_database.postgres_by_path.pam_settings)
  db_protocol = try(local.db_settings[0].connection[0].protocol, "unknown")
  db_name = try(local.db_settings[0].connection[0].database, "")
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

# Example: Build database connection string
locals {
  db_host = data.secretsmanager_pam_database.postgres_by_path.pam_hostname[0].hostname
  db_port = data.secretsmanager_pam_database.postgres_by_path.pam_hostname[0].port
  db_user = data.secretsmanager_pam_database.postgres_by_path.login[0].value
  db_pass = data.secretsmanager_pam_database.postgres_by_path.password[0].value
  db_ssl = try(data.secretsmanager_pam_database.postgres_by_path.use_ssl[0].value[0], false) ? "require" : "disable"

  # PostgreSQL connection string
  postgres_connection_string = "postgresql://${local.db_user}:${local.db_pass}@${local.db_host}:${local.db_port}/${local.db_name}?sslmode=${local.db_ssl}"
}

output "connection_string" {
  value     = local.postgres_connection_string
  sensitive = true
}

# Example: Access cloud database metadata
output "cloud_db_info" {
  value = {
    database_id = try(data.secretsmanager_pam_database.postgres_by_path.database_id[0].value[0], "")
    provider = try(data.secretsmanager_pam_database.postgres_by_path.provider_group[0].value[0], "")
    region = try(data.secretsmanager_pam_database.postgres_by_path.provider_region[0].value[0], "")
  }
}

# Example: Use in database provider configuration
# This would connect to the database using credentials from Keeper
/*
provider "postgresql" {
  host     = data.secretsmanager_pam_database.postgres_by_path.pam_hostname[0].hostname
  port     = tonumber(data.secretsmanager_pam_database.postgres_by_path.pam_hostname[0].port)
  username = data.secretsmanager_pam_database.postgres_by_path.login[0].value
  password = data.secretsmanager_pam_database.postgres_by_path.password[0].value
  sslmode  = try(data.secretsmanager_pam_database.postgres_by_path.use_ssl[0].value[0], false) ? "require" : "disable"
  database = local.db_name
}

resource "postgresql_database" "new_db" {
  name = "my_new_database"
}
*/
