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

data "secretsmanager_bank_account" "my_account" {
  path        = "<record UID>"
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
UID:    ${ data.secretsmanager_bank_account.my_account.path }
Type:   ${ data.secretsmanager_bank_account.my_account.type }
Title:  ${ data.secretsmanager_bank_account.my_account.title }
Notes:  ${ data.secretsmanager_bank_account.my_account.notes }
======

Bank Account:
-------------
%{ for a in data.secretsmanager_bank_account.my_account.bank_account ~}
Account Type:   ${ a.account_type }
Other Type:     ${ a.other_type }
Routing Number: ${ a.routing_number }
Account Number: ${ a.account_number }

%{ endfor ~}

Name:
-----
%{ for n in data.secretsmanager_bank_account.my_account.name ~}
First Name:   ${ n.first }
Midlle Name:  ${ n.middle }
Last Name:    ${ n.last }

%{ endfor ~}

Login:    ${ data.secretsmanager_bank_account.my_account.login }
Password: ${ data.secretsmanager_bank_account.my_account.password }
URL:      ${ data.secretsmanager_bank_account.my_account.url }

Card Ref:
---------
%{ for r in data.secretsmanager_bank_account.my_account.card_ref ~}
 Card Reference#:  ${ r.uid }
 ----------------
%{ for cc in r.payment_card ~}
  Card Number:          ${ cc.card_number }
  Card Expiration Date: ${ cc.card_expiration_date }
  Card Security Code:   ${ cc.card_security_code }
%{ endfor ~}
 Cardholder Name:  ${ r.cardholder_name }
 PIN Code:         ${ r.pin_code }
%{ endfor }

FileRefs:
---------
%{ for fr in data.secretsmanager_bank_account.my_account.file_ref ~}
UID:      ${ fr.uid }
Title:    ${ fr.title }
Name:     ${ fr.name }
Type:     ${ fr.type }
Size:     ${ fr.size }
Last Modified:  ${ fr.last_modified }
URL:            ${ fr.url }

Content/Base64: ${ fr.content_base64 }


%{ endfor ~}

TOTP:
-----
%{ for t in data.secretsmanager_bank_account.my_account.totp ~}
URL:    ${ t.url }
Token:  ${ t.token }
TTL:    ${ t.ttl }

%{ endfor ~}

EOT
}

output "name" {
  value = data.secretsmanager_bank_account.my_account.name
}
output "totp" {
  value = data.secretsmanager_bank_account.my_account.totp
}
output "totp_token" {
  value = length(data.secretsmanager_bank_account.my_account.totp) < 1 ? "" : data.secretsmanager_bank_account.my_account.totp.0.token
}
