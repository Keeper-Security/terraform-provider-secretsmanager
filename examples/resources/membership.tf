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

resource "secretsmanager_membership" "my_membership" {
  folder_uid = "<folder UID>"
  title = "My Title"
  notes = "My Notes"
  account_number {
    label = "My Account"
    required = true
    privacy_screen = true
    value = "My Account# 1234"
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
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
FUID:   ${ secretsmanager_membership.my_membership.folder_uid }
UID:    ${ secretsmanager_membership.my_membership.uid }
Type:   ${ secretsmanager_membership.my_membership.type }
Title:  ${ secretsmanager_membership.my_membership.title }
Notes:  ${ secretsmanager_membership.my_membership.notes }
======

Account Number:
---------------
%{ for n in secretsmanager_membership.my_membership.account_number ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Value:    ${ n.value }
%{ endfor }

Name:
-----
%{ for n in secretsmanager_membership.my_membership.name ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
First Name:     ${ n.value.0.first }
Middle Name:    ${ n.value.0.middle }
Last Name:      ${ n.value.0.last }
%{ endfor }

Password:
---------
%{ for n in secretsmanager_membership.my_membership.password ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Enforce Generation: ${ n.enforce_generation }
Generate: %{ if n.generate != null }${n.generate}%{ endif }
Complexity: Length = ${ n.complexity.0.length }
Value:    ${ n.value }
%{ endfor }

EOT
}

output "record_uid" {
  value = secretsmanager_membership.my_membership.uid
}
output "record_title" {
  value = secretsmanager_membership.my_membership.title
}
