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

resource "secretsmanager_ssh_keys" "my_ssh_keys" {
  folder_uid = "<folder UID>"
  title = "My Title"
  notes = "My Notes"
  login {
    label = "My Login"
    required = true
    privacy_screen = true
    value = "MyLogin"
  }
  passphrase {
    label = "My Pass"
    required = true
    privacy_screen = true
    enforce_generation = true
    #generate = "yes"
    complexity {
      length = 20
      caps = 5
      lowercase = 5
      digits = 5
      special = 5
    }
    value = "<SSH PASSPHRASE>"
  }
  host {
    label = "My Host"
    required = true
    privacy_screen = true
    value {
      host_name = "127.0.0.1"
      port = "22"
    }
  }
  key_pair {
    label = "My Keys"
    required = true
    privacy_screen = true
    value {
      public_key = "<PUBLIC KEY>"
      private_key = "<PRIVATE KEY>"
    }
  }
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
FUID:   ${ secretsmanager_ssh_keys.my_ssh_keys.folder_uid }
UID:    ${ secretsmanager_ssh_keys.my_ssh_keys.uid }
Type:   ${ secretsmanager_ssh_keys.my_ssh_keys.type }
Title:  ${ secretsmanager_ssh_keys.my_ssh_keys.title }
Notes:  ${ secretsmanager_ssh_keys.my_ssh_keys.notes }
======

Login:
------
%{ for n in secretsmanager_ssh_keys.my_ssh_keys.login ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Value:    ${ n.value }
%{ endfor }

Passphrase:
-----------
%{ for n in secretsmanager_ssh_keys.my_ssh_keys.passphrase ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Enforce Generation: ${ n.enforce_generation }
Generate: %{ if n.generate != null }${n.generate}%{ endif }
Complexity: Length = ${ n.complexity.0.length }
Value:    ${ n.value }
%{ endfor }

Host:
-----
%{ for n in secretsmanager_ssh_keys.my_ssh_keys.host ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Host Name: ${ n.value.0.host_name }
Port:      ${ n.value.0.port }
%{ endfor }

Key Pair:
---------
%{ for n in secretsmanager_ssh_keys.my_ssh_keys.key_pair ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Public Key:  ${ n.value.0.public_key }
Private Key: ${ n.value.0.private_key }
%{ endfor }

EOT
}

output "record_uid" {
  value = secretsmanager_ssh_keys.my_ssh_keys.uid
}
output "record_title" {
  value = secretsmanager_ssh_keys.my_ssh_keys.title
}
