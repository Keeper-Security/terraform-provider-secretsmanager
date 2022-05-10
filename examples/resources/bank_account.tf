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

resource "secretsmanager_bank_account" "my_bank_account" {
  folder_uid = "<folder UID>"
  title = "My Title"
  notes = "My Notes"
  bank_account {
    label = "My Bank Account"
    required = true
    privacy_screen = true
    value {
      account_type = "Other"
      routing_number = "1234567890"
      account_number = "0987654321"
      other_type = "Investment"
    }
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
  login {
    label = "My Login"
    required = true
    privacy_screen = true
    value = "MyLogin"
  }
  password {
    label = "My Pass"
    required = true
    privacy_screen = true
    enforce_generation = true
    generate = "yes"
    complexity {
      length = 20
      caps = 5
      lowercase = 5
      digits = 5
      special = 5
    }
    #value = "to_be_generated"
  }
  url {
    label = "My Url"
    required = true
    privacy_screen = true
    value = "https://192.168.1.1/"
  }
  card_ref {
    label = "My Card Ref"
    required = true
    privacy_screen = true
    value = "<card ref UID>"
  }
  totp {
    label = "My TOTP"
    required = true
    privacy_screen = true
    value = "otpauth://totp/Acme:Buster?secret=6I4PI5EUKS66GPRY5TMLJJP25MAYWAVL&issuer=Acme&algorithm=SHA1&digits=6&period=30"
  }
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
FUID:   ${ secretsmanager_bank_account.my_bank_account.folder_uid }
UID:    ${ secretsmanager_bank_account.my_bank_account.uid }
Type:   ${ secretsmanager_bank_account.my_bank_account.type }
Title:  ${ secretsmanager_bank_account.my_bank_account.title }
Notes:  ${ secretsmanager_bank_account.my_bank_account.notes }
======

Bank Account:
-------------
%{ for n in secretsmanager_bank_account.my_bank_account.bank_account ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Account Type:   ${ n.value.0.account_type }
Other Type:     ${ n.value.0.other_type }
Routing Number Expiration Date: ${ n.value.0.routing_number }
Account Number Security Code:   ${ n.value.0.account_number }
%{ endfor }

Name:
-----
%{ for n in secretsmanager_bank_account.my_bank_account.name ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen:   ${ n.privacy_screen }
First Name:  ${ n.value.0.first }
Middle Name: ${ n.value.0.middle }
Last Name:   ${ n.value.0.last }
%{ endfor }

Login:
------
%{ for n in secretsmanager_bank_account.my_bank_account.login ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen:   ${ n.privacy_screen }
Value:   ${ n.value }
%{ endfor }

Password:
---------
%{ for n in secretsmanager_bank_account.my_bank_account.password ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen:     ${ n.privacy_screen }
Enforce Generation: ${ n.enforce_generation }
Generate: %{ if n.generate != null }${n.generate}%{ endif }
Complexity: Length = ${ n.complexity.0.length }
Value:    ${ n.value }
%{ endfor }

URL:
----
%{ for n in secretsmanager_bank_account.my_bank_account.url ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Value:    ${ n.value }
%{ endfor }

Card Ref:
---------
%{ for n in secretsmanager_bank_account.my_bank_account.card_ref ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Value:    ${ n.value }
%{ endfor }

TOTP:
-----
%{ for n in secretsmanager_bank_account.my_bank_account.totp ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Value:    ${ n.value }
%{ endfor }

EOT
}

output "record_uid" {
  value = secretsmanager_bank_account.my_bank_account.uid
}
output "record_title" {
  value = secretsmanager_bank_account.my_bank_account.title
}
