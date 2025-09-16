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

resource "secretsmanager_ssn_card" "my_ssn_card" {
  folder_uid = "<folder UID>"
  title = "My Title"
  notes = "My Notes"
  identity_number {
    label = "My ID Number"
    required = true
    privacy_screen = true
    value = "My ID# 1234"
  }
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
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
FUID:   ${ secretsmanager_ssn_card.my_ssn_card.folder_uid }
UID:    ${ secretsmanager_ssn_card.my_ssn_card.uid }
Type:   ${ secretsmanager_ssn_card.my_ssn_card.type }
Title:  ${ secretsmanager_ssn_card.my_ssn_card.title }
Notes:  ${ secretsmanager_ssn_card.my_ssn_card.notes }
======

Identity Number:
----------------
%{ for n in secretsmanager_ssn_card.my_ssn_card.identity_number ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Value:    ${ n.value }
%{ endfor }

Name:
-----
%{ for n in secretsmanager_ssn_card.my_ssn_card.name ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
First Name:  ${ n.value.0.first }
Middle Name: ${ n.value.0.middle }
Last Name:   ${ n.value.0.last }
%{ endfor }

EOT
}

output "record_uid" {
  value = secretsmanager_ssn_card.my_ssn_card.uid
}
output "record_title" {
  value = secretsmanager_ssn_card.my_ssn_card.title
}
