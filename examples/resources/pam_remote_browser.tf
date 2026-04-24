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

# Example 1: PAM Remote Browser with RBI URL
resource "secretsmanager_pam_remote_browser" "rbi_example" {
  folder_uid = "<folder UID>"
  title      = "Production RBI Session"
  notes      = "Remote Browser Isolation for secure web access"

  rbi_url {
    label = "RBI URL"
    value = "https://rbi.example.com/portal"
  }

  traffic_encryption_seed {
    label = "Encryption Seed"
    value = "seed-value-here"
  }

  pam_remote_browser_settings = jsonencode([{
    recordingEnabled = true
    sessionTimeout   = 3600
  }])

  # Optional TOTP for multi-factor authentication
  totp {
    value = "JBSWY3DPEBLW64TMMQ======"
  }

  # Custom fields — attach arbitrary typed data to the record
  custom {
    type  = "text"
    label = "Environment"
    value = "production"
  }
  custom {
    type  = "email"
    label = "Admin Contact"
    value = "admin@example.com"
  }
}

# Example 2: PAM Remote Browser with file reference
resource "secretsmanager_pam_remote_browser" "rbi_with_cert" {
  folder_uid = "<folder UID>"
  title      = "Staging RBI with Certificate"
  notes      = "Staging environment with client certificate"

  rbi_url {
    label = "RBI URL"
    value = "https://staging-rbi.example.com"
  }

  # File reference for client certificate
  file_ref {
    label = "Client Certificate"
    value = "<file UID>"
  }

  pam_remote_browser_settings = jsonencode([{
    recordingEnabled      = false
    sessionTimeout        = 1800
    requireClientCertificate = true
  }])

  custom {
    type  = "text"
    label = "Environment"
    value = "staging"
  }
}

# Output the created PAM Remote Browser records
output "rbi_example_uid" {
  value = secretsmanager_pam_remote_browser.rbi_example.uid
}

output "rbi_example_url" {
  value = secretsmanager_pam_remote_browser.rbi_example.rbi_url[0].value
}

output "rbi_with_cert_uid" {
  value = secretsmanager_pam_remote_browser.rbi_with_cert.uid
}
