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

data "keeper_secret_address" "address" {
  path        = "<record UID>"
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
UID:    ${ data.keeper_secret_address.address.path }
Type:   ${ data.keeper_secret_address.address.type }
Title:  ${ data.keeper_secret_address.address.title }
Notes:  ${ data.keeper_secret_address.address.notes }
======

Address:
--------
%{ for a in data.keeper_secret_address.address.address ~}
Street1:  ${ a.street1 }
Street2:  ${ a.street2 }
City:     ${ a.city }
State:    ${ a.state }
Zip:      ${ a.zip }
Country:  ${ a.country }
%{ endfor ~}

FileRefs:
---------
%{ for fr in data.keeper_secret_address.address.file_ref ~}
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

output "files" {
  value = data.keeper_secret_address.address.file_ref
  sensitive = true
}

output "first_file" {
  value = length(data.keeper_secret_address.address.file_ref) < 1 ? "" : textdecodebase64(element(data.keeper_secret_address.address.file_ref, 0).content_base64, "UTF-8")
  sensitive = true
}

output "zip_code" {
  value = length(data.keeper_secret_address.address.address) < 1 ? "" : data.keeper_secret_address.address.address.0.zip
}
