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

resource "secretsmanager_contact" "my_contact" {
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
  company {
    label = "My Company"
    required = true
    privacy_screen = true
    value = "My Company"
  }
  email {
    label = "My Email"
    required = true
    privacy_screen = true
    value = "My Email"
  }
  phone {
    label = "My Phone"
    required = true
    privacy_screen = true
    value {
      region = "US"
      number = "202-555-0130"
      ext = "9987"
      type = "Work"
    }
  }
  address_ref {
    label = "My Address Ref"
    required = true
    privacy_screen = true
    value = "<address ref UID>"
  }
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
FUID:   ${ secretsmanager_contact.my_contact.folder_uid }
UID:    ${ secretsmanager_contact.my_contact.uid }
Type:   ${ secretsmanager_contact.my_contact.type }
Title:  ${ secretsmanager_contact.my_contact.title }
Notes:  ${ secretsmanager_contact.my_contact.notes }
======

Name:
-----
%{ for n in secretsmanager_contact.my_contact.name ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
First Name:  ${ n.value.0.first }
Middle Name: ${ n.value.0.middle }
Last Name:   ${ n.value.0.last }
%{ endfor }

Company:
--------
%{ for n in secretsmanager_contact.my_contact.company ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Value:    ${ n.value }
%{ endfor }

E-mail:
-------
%{ for n in secretsmanager_contact.my_contact.email ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Value:    ${ n.value }
%{ endfor }

Phone:
------
%{ for n in secretsmanager_contact.my_contact.phone ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Region: ${ n.value.0.region }
Number: ${ n.value.0.number }
Ext.:   ${ n.value.0.ext }
Type:   ${ n.value.0.type }
%{ endfor }

Address Ref:
------------
%{ for n in secretsmanager_contact.my_contact.address_ref ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Value:    ${ n.value }
%{ endfor }

EOT
}

output "record_uid" {
  value = secretsmanager_contact.my_contact.uid
}
output "record_title" {
  value = secretsmanager_contact.my_contact.title
}
