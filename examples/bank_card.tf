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

data "secretsmanager_bank_card" "my_card" {
  path        = "<record UID>"
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
UID:    ${ data.secretsmanager_bank_card.my_card.path }
Type:   ${ data.secretsmanager_bank_card.my_card.type }
Title:  ${ data.secretsmanager_bank_card.my_card.title }
Notes:  ${ data.secretsmanager_bank_card.my_card.notes }
======

Payment Card:
-------------
%{ for cc in data.secretsmanager_bank_card.my_card.payment_card ~}
Card Number:          ${ cc.card_number }
Card Expiration Date: ${ cc.card_expiration_date }
Card Security Code:   ${ cc.card_security_code }

%{ endfor ~}

Cardholder Name:  ${ data.secretsmanager_bank_card.my_card.cardholder_name }
PIN Code:         ${ data.secretsmanager_bank_card.my_card.pin_code }

AddressRefs:
--------
%{ if data.secretsmanager_bank_card.my_card.address_ref != null }
%{ for a in data.secretsmanager_bank_card.my_card.address_ref ~}
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
%{ for fr in data.secretsmanager_bank_card.my_card.file_ref ~}
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

output "cardholder_name" {
  value = data.secretsmanager_bank_card.my_card.cardholder_name
}
