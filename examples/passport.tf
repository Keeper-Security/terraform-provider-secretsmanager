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

data "keeper_secret_passport" "my_passport" {
  path        = "<record UID>"
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
UID:    ${ data.keeper_secret_passport.my_passport.path }
Type:   ${ data.keeper_secret_passport.my_passport.type }
Title:  ${ data.keeper_secret_passport.my_passport.title }
Notes:  ${ data.keeper_secret_passport.my_passport.notes }
======

Passport#: ${ data.keeper_secret_passport.my_passport.passport_number }

Name:
-----
%{ for n in data.keeper_secret_passport.my_passport.name ~}
First Name:   ${ n.first }
Midlle Name:  ${ n.middle }
Last Name:    ${ n.last }

%{ endfor ~}

Birth Date:   %{ if data.keeper_secret_passport.my_passport.birth_date != null ~}${ data.keeper_secret_passport.my_passport.birth_date }%{ endif ~}

Exp.  Date:   %{ if data.keeper_secret_passport.my_passport.expiration_date != null ~}${ data.keeper_secret_passport.my_passport.expiration_date }%{ endif ~}

Date Issued:  %{ if data.keeper_secret_passport.my_passport.date_issued != null ~}${ data.keeper_secret_passport.my_passport.date_issued }%{ endif ~}


Password:     ${ data.keeper_secret_passport.my_passport.password }

AddressRefs:
------------
%{ if data.keeper_secret_passport.my_passport.address_ref != null }
%{ for a in data.keeper_secret_passport.my_passport.address_ref ~}
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
%{ for fr in data.keeper_secret_passport.my_passport.file_ref ~}
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

output "passport_number" {
  value = data.keeper_secret_passport.my_passport.passport_number
}
