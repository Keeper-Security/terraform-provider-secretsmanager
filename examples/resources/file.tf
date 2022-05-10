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

resource "secretsmanager_file" "my_files" {
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
FUID:   ${ secretsmanager_file.my_files.folder_uid }
UID:    ${ secretsmanager_file.my_files.uid }
Type:   ${ secretsmanager_file.my_files.type }
Title:  ${ secretsmanager_file.my_files.title }
Notes:  ${ secretsmanager_file.my_files.notes }
======

FileRefs:
---------
%{ for fr in secretsmanager_file.my_files.file_ref ~}
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
  value = secretsmanager_file.my_files.uid
}
output "record_title" {
  value = secretsmanager_file.my_files.title
}

output "files" {
  value = secretsmanager_file.my_files.file_ref
}

output "file" {
  value = length(secretsmanager_file.my_files.file_ref.0.value) < 1 ? "" : textdecodebase64(secretsmanager_file.my_files.file_ref.0.value.0.content_base64, "UTF-8")
}
