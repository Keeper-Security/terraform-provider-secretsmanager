terraform {
  required_version = ">= 1.0.0"
  required_providers {
    keeper = {
      source  = "github.com/keeper-security/keeper"
      version = ">= 0.1.0"
    }
    local = {
      source = "hashicorp/local"
      version = "2.1.0"
    }
  }
}

provider "local" { }
provider "keeper" {
  credential = "<CREDENTIAL>"
  # credential = file("~/.keeper/credential")
}

data "keeper_secret_encrypted_notes" "my_notes" {
  path        = "<record UID>"
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
UID:    ${ data.keeper_secret_encrypted_notes.my_notes.path }
Type:   ${ data.keeper_secret_encrypted_notes.my_notes.type }
Title:  ${ data.keeper_secret_encrypted_notes.my_notes.title }
Notes:  ${ data.keeper_secret_encrypted_notes.my_notes.notes }
======

Notes:  ${ data.keeper_secret_encrypted_notes.my_notes.note }

Date:   %{ if data.keeper_secret_encrypted_notes.my_notes.date != null ~}${ data.keeper_secret_encrypted_notes.my_notes.date }%{ endif ~}


FileRefs:
---------
%{ for fr in data.keeper_secret_encrypted_notes.my_notes.file_ref ~}
UID:      ${ fr.uid }
Title:    ${ fr.title }
Name:     ${ fr.name }
Type:     ${ fr.type }
Size:     ${ fr.size }
Last Modified:  ${ fr.last_modified }
URL:            ${ fr.url }

Content/Base64: ${ fr.content_base64 }


%{ endfor ~}
EOT
}

output "notes" {
  value = data.keeper_secret_encrypted_notes.my_notes.note
  sensitive = true
}
output "date" {
  value = data.keeper_secret_encrypted_notes.my_notes.date
}
