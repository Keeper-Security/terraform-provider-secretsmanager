terraform {
  required_version = ">= 1.0.0"
  required_providers {
    keeper = {
      source  = "github.com/keeper-security/keeper"
      version = ">= 0.1.0"
    }
    local = {
      source = "hashicorp/local"
      version = "2.1.0"
    }
  }
}

provider "local" { }
provider "keeper" {
  credential = "<CREDENTIAL>"
  # credential = file("~/.keeper/credential")
}

data "keeper_secret_software_license" "my_license" {
  path        = "<record UID>"
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
UID:    ${ data.keeper_secret_software_license.my_license.path }
Type:   ${ data.keeper_secret_software_license.my_license.type }
Title:  ${ data.keeper_secret_software_license.my_license.title }
Notes:  ${ data.keeper_secret_software_license.my_license.notes }
======

License#:  ${ data.keeper_secret_software_license.my_license.license_number }
Activation Date:  %{ if data.keeper_secret_software_license.my_license.activation_date != null ~}${ data.keeper_secret_software_license.my_license.activation_date }%{ endif }
Expiration Date:  %{ if data.keeper_secret_software_license.my_license.expiration_date != null ~}${ data.keeper_secret_software_license.my_license.expiration_date }%{ endif }

FileRefs:
---------
%{ for fr in data.keeper_secret_software_license.my_license.file_ref ~}
UID:      ${ fr.uid }
Title:    ${ fr.title }
Name:     ${ fr.name }
Type:     ${ fr.type }
Size:     ${ fr.size }
Last Modified:  ${ fr.last_modified }
URL:            ${ fr.url }

Content/Base64: ${ fr.content_base64 }


%{ endfor ~}
EOT
}

output "driver_license_number" {
  value = data.keeper_secret_software_license.my_license.license_number
}
output "expiration_date" {
  value = data.keeper_secret_software_license.my_license.expiration_date
}
