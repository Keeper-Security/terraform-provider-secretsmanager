terraform {
  required_providers {
    secretsmanager = {
      source  = "keeper-security/secretsmanager"
      version = ">= 1.1.5"
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

data "secretsmanager_folder" "my_folder" {
  name = "<Folder Name>"
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
UID:       ${ data.secretsmanager_folder.my_folder.uid }
Name:      ${ data.secretsmanager_folder.my_folder.name }
ParentUID: ${ data.secretsmanager_folder.my_folder.parent_uid }
Shared:    ${ data.secretsmanager_folder.my_folder.shared }
EOT
}

output "my_folder_uid" {
  value = data.secretsmanager_folder.my_folder.uid
}
output "my_folder_name" {
  value = data.secretsmanager_folder.my_folder.name
}
