terraform {
  required_providers {
    secretsmanager = {
      source  = "keeper-security/secretsmanager"
      version = ">= 1.1.1"
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

resource "secretsmanager_photo" "my_photos" {
  folder_uid = "<folder UID>"
  title = "My Title"
  notes = "My Notes"
  file_ref {
    value { uid = "<file1 UID>" }
    value { uid = "<file2 UID>" }
  }
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
FUID:   ${ secretsmanager_photo.my_photos.folder_uid }
UID:    ${ secretsmanager_photo.my_photos.uid }
Type:   ${ secretsmanager_photo.my_photos.type }
Title:  ${ secretsmanager_photo.my_photos.title }
Notes:  ${ secretsmanager_photo.my_photos.notes }
======

FileRefs:
---------
%{ for fr in secretsmanager_photo.my_photos.file_ref ~}
Type:     ${ fr.type }
Label:    ${ fr.label }
Required: ${ fr.required }

%{ for fv in fr.value ~}
UID:      ${ fv.uid }
Title:    ${ fv.title }
Name:     ${ fv.name }
Type:     ${ fv.type }
Size:     ${ fv.size }
Last Modified:  ${ fv.last_modified }

Content/Base64: ${ fv.content_base64 }

%{ endfor ~}
%{ endfor ~}

EOT
}

output "record_uid" {
  value = secretsmanager_photo.my_photos.uid
}
output "record_title" {
  value = secretsmanager_photo.my_photos.title
}

output "files" {
  value = secretsmanager_photo.my_photos.file_ref
}

output "file" {
  value = length(secretsmanager_photo.my_photos.file_ref.0.value) < 1 ? "" : textdecodebase64(secretsmanager_photo.my_photos.file_ref.0.value.0.content_base64, "UTF-8")
}
