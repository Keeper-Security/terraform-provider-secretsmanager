terraform {
  required_providers {
    secretsmanager = {
      source  = "keeper-security/secretsmanager"
      version = ">= 1.1.7"
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

resource "secretsmanager_encrypted_notes" "my_encrypted_notes" {
  folder_uid = "<folder UID>"
  title = "My Title"
  notes = "My Notes"
  note {
    label = "My Note"
    required = true
    privacy_screen = true
    value = "My Note"
  }
  date {
    label = "Date"
    required = true
    privacy_screen = true
    value = 1651186276
    # unix time in milliseconds
  }
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
FUID:   ${ secretsmanager_encrypted_notes.my_encrypted_notes.folder_uid }
UID:    ${ secretsmanager_encrypted_notes.my_encrypted_notes.uid }
Type:   ${ secretsmanager_encrypted_notes.my_encrypted_notes.type }
Title:  ${ secretsmanager_encrypted_notes.my_encrypted_notes.title }
Notes:  ${ secretsmanager_encrypted_notes.my_encrypted_notes.notes }
======

Note:
-----
%{ for n in secretsmanager_encrypted_notes.my_encrypted_notes.note ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Value:    ${ n.value }
%{ endfor }

Date:
-----
%{ for n in secretsmanager_encrypted_notes.my_encrypted_notes.date ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Value:    ${ n.value }
%{ endfor }

EOT
}

output "record_uid" {
  value = secretsmanager_encrypted_notes.my_encrypted_notes.uid
}
output "record_title" {
  value = secretsmanager_encrypted_notes.my_encrypted_notes.title
}
