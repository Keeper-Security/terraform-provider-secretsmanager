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

data "secretsmanager_software_license" "my_license" {
  path        = "<record UID>"
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
UID:    ${ data.secretsmanager_software_license.my_license.path }
Type:   ${ data.secretsmanager_software_license.my_license.type }
Title:  ${ data.secretsmanager_software_license.my_license.title }
Notes:  ${ data.secretsmanager_software_license.my_license.notes }
======

License#:  ${ data.secretsmanager_software_license.my_license.license_number }
Activation Date:  %{ if data.secretsmanager_software_license.my_license.activation_date != null ~}${ data.secretsmanager_software_license.my_license.activation_date }%{ endif }
Expiration Date:  %{ if data.secretsmanager_software_license.my_license.expiration_date != null ~}${ data.secretsmanager_software_license.my_license.expiration_date }%{ endif }

FileRefs:
---------
%{ for fr in data.secretsmanager_software_license.my_license.file_ref ~}
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
  value = data.secretsmanager_software_license.my_license.license_number
}
output "expiration_date" {
  value = data.secretsmanager_software_license.my_license.expiration_date
}
