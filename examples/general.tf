terraform {
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

data "keeper_secret_general" "db_server" {
  path       = "<record UID>"
}

resource "local_file" "out" {
    filename        = "${path.module}/out.txt"
    file_permission = "0644"
    content         = <<EOT
UID:    ${ data.keeper_secret_general.db_server.path }
Type:   ${ data.keeper_secret_general.db_server.type }
Title:  ${ data.keeper_secret_general.db_server.title }
Notes:  ${ data.keeper_secret_general.db_server.notes }
======

Login:    ${ data.keeper_secret_general.db_server.login }
Password: ${ data.keeper_secret_general.db_server.password }
URL:      ${ data.keeper_secret_general.db_server.url }

TOTP:
-----
%{ for t in data.keeper_secret_general.db_server.totp ~}
URL:    ${ t.url }
Token:  ${ t.token }
TTL:    ${ t.ttl }

%{ endfor ~}

FileRefs:
---------
%{ for fr in data.keeper_secret_general.db_server.file_ref ~}
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

output "db_secret_login" {
  value = data.keeper_secret_general.db_server.login
}
