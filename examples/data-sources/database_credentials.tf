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

data "secretsmanager_database_credentials" "my_db_creds" {
  path        = "<record UID>"
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
UID:    ${ data.secretsmanager_database_credentials.my_db_creds.path }
Type:   ${ data.secretsmanager_database_credentials.my_db_creds.type }
Title:  ${ data.secretsmanager_database_credentials.my_db_creds.title }
Notes:  ${ data.secretsmanager_database_credentials.my_db_creds.notes }
======

DB Type:  ${ data.secretsmanager_database_credentials.my_db_creds.db_type }
Login:    ${ data.secretsmanager_database_credentials.my_db_creds.login }
Password: ${ data.secretsmanager_database_credentials.my_db_creds.password }

Host: %{ for h in data.secretsmanager_database_credentials.my_db_creds.host ~}Hostname:   ${ h.host_name }   Port:  ${ h.port } %{ endfor }


FileRefs:
---------
%{ for fr in data.secretsmanager_database_credentials.my_db_creds.file_ref ~}
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

output "db_type" {
  value = data.secretsmanager_database_credentials.my_db_creds.db_type
}
output "login" {
  value = data.secretsmanager_database_credentials.my_db_creds.login
}
