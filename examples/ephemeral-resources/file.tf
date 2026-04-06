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

ephemeral "secretsmanager_file" "my_files" {
  path = "<record UID>"
}

output "file_count" {
  value     = length(ephemeral.secretsmanager_file.my_files.file_ref)
  ephemeral = true
}

output "title" {
  value     = ephemeral.secretsmanager_file.my_files.title
  ephemeral = true
}
