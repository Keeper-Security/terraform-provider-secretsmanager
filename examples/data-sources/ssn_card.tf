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

data "secretsmanager_ssn_card" "my_ssn" {
  path        = "<record UID>"
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
UID:    ${ data.secretsmanager_ssn_card.my_ssn.path }
Type:   ${ data.secretsmanager_ssn_card.my_ssn.type }
Title:  ${ data.secretsmanager_ssn_card.my_ssn.title }
Notes:  ${ data.secretsmanager_ssn_card.my_ssn.notes }
======

Identity Number:  ${ data.secretsmanager_ssn_card.my_ssn.identity_number }

Name:
-----
%{ for n in data.secretsmanager_ssn_card.my_ssn.name ~}
First Name:   ${ n.first }
Midlle Name:  ${ n.middle }
Last Name:    ${ n.last }
%{ endfor }

FileRefs:
---------
%{ for fr in data.secretsmanager_ssn_card.my_ssn.file_ref ~}
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

output "identity_number" {
  value = data.secretsmanager_ssn_card.my_ssn.identity_number
}
