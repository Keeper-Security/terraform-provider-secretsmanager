terraform {
  required_providers {
    secretsmanager = {
      source  = "keeper-security/secretsmanager"
      version = ">= 1.1.2"
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

resource "secretsmanager_software_license" "my_software_license" {
  folder_uid = "<folder UID>"
  title = "My Title"
  notes = "My Notes"
  license_number {
    label = "My License"
    required = true
    privacy_screen = true
    value = "My License# 1234"
  }
  activation_date {
    label = "License Activation Date"
    required = true
    privacy_screen = true
    value = 1651186276
    # unix time in milliseconds
  }
  expiration_date {
    label = "License Expiration Date"
    required = true
    privacy_screen = true
    value = 21651186276
    # unix time in milliseconds
  }
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
FUID:   ${ secretsmanager_software_license.my_software_license.folder_uid }
UID:    ${ secretsmanager_software_license.my_software_license.uid }
Type:   ${ secretsmanager_software_license.my_software_license.type }
Title:  ${ secretsmanager_software_license.my_software_license.title }
Notes:  ${ secretsmanager_software_license.my_software_license.notes }
======

License Number:
---------------
%{ for n in secretsmanager_software_license.my_software_license.license_number ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Value:    ${ n.value }
%{ endfor }

Activation Date:
----------------
%{ for n in secretsmanager_software_license.my_software_license.activation_date ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Value:    ${ n.value }
%{ endfor }

Expiration Date:
----------------
%{ for n in secretsmanager_software_license.my_software_license.expiration_date ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Value:    ${ n.value }
%{ endfor }

EOT
}

output "record_uid" {
  value = secretsmanager_software_license.my_software_license.uid
}
output "record_title" {
  value = secretsmanager_software_license.my_software_license.title
}
