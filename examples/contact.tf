terraform {
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

data "keeper_secret_contact" "my_contact" {
  path        = "<record UID>"
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
UID:    ${ data.keeper_secret_contact.my_contact.path }
Type:   ${ data.keeper_secret_contact.my_contact.type }
Title:  ${ data.keeper_secret_contact.my_contact.title }
Notes:  ${ data.keeper_secret_contact.my_contact.notes }
======

Name:
-----
%{ for n in data.keeper_secret_contact.my_contact.name ~}
First Name:   ${ n.first }
Midlle Name:  ${ n.middle }
Last Name:    ${ n.last }

%{ endfor ~}

Company:  ${ data.keeper_secret_contact.my_contact.company }
Email:    ${ data.keeper_secret_contact.my_contact.email }

Phone:
------
%{ if data.keeper_secret_contact.my_contact.phone != null ~}
%{ for p in data.keeper_secret_contact.my_contact.phone ~}
Region: ${ p.region }
Number: ${ p.number }
Ext.:   ${ p.ext }
Type:   ${ p.type }

%{ endfor ~}
%{ endif ~}

AddressRefs:
------------
%{ if data.keeper_secret_contact.my_contact.address_ref != null }
%{ for a in data.keeper_secret_contact.my_contact.address_ref ~}
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
%{ for fr in data.keeper_secret_contact.my_contact.file_ref ~}
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

output "name" {
  value = data.keeper_secret_contact.my_contact.name
}
