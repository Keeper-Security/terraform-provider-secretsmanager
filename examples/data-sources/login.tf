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

data "secretsmanager_login" "db_server" {
  path       = "<record UID>"
}

resource "local_file" "out" {
    filename        = "${path.module}/out.txt"
    file_permission = "0644"
    content         = <<EOT
UID:    ${ data.secretsmanager_login.db_server.path }
Type:   ${ data.secretsmanager_login.db_server.type }
Title:  ${ data.secretsmanager_login.db_server.title }
Notes:  ${ data.secretsmanager_login.db_server.notes }
======

Login:    ${ data.secretsmanager_login.db_server.login }
Password: ${ data.secretsmanager_login.db_server.password }
URL:      ${ data.secretsmanager_login.db_server.url }

TOTP:
-----
%{ for t in data.secretsmanager_login.db_server.totp ~}
URL:    ${ t.url }
Token:  ${ t.token }
TTL:    ${ t.ttl }

%{ endfor ~}

FileRefs:
---------
%{ for fr in data.secretsmanager_login.db_server.file_ref ~}
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

output "db_secret_login" {
  value = data.secretsmanager_login.db_server.login
}
