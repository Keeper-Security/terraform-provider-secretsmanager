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

resource "secretsmanager_passport" "my_passport" {
  folder_uid = "<folder UID>"
  title = "My Title"
  notes = "My Notes"
  passport_number {
    label = "My Passport"
    required = true
    privacy_screen = true
    value = "Passport# 1234"
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
  birth_date {
    label = "Birth Date"
    required = true
    privacy_screen = true
    value = 1651186276
    # unix time in milliseconds
  }
  date_issued {
    label = "Date Issued"
    required = true
    privacy_screen = true
    value = 4651186276
    # unix time in milliseconds
  }
  expiration_date {
    label = "Passport Expiration Date"
    required = true
    privacy_screen = true
    value = 21651186276
    # unix time in milliseconds
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
FUID:   ${ secretsmanager_passport.my_passport.folder_uid }
UID:    ${ secretsmanager_passport.my_passport.uid }
Type:   ${ secretsmanager_passport.my_passport.type }
Title:  ${ secretsmanager_passport.my_passport.title }
Notes:  ${ secretsmanager_passport.my_passport.notes }
======

Passport Number:
----------------
%{ for n in secretsmanager_passport.my_passport.passport_number ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Value:    ${ n.value }
%{ endfor }

Name:
-----
%{ for n in secretsmanager_passport.my_passport.name ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
First Name:     ${ n.value.0.first }
Middle Name:    ${ n.value.0.middle }
Last Name:      ${ n.value.0.last }
%{ endfor }

Birth Date:
-----------
%{ for n in secretsmanager_passport.my_passport.birth_date ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Value:    ${ n.value }
%{ endfor }

Expiration Date:
----------------
%{ for n in secretsmanager_passport.my_passport.expiration_date ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Value:    ${ n.value }
%{ endfor }

Date Issued:
------------
%{ for n in secretsmanager_passport.my_passport.date_issued ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Value:    ${ n.value }
%{ endfor }

Password:
---------
%{ for n in secretsmanager_passport.my_passport.password ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen:     ${ n.privacy_screen }
Enforce Generation: ${ n.enforce_generation }
Generate: %{ if n.generate != null }${n.generate}%{ endif }
Complexity: Length = ${ n.complexity.0.length }
Value:    ${ n.value }
%{ endfor }

Address Ref:
------------
%{ for n in secretsmanager_passport.my_passport.address_ref ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Value:    ${ n.value }
%{ endfor }

EOT
}

output "record_uid" {
  value = secretsmanager_passport.my_passport.uid
}
output "record_title" {
  value = secretsmanager_passport.my_passport.title
}
