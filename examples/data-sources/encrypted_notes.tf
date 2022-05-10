terraform {
  required_providers {
    secretsmanager = {
      source  = "keeper-security/secretsmanager"
      version = ">= 1.1.0"
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

data "secretsmanager_encrypted_notes" "my_notes" {
  path        = "<record UID>"
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
UID:    ${ data.secretsmanager_encrypted_notes.my_notes.path }
Type:   ${ data.secretsmanager_encrypted_notes.my_notes.type }
Title:  ${ data.secretsmanager_encrypted_notes.my_notes.title }
Notes:  ${ data.secretsmanager_encrypted_notes.my_notes.notes }
======

Notes:  ${ data.secretsmanager_encrypted_notes.my_notes.note }

Date:   %{ if data.secretsmanager_encrypted_notes.my_notes.date != null ~}${ data.secretsmanager_encrypted_notes.my_notes.date }%{ endif ~}


FileRefs:
---------
%{ for fr in data.secretsmanager_encrypted_notes.my_notes.file_ref ~}
UID:      ${ fr.uid }
Title:    ${ fr.title }
Name:     ${ fr.name }
Type:     ${ fr.type }
Size:     ${ fr.size }
Last Modified:  ${ fr.last_modified }

Content/Base64: ${ fr.content_base64 }


%{ endfor ~}
EOT
}

output "notes" {
  value = data.secretsmanager_encrypted_notes.my_notes.note
  sensitive = true
}
output "date" {
  value = data.secretsmanager_encrypted_notes.my_notes.date
}
