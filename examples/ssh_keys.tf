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

data "keeper_secret_ssh_keys" "my_ssh_keys" {
  path        = "<record UID>"
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
UID:    ${ data.keeper_secret_ssh_keys.my_ssh_keys.path }
Type:   ${ data.keeper_secret_ssh_keys.my_ssh_keys.type }
Title:  ${ data.keeper_secret_ssh_keys.my_ssh_keys.title }
Notes:  ${ data.keeper_secret_ssh_keys.my_ssh_keys.notes }
======

Login:            ${ data.keeper_secret_ssh_keys.my_ssh_keys.login }

Key Pair:
---------
%{ for k in data.keeper_secret_ssh_keys.my_ssh_keys.key_pair }
Public Key:  ${ k.public_key }
Private Key: ${ k.private_key }
%{ endfor }

Passphrase:       ${ data.keeper_secret_ssh_keys.my_ssh_keys.passphrase }

Host: %{ for h in data.keeper_secret_ssh_keys.my_ssh_keys.host ~}Hostname:   ${ h.host_name }   Port:  ${ h.port } %{ endfor }


FileRefs:
---------
%{ for fr in data.keeper_secret_ssh_keys.my_ssh_keys.file_ref ~}
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
  value = data.keeper_secret_ssh_keys.my_ssh_keys.login
}
