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

ephemeral "secretsmanager_photo" "my_photos" {
  path = "<record UID>"
}

output "photo_count" {
  value     = length(ephemeral.secretsmanager_photo.my_photos.file_ref)
  ephemeral = true
}

output "title" {
  value     = ephemeral.secretsmanager_photo.my_photos.title
  ephemeral = true
}
