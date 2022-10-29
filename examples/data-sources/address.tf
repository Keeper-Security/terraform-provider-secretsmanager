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

data "secretsmanager_address" "my_address" {
  path        = "<record UID>"
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
UID:    ${ data.secretsmanager_address.my_address.path }
Type:   ${ data.secretsmanager_address.my_address.type }
Title:  ${ data.secretsmanager_address.my_address.title }
Notes:  ${ data.secretsmanager_address.my_address.notes }
======

Address:
--------
%{ for a in data.secretsmanager_address.my_address.address ~}
Street1:  ${ a.street1 }
Street2:  ${ a.street2 }
City:     ${ a.city }
State:    ${ a.state }
Zip:      ${ a.zip }
Country:  ${ a.country }
%{ endfor ~}

FileRefs:
---------
%{ for fr in data.secretsmanager_address.my_address.file_ref ~}
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

output "files" {
  value = data.secretsmanager_address.my_address.file_ref
  sensitive = true
}

output "first_file" {
  value = length(data.secretsmanager_address.my_address.file_ref) < 1 ? "" : textdecodebase64(element(data.secretsmanager_address.my_address.file_ref, 0).content_base64, "UTF-8")
  sensitive = true
}

output "zip_code" {
  value = length(data.secretsmanager_address.my_address.address) < 1 ? "" : data.secretsmanager_address.my_address.address.0.zip
}
