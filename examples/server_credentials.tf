terraform {
  required_providers {
    secretsmanager = {
      source  = "keeper-security/secretsmanager"
      version = ">= 1.0.0"
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

data "secretsmanager_server_credentials" "my_server_creds" {
  path        = "<record UID>"
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
UID:    ${ data.secretsmanager_server_credentials.my_server_creds.path }
Type:   ${ data.secretsmanager_server_credentials.my_server_creds.type }
Title:  ${ data.secretsmanager_server_credentials.my_server_creds.title }
Notes:  ${ data.secretsmanager_server_credentials.my_server_creds.notes }
======

Login:    ${ data.secretsmanager_server_credentials.my_server_creds.login }
Password: ${ data.secretsmanager_server_credentials.my_server_creds.password }

Host: %{ for h in data.secretsmanager_server_credentials.my_server_creds.host ~}Hostname:   ${ h.host_name }   Port:  ${ h.port } %{ endfor }


FileRefs:
---------
%{ for fr in data.secretsmanager_server_credentials.my_server_creds.file_ref ~}
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

output "login" {
  value = data.secretsmanager_server_credentials.my_server_creds.login
}
