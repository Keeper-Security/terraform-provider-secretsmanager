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

data "secretsmanager_ssh_keys" "my_ssh_keys" {
  path        = "<record UID>"
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
UID:    ${ data.secretsmanager_ssh_keys.my_ssh_keys.path }
Type:   ${ data.secretsmanager_ssh_keys.my_ssh_keys.type }
Title:  ${ data.secretsmanager_ssh_keys.my_ssh_keys.title }
Notes:  ${ data.secretsmanager_ssh_keys.my_ssh_keys.notes }
======

Login:            ${ data.secretsmanager_ssh_keys.my_ssh_keys.login }

Key Pair:
---------
%{ for k in data.secretsmanager_ssh_keys.my_ssh_keys.key_pair }
Public Key:  ${ k.public_key }
Private Key: ${ k.private_key }
%{ endfor }

Passphrase:       ${ data.secretsmanager_ssh_keys.my_ssh_keys.passphrase }

Host: %{ for h in data.secretsmanager_ssh_keys.my_ssh_keys.host ~}Hostname:   ${ h.host_name }   Port:  ${ h.port } %{ endfor }


FileRefs:
---------
%{ for fr in data.secretsmanager_ssh_keys.my_ssh_keys.file_ref ~}
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

output "login" {
  value = data.secretsmanager_ssh_keys.my_ssh_keys.login
}
