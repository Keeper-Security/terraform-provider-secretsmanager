terraform {
  required_providers {
    secretsmanager = {
      source  = "keeper-security/secretsmanager"
      version = ">= 1.1.7"
    }
    local = {
      source  = "hashicorp/local"
      version = "2.1.0"
    }
  }
}

provider "local" {}
provider "secretsmanager" {
  credential = "<CREDENTIAL>"
  # credential = file("~/.keeper/credential")
}

# Example 1: Fetch multiple records by UID
data "secretsmanager_records" "by_uids" {
  uids = [
    "<record_uid_1>",
    "<record_uid_2>",
    "<record_uid_3>"
  ]
}

# Example 2: Fetch records by title
data "secretsmanager_records" "by_titles" {
  titles = [
    "Production Database",
    "Staging Database"
  ]
}

# Example 3: Mix UIDs and titles
data "secretsmanager_records" "mixed" {
  uids = [
    "<record_uid_1>"
  ]
  titles = [
    "API Gateway Config"
  ]
}

# Example 4: Large batch for infrastructure secrets
locals {
  # In real usage, these would be actual UIDs for your infrastructure
  infrastructure_uids = [
    "<db_prod_uid>",
    "<db_staging_uid>",
    "<cache_redis_uid>",
    "<api_gateway_uid>",
    "<service_account_1_uid>",
    "<service_account_2_uid>"
    # Can include hundreds of UIDs here
  ]
}

data "secretsmanager_records" "infrastructure" {
  uids = local.infrastructure_uids
}

# Output examples showing different ways to access the data
output "all_record_titles" {
  value       = [for r in data.secretsmanager_records.by_uids.records : r.title]
  description = "List of all record titles"
}

output "first_record" {
  value       = data.secretsmanager_records.by_uids.records[0]
  sensitive   = true
  description = "Complete first record (by index)"
}

output "specific_record_by_uid" {
  value       = jsondecode(data.secretsmanager_records.by_uids.records_by_uid["<record_uid_1>"])
  sensitive   = true
  description = "Access specific record by UID from map"
}

# Example: Extract specific field values
locals {
  # Decode a record from the map
  db_record = jsondecode(
    data.secretsmanager_records.infrastructure.records_by_uid["<db_prod_uid>"]
  )
  
  # Extract username (assuming it's the first field)
  db_username = local.db_record.fields[0].value
  
  # Extract password (assuming it's the second field)
  db_password = sensitive(local.db_record.fields[1].value)
}

# Write output to file demonstrating batch efficiency
resource "local_file" "batch_results" {
  filename        = "${path.module}/batch_results.txt"
  file_permission = "0644"
  content         = <<-EOT
    Batch Fetch Results
    ===================
    
    Total records fetched: ${length(data.secretsmanager_records.by_uids.records)}
    
    This single data source replaced ${length(data.secretsmanager_records.by_uids.records)} individual API calls!
    
    Records by Title:
    ${join("\n", [for r in data.secretsmanager_records.by_uids.records : "- ${r.title} (${r.type})"])}
    
    Available UIDs in map:
    ${join("\n", [for uid, _ in data.secretsmanager_records.by_uids.records_by_uid : "- ${uid}"])}
  EOT
}

# Example: Using with other resources
# This shows how batch-fetched secrets can be used with AWS resources
/*
resource "aws_db_instance" "example" {
  identifier = "mydb-instance"
  engine     = "postgres"
  
  # Access credentials from batch-fetched records
  username = local.db_username
  password = local.db_password
  
  # ... other configuration
}

resource "aws_elasticache_cluster" "example" {
  cluster_id = "my-cache"
  
  # Access Redis password from batch fetch
  auth_token = sensitive(
    jsondecode(
      data.secretsmanager_records.infrastructure.records_by_uid["<cache_redis_uid>"]
    ).fields[0].value
  )
  
  # ... other configuration
}
*/
