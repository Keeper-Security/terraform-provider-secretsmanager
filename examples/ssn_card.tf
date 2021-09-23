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

data "keeper_secret_ssn_card" "my_ssn" {
  path        = "<record UID>"
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
UID:    ${ data.keeper_secret_ssn_card.my_ssn.path }
Type:   ${ data.keeper_secret_ssn_card.my_ssn.type }
Title:  ${ data.keeper_secret_ssn_card.my_ssn.title }
Notes:  ${ data.keeper_secret_ssn_card.my_ssn.notes }
======

Identity Number:  ${ data.keeper_secret_ssn_card.my_ssn.identity_number }

Name:
-----
%{ for n in data.keeper_secret_ssn_card.my_ssn.name ~}
First Name:   ${ n.first }
Midlle Name:  ${ n.middle }
Last Name:    ${ n.last }
%{ endfor }

FileRefs:
---------
%{ for fr in data.keeper_secret_ssn_card.my_ssn.file_ref ~}
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

output "identity_number" {
  value = data.keeper_secret_ssn_card.my_ssn.identity_number
}
