terraform {
  required_providers {
    secretsmanager = {
      source  = "keeper-security/secretsmanager"
      version = ">= 1.4.0"
    }
    local = {
      source  = "hashicorp/local"
      version = "2.1.0"
    }
  }
}

provider "local" {}
provider "secretsmanager" {
  credential = "<CREDENTIAL>"
  # credential = file("~/.keeper/credential")
}

data "secretsmanager_metadata" "db_server" {
  path = "<record UID>"
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
UID:        ${data.secretsmanager_metadata.db_server.uid}
Type:       ${data.secretsmanager_metadata.db_server.type}
Title:      ${data.secretsmanager_metadata.db_server.title}
Notes:      ${data.secretsmanager_metadata.db_server.notes}
Revision:   ${data.secretsmanager_metadata.db_server.revision}
FolderUID:  ${data.secretsmanager_metadata.db_server.folder_uid}
IsEditable: ${data.secretsmanager_metadata.db_server.is_editable}
EOT
}

output "db_server_revision" {
  value = data.secretsmanager_metadata.db_server.revision
}
