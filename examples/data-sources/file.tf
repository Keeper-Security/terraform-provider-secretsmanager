terraform {
  required_providers {
    secretsmanager = {
      source  = "keeper-security/secretsmanager"
      version = ">= 1.1.2"
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

data "secretsmanager_file" "my_files" {
  path        = "<record UID>"
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
UID:    ${ data.secretsmanager_file.my_files.path }
Type:   ${ data.secretsmanager_file.my_files.type }
Title:  ${ data.secretsmanager_file.my_files.title }
Notes:  ${ data.secretsmanager_file.my_files.notes }
======

FileRefs:
---------
%{ for fr in data.secretsmanager_file.my_files.file_ref ~}
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

output "file_count" {
  value = length(data.secretsmanager_file.my_files.file_ref.*)
  sensitive = true
}
