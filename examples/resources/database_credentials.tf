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

resource "secretsmanager_database_credentials" "my_database_credentials" {
  folder_uid = "<folder UID>"
  title = "My Title"
  notes = "My Notes"
  login {
    label = "My Login"
    required = true
    privacy_screen = true
    value = "MyLogin"
  }
  password {
    label = "My Pass"
    required = true
    privacy_screen = true
    enforce_generation = true
    generate = "yes"
    complexity {
      length = 20
      caps = 5
      lowercase = 5
      digits = 5
      special = 5
    }
    #value = "to_be_generated"
  }
  host {
    label = "My Host"
    required = true
    privacy_screen = true
    value {
      host_name = "127.0.0.1"
      port = "3306"
    }
  }
  db_type {
    label = "My DB Type"
    required = true
    privacy_screen = true
    value = "MySQL"
  }
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
FUID:   ${ secretsmanager_database_credentials.my_database_credentials.folder_uid }
UID:    ${ secretsmanager_database_credentials.my_database_credentials.uid }
Type:   ${ secretsmanager_database_credentials.my_database_credentials.type }
Title:  ${ secretsmanager_database_credentials.my_database_credentials.title }
Notes:  ${ secretsmanager_database_credentials.my_database_credentials.notes }
======

DB Type:
--------
%{ for n in secretsmanager_database_credentials.my_database_credentials.db_type ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Value:    ${ n.value }
%{ endfor }

Login:
------
%{ for n in secretsmanager_database_credentials.my_database_credentials.login ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Value:    ${ n.value }
%{ endfor }

Password:
---------
%{ for n in secretsmanager_database_credentials.my_database_credentials.password ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Enforce Generation:  ${ n.enforce_generation }
Generate:            %{ if n.generate != null }${n.generate}%{ endif }
Complexity: Length = ${ n.complexity.0.length }
Value:    ${ n.value }
%{ endfor }

Host:
-----
%{ for n in secretsmanager_database_credentials.my_database_credentials.host ~}
Type:      ${ n.type }
Label:     ${ n.label }
Required:  ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Host Name: ${ n.value.0.host_name }
Port:      ${ n.value.0.port }
%{ endfor }

EOT
}

output "record_uid" {
  value = secretsmanager_database_credentials.my_database_credentials.uid
}
output "record_title" {
  value = secretsmanager_database_credentials.my_database_credentials.title
}
