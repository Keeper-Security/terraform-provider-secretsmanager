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

data "secretsmanager_passport" "my_passport" {
  path        = "<record UID>"
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
UID:    ${ data.secretsmanager_passport.my_passport.path }
Type:   ${ data.secretsmanager_passport.my_passport.type }
Title:  ${ data.secretsmanager_passport.my_passport.title }
Notes:  ${ data.secretsmanager_passport.my_passport.notes }
======

Passport#: ${ data.secretsmanager_passport.my_passport.passport_number }

Name:
-----
%{ for n in data.secretsmanager_passport.my_passport.name ~}
First Name:   ${ n.first }
Midlle Name:  ${ n.middle }
Last Name:    ${ n.last }

%{ endfor ~}

Birth Date:   %{ if data.secretsmanager_passport.my_passport.birth_date != null ~}${ data.secretsmanager_passport.my_passport.birth_date }%{ endif ~}

Exp.  Date:   %{ if data.secretsmanager_passport.my_passport.expiration_date != null ~}${ data.secretsmanager_passport.my_passport.expiration_date }%{ endif ~}

Date Issued:  %{ if data.secretsmanager_passport.my_passport.date_issued != null ~}${ data.secretsmanager_passport.my_passport.date_issued }%{ endif ~}


Password:     ${ data.secretsmanager_passport.my_passport.password }

AddressRefs:
------------
%{ if data.secretsmanager_passport.my_passport.address_ref != null }
%{ for a in data.secretsmanager_passport.my_passport.address_ref ~}
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
%{ for fr in data.secretsmanager_passport.my_passport.file_ref ~}
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

output "passport_number" {
  value = data.secretsmanager_passport.my_passport.passport_number
}
