terraform {
  required_providers {
    secretsmanager = {
      source  = "keeper-security/secretsmanager"
      version = ">= 1.3.0"
    }
  }
}

provider "secretsmanager" {
  # credential = "<CREDENTIAL>"
  # can also be set via KEEPER_CREDENTIAL env variable
}

# Read a PAM Remote Browser record as ephemeral (never stored in state)
ephemeral "secretsmanager_pam_remote_browser" "browser" {
  path = "<record UID>"
}
