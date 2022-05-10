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

resource "secretsmanager_driver_license" "my_driver_license" {
  folder_uid = "<folder UID>"
  title = "My Title"
  notes = "My Notes"
  driver_license_number {
    label = "My Driver License"
    required = true
    privacy_screen = true
    value = "My Driver License# 1234"
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
  expiration_date {
    label = "Driver License Expiration Date"
    required = true
    privacy_screen = true
    value = 21651186276
    # unix time in milliseconds
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
UID:    ${ secretsmanager_driver_license.my_driver_license.path }
Type:   ${ secretsmanager_driver_license.my_driver_license.type }
Title:  ${ secretsmanager_driver_license.my_driver_license.title }
Notes:  ${ secretsmanager_driver_license.my_driver_license.notes }
======

DL Number:  ${ secretsmanager_driver_license.my_driver_license.driver_license_number }

Name:
-----
%{ for n in secretsmanager_driver_license.my_driver_license.name ~}
First Name:   ${ n.first }
Midlle Name:  ${ n.middle }
Last Name:    ${ n.last }

%{ endfor ~}

Birth Date:       %{ if secretsmanager_driver_license.my_driver_license.birth_date != null ~}${ secretsmanager_driver_license.my_driver_license.birth_date }%{ endif ~}

Expiration Date:  %{ if secretsmanager_driver_license.my_driver_license.expiration_date != null ~}${ secretsmanager_driver_license.my_driver_license.expiration_date }%{ endif ~}


AddressRefs:
------------
%{ if secretsmanager_driver_license.my_driver_license.address_ref != null }
%{ for a in secretsmanager_driver_license.my_driver_license.address_ref ~}
UID:      ${ a.uid }
Street1:  ${ a.street1 }
Street2:  ${ a.street2 }
City:     ${ a.city }
State:    ${ a.state }
Zip:      ${ a.zip }
Country:  ${ a.country }

%{ endfor ~}
%{ endif }

FileRefs:
---------
%{ for fr in secretsmanager_driver_license.my_driver_license.file_ref ~}
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

output "driver_license_number" {
  value = secretsmanager_driver_license.my_driver_license.driver_license_number
}
output "expiration_date" {
  value = secretsmanager_driver_license.my_driver_license.expiration_date
}
