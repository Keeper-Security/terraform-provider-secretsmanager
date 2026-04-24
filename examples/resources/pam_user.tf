terraform {
  required_providers {
    secretsmanager = {
      source  = "keeper-security/secretsmanager"
      version = ">= 1.3.0"
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

# Example 1: PAM User with generated password and LDAP distinguished name
resource "secretsmanager_pam_user" "db_admin" {
  folder_uid = "<folder UID>"
  title      = "Database Admin User"
  notes      = "Production database administrator account"

  login {
    label    = "Username"
    required = true
    value    = "dbadmin"
  }

  password {
    label              = "Password"
    required           = true
    privacy_screen     = true
    enforce_generation = true
    generate           = "yes"
    complexity {
      length    = 32
      caps      = 8
      lowercase = 8
      digits    = 8
      special   = 8
    }
  }

  distinguished_name {
    label = "Distinguished Name"
    value = "CN=dbadmin,OU=Database Admins,DC=example,DC=com"
  }

  connect_database {
    label = "Connect Database"
    value = "production_db"
  }

  managed {
    label = "Managed"
    value = true
  }

  # Custom fields — attach arbitrary typed data to the record
  custom {
    type  = "text"
    label = "Environment"
    value = "production"
  }

}

# Example 2: PAM User with rotation scripts
resource "secretsmanager_pam_user" "service_account" {
  folder_uid = "<folder UID>"
  title      = "Service Account User"
  notes      = "Automated service account with rotation"

  login {
    value = "svc_automation"
  }

  password {
    value          = "CurrentP@ssw0rd123!"
    privacy_screen = true
  }

  rotation_scripts {
    label = "Rotation Scripts"
    value {
      file_ref = "fileRef123"
      command  = "rotate"
    }
  }

  connect_database {
    value = "app_database"
  }

  managed {
    value = true
  }
}

# Example 3: PAM User with TOTP/2FA
resource "secretsmanager_pam_user" "ldap_admin" {
  folder_uid = "<folder UID>"
  title      = "LDAP Admin User"
  notes      = "LDAP directory administrator with 2FA"

  login {
    value = "ldapadmin"
  }

  password {
    enforce_generation = true
    generate           = "yes"
    complexity {
      length    = 24
      caps      = 6
      lowercase = 6
      digits    = 6
      special   = 6
    }
  }

  distinguished_name {
    value = "CN=ldapadmin,OU=IT,OU=Admins,DC=corp,DC=example,DC=com"
  }

  totp {
    label = "Two-Factor Code"
    value = "otpauth://totp/Example:ldapadmin?secret=JBSWY3DPEHPK3PXP&issuer=Example"
  }

  managed {
    value = false
  }
}

# Example 4: Simple PAM User with minimal configuration
resource "secretsmanager_pam_user" "readonly_user" {
  folder_uid = "<folder UID>"
  title      = "Read-Only Database User"

  login {
    value = "readonly"
  }

  password {
    value = "ReadOnlyP@ss123"
  }

  connect_database {
    value = "analytics_db"
  }
}

# Output the created user records
output "db_admin_uid" {
  value = secretsmanager_pam_user.db_admin.uid
}

output "db_admin_login" {
  value = secretsmanager_pam_user.db_admin.login[0].value
}

output "service_account_uid" {
  value = secretsmanager_pam_user.service_account.uid
}

output "ldap_admin_dn" {
  value = secretsmanager_pam_user.ldap_admin.distinguished_name[0].value
}
