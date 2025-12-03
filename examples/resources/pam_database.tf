terraform {
  required_providers {
    secretsmanager = {
      source  = "keeper-security/secretsmanager"
      version = ">= 1.1.8"
    }
    local = {
      source = "hashicorp/local"
      version = "2.1.0"
    }
  }
}

provider "local" { }
provider "secretsmanager" {
  credential = "<CREDENTIAL>"
  # credential = file("~/.keeper/credential")
}

# Example 1: PostgreSQL Database
resource "secretsmanager_pam_database" "postgres_prod" {
  folder_uid = "<folder UID>"
  title = "Production PostgreSQL"
  notes = "Main production database cluster"

  pam_hostname {
    value {
      hostname = "postgres.prod.example.com"
      port = "5432"
    }
  }

  # PostgreSQL-specific connection settings
  pam_settings = jsonencode([{
    connection = [{
      protocol = "postgresql"
      port = "5432"
      recordingIncludeKeys = true
      allowSupplyUser = false
      database = "production"
      allowSupplyHost = false
    }]
  }])

  database_type = "postgresql"

  use_ssl {
    label = "useSSL"
    value = [true]
  }

  login {
    label = "DB Username"
    required = true
    value = "postgres_admin"
  }

  password {
    label = "DB Password"
    required = true
    privacy_screen = true
    enforce_generation = true
    generate = "yes"
    complexity {
      length = 32
      caps = 8
      lowercase = 8
      digits = 8
      special = 8
    }
  }

  # Optional: Connect to specific database
  connect_database {
    label = "Default Database"
    value = ["production"]
  }
}

# Example 2: MySQL Database
resource "secretsmanager_pam_database" "mysql_staging" {
  folder_uid = "<folder UID>"
  title = "Staging MySQL"
  notes = "Staging database for testing"

  pam_hostname {
    value {
      hostname = "mysql.staging.example.com"
      port = "3306"
    }
  }

  # MySQL-specific connection settings
  pam_settings = jsonencode([{
    connection = [{
      protocol = "mysql"
      port = "3306"
      recordingIncludeKeys = false
      allowSupplyUser = true
      database = "staging"
    }]
  }])

  database_type = "mysql"

  use_ssl {
    value = [false]
  }

  login {
    value = "mysql_user"
  }

  password {
    value = "Staging!P@ss123"
    privacy_screen = true
  }
}

# Example 3: AWS RDS PostgreSQL with metadata
resource "secretsmanager_pam_database" "aws_rds_postgres" {
  folder_uid = "<folder UID>"
  title = "AWS RDS PostgreSQL"
  notes = "Production RDS instance"

  pam_hostname {
    value {
      hostname = "mydb.cluster-abc123.us-east-1.rds.amazonaws.com"
      port = "5432"
    }
  }

  pam_settings = jsonencode([{
    connection = [{
      protocol = "postgresql"
      port = "5432"
      recordingIncludeKeys = true
      database = "appdb"
    }]
  }])

  database_type = "postgresql"

  use_ssl {
    value = [true]
  }

  login {
    value = "dbadmin"
  }

  password {
    enforce_generation = true
    generate = "yes"
    complexity {
      length = 32
      caps = 8
      lowercase = 8
      digits = 8
      special = 8
    }
  }

  # Database ID (RDS instance identifier)
  database_id {
    label = "RDS Instance ID"
    value = ["mydb-cluster"]
  }

  # Cloud provider metadata
  provider_group {
    label = "Provider"
    value = ["AWS"]
  }

  provider_region {
    label = "Region"
    value = ["us-east-1"]
  }
}

# Example 4: MongoDB with port forwarding
resource "secretsmanager_pam_database" "mongodb_dev" {
  folder_uid = "<folder UID>"
  title = "Development MongoDB"
  notes = "Local MongoDB instance for development"

  pam_hostname {
    value {
      hostname = "mongodb.dev.local"
      port = "27017"
    }
  }

  pam_settings = jsonencode([{
    connection = [{
      protocol = "mongodb"
      port = "27017"
      database = "devdb"
    }]
    portForward = [{
      port = "27018"
      reusePort = true
    }]
  }])

  database_type = "mongodb"

  login {
    value = "mongo_dev"
  }

  password {
    value = "DevP@ssw0rd123"
  }
}

# Output the created database records
output "postgres_prod_uid" {
  value = secretsmanager_pam_database.postgres_prod.uid
}

output "postgres_hostname" {
  value = secretsmanager_pam_database.postgres_prod.pam_hostname[0].hostname
}

output "mysql_staging_uid" {
  value = secretsmanager_pam_database.mysql_staging.uid
}

output "aws_rds_database_id" {
  value = secretsmanager_pam_database.aws_rds_postgres.database_id[0].value[0]
}
