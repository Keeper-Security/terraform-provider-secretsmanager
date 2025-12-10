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

# Example 1: Active Directory PAM Directory
resource "secretsmanager_pam_directory" "active_directory" {
  folder_uid = "<folder UID>"
  title = "Corporate Active Directory"
  notes = "Main AD server for authentication"

  pam_hostname {
    value {
      hostname = "ad.corp.example.com"
      port = "636"  # LDAPS port
    }
  }

  # Active Directory connection settings
  pam_settings = jsonencode([{
    connection = [{
      protocol = "ldaps"
      port = "636"
      recordingIncludeKeys = false
    }]
  }])

  directory_type = "Active Directory"

  distinguished_name {
    label = "Base DN"
    value = ["DC=corp,DC=example,DC=com"]
  }

  use_ssl {
    label = "Use SSL"
    value = [true]
  }
}

# Example 2: OpenLDAP PAM Directory
resource "secretsmanager_pam_directory" "openldap" {
  folder_uid = "<folder UID>"
  title = "Development OpenLDAP"
  notes = "OpenLDAP directory for dev environment"

  pam_hostname {
    value {
      hostname = "ldap.dev.example.com"
      port = "389"
    }
  }

  # OpenLDAP connection settings
  pam_settings = jsonencode([{
    connection = [{
      protocol = "ldap"
      port = "389"
      recordingIncludeKeys = false
    }]
  }])

  directory_type = "OpenLDAP"

  distinguished_name {
    label = "Base DN"
    value = ["dc=dev,dc=example,dc=com"]
  }

  use_ssl {
    value = [false]  # Plain LDAP for dev
  }
}

# Example 3: Active Directory with LDAPS and custom settings
resource "secretsmanager_pam_directory" "secure_ad" {
  folder_uid = "<folder UID>"
  title = "Secure Active Directory"
  notes = "Production AD with LDAPS and SSL verification"

  pam_hostname {
    value {
      hostname = "secure-ad.prod.example.com"
      port = "636"
    }
  }

  pam_settings = jsonencode([{
    connection = [{
      protocol = "ldaps"
      port = "636"
      recordingIncludeKeys = true
      allowSupplyUser = false
    }]
    portForward = [{
      port = "1636"
      reusePort = true
    }]
  }])

  directory_type = "Active Directory"

  distinguished_name {
    label = "Search Base DN"
    value = ["OU=Users,DC=prod,DC=example,DC=com"]
  }

  use_ssl {
    label = "Enforce SSL"
    value = [true]
  }

  # Optional: Rotation scripts
  # rotation_scripts {
  #   value = [{
  #     command = "pwsh -File /opt/keeper/scripts/rotate-ad-password.ps1"
  #   }]
  # }
}

# Output the created directory records
output "ad_server_uid" {
  value = secretsmanager_pam_directory.active_directory.uid
}

output "ad_server_hostname" {
  value = secretsmanager_pam_directory.active_directory.pam_hostname[0].hostname
}

output "openldap_uid" {
  value = secretsmanager_pam_directory.openldap.uid
}

output "secure_ad_directory_type" {
  value = secretsmanager_pam_directory.secure_ad.directory_type
}
