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

data "secretsmanager_driver_license" "my_license" {
  path        = "<record UID>"
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
UID:    ${ data.secretsmanager_driver_license.my_license.path }
Type:   ${ data.secretsmanager_driver_license.my_license.type }
Title:  ${ data.secretsmanager_driver_license.my_license.title }
Notes:  ${ data.secretsmanager_driver_license.my_license.notes }
======

DL Number:  ${ data.secretsmanager_driver_license.my_license.driver_license_number }

Name:
-----
%{ for n in data.secretsmanager_driver_license.my_license.name ~}
First Name:   ${ n.first }
Midlle Name:  ${ n.middle }
Last Name:    ${ n.last }

%{ endfor ~}

Birth Date:       %{ if data.secretsmanager_driver_license.my_license.birth_date != null ~}${ data.secretsmanager_driver_license.my_license.birth_date }%{ endif ~}

Expiration Date:  %{ if data.secretsmanager_driver_license.my_license.expiration_date != null ~}${ data.secretsmanager_driver_license.my_license.expiration_date }%{ endif ~}


AddressRefs:
--------
%{ if data.secretsmanager_driver_license.my_license.address_ref != null }
%{ for a in data.secretsmanager_driver_license.my_license.address_ref ~}
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
%{ for fr in data.secretsmanager_driver_license.my_license.file_ref ~}
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
  value = data.secretsmanager_driver_license.my_license.driver_license_number
}
output "expiration_date" {
  value = data.secretsmanager_driver_license.my_license.expiration_date
}
