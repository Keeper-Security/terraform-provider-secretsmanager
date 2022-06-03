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

resource "secretsmanager_bank_card" "my_bank_card" {
  folder_uid = "<folder UID>"
  title = "My Title"
  notes = "My Notes"
  payment_card {
    label = "My Card"
    required = true
    privacy_screen = true
    value {
      card_number = "123456780"
      card_expiration_date = "12/2121"
      card_security_code = "787"
    }
  }
  cardholder_name {
    label = "My Card Name"
    required = true
    privacy_screen = true
    value = "John Doe"
  }
  pin_code {
    label = "My Pin Code"
    required = true
    privacy_screen = true
    value = "7870"
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
FUID:   ${ secretsmanager_bank_card.my_bank_card.folder_uid }
UID:    ${ secretsmanager_bank_card.my_bank_card.uid }
Type:   ${ secretsmanager_bank_card.my_bank_card.type }
Title:  ${ secretsmanager_bank_card.my_bank_card.title }
Notes:  ${ secretsmanager_bank_card.my_bank_card.notes }
======

Payment Card:
-------------
%{ for n in secretsmanager_bank_card.my_bank_card.payment_card ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Card Number:    ${ n.value.0.card_number }
Card Expiration Date: ${ n.value.0.card_expiration_date }
Card Security Code:   ${ n.value.0.card_security_code }
%{ endfor }

Cardholder Name:
----------------
%{ for n in secretsmanager_bank_card.my_bank_card.cardholder_name ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Value:    ${ n.value }
%{ endfor }

Pin Code:
---------
%{ for n in secretsmanager_bank_card.my_bank_card.pin_code ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Value:    ${ n.value }
%{ endfor }

Address Ref:
------------
%{ for n in secretsmanager_bank_card.my_bank_card.address_ref ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Value:    ${ n.value }
%{ endfor }

EOT
}

output "record_uid" {
  value = secretsmanager_bank_card.my_bank_card.uid
}
output "record_title" {
  value = secretsmanager_bank_card.my_bank_card.title
}
