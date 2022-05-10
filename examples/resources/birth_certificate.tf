terraform {
  required_providers {
    secretsmanager = {
      source  = "keeper-security/secretsmanager"
      version = ">= 1.1.0"
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

resource "secretsmanager_birth_certificate" "my_birth_certificate" {
  folder_uid = "<folder UID>"
  title = "My Title"
  notes = "My Notes"
  name {
    label = "John"
    required = true
    privacy_screen = true
    value {
      first = "John"
      middle = "D"
      last = "Doe"
    }
  }
  birth_date {
    label = "Birth Date"
    required = true
    privacy_screen = true
    value = 1651186276
    # unix time in milliseconds
    # unix time seconds can be produced using time_static resource from hashicorp/time provider
  }
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content     = <<EOT
FUID:  ${ secretsmanager_birth_certificate.my_birth_certificate.folder_uid }
UID:   ${ secretsmanager_birth_certificate.my_birth_certificate.uid }
Type:  ${ secretsmanager_birth_certificate.my_birth_certificate.type }
Title: ${ secretsmanager_birth_certificate.my_birth_certificate.title }
Notes: ${ secretsmanager_birth_certificate.my_birth_certificate.notes }
======

Name:
-----
%{ for n in secretsmanager_birth_certificate.my_birth_certificate.name ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
First Name:  ${ n.value.0.first }
Middle Name: ${ n.value.0.middle }
Last Name:   ${ n.value.0.last }
%{ endfor }

Birth Date:
-----------
%{ for n in secretsmanager_birth_certificate.my_birth_certificate.birth_date ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Value:    ${ n.value }
%{ endfor }

EOT
}

output "record_uid" {
  value = secretsmanager_birth_certificate.my_birth_certificate.uid
}
output "record_title" {
  value = secretsmanager_birth_certificate.my_birth_certificate.title
}
