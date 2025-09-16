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

data "secretsmanager_folders" "my_folders" { }

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT

Folders:
--------
%{ for f in data.secretsmanager_folders.my_folders.folders ~}
UID:       ${ f.uid }
Name:      ${ f.name }
ParentUID: ${ f.parent_uid }
Shared:    ${ f.shared }

%{ endfor ~}
EOT
}

output "my_folders_count" {
  value = length(data.secretsmanager_folders.my_folders.folders)
}
