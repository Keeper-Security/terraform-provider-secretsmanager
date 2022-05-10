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

data "secretsmanager_membership" "my_membership" {
  path        = "<record UID>"
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
UID:    ${ data.secretsmanager_membership.my_membership.path }
Type:   ${ data.secretsmanager_membership.my_membership.type }
Title:  ${ data.secretsmanager_membership.my_membership.title }
Notes:  ${ data.secretsmanager_membership.my_membership.notes }
======

Acct.#:   ${ data.secretsmanager_membership.my_membership.account_number }
Password: ${ data.secretsmanager_membership.my_membership.password }

Name:
-----
%{ for n in data.secretsmanager_membership.my_membership.name ~}
First Name:   ${ n.first }
Midlle Name:  ${ n.middle }
Last Name:    ${ n.last }

%{ endfor ~}

FileRefs:
---------
%{ for fr in data.secretsmanager_membership.my_membership.file_ref ~}
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

output "account_number" {
  value = data.secretsmanager_membership.my_membership.account_number
}
