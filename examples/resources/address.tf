terraform {
  required_providers {
    secretsmanager = {
      source  = "keeper-security/secretsmanager"
      version = ">= 1.1.7"
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

resource "secretsmanager_address" "my_address" {
  folder_uid = "<folder UID>"
  title = "My Title"
  notes = "My Notes"
  address {
    label = "My Address"
    required = true
    privacy_screen = true
    value {
      street1 = "7422 Avalon"
      street2 = "Apt 21"
      city = "Los Angeles"
      state = "CA"
      country = "United States"
      zip = "90003-2334"
    }
  }
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
FUID:   ${ secretsmanager_address.my_address.folder_uid }
UID:    ${ secretsmanager_address.my_address.uid }
Type:   ${ secretsmanager_address.my_address.type }
Title:  ${ secretsmanager_address.my_address.title }
Notes:  ${ secretsmanager_address.my_address.notes }
======

Address:
--------
%{ for n in secretsmanager_address.my_address.address ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen:   ${ n.privacy_screen }

Street1:  ${ n.value.0.street1 }
Street2:  ${ n.value.0.street2 }
City:     ${ n.value.0.city }
State:    ${ n.value.0.state }
Country:  ${ n.value.0.country }
ZIP:      ${ n.value.0.zip }
%{ endfor }

EOT
}

output "record_uid" {
  value = secretsmanager_address.my_address.uid
}
output "record_title" {
  value = secretsmanager_address.my_address.title
}
