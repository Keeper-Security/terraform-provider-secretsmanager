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

data "keeper_secret_membership" "my_membership" {
  path        = "<record UID>"
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
UID:    ${ data.keeper_secret_membership.my_membership.path }
Type:   ${ data.keeper_secret_membership.my_membership.type }
Title:  ${ data.keeper_secret_membership.my_membership.title }
Notes:  ${ data.keeper_secret_membership.my_membership.notes }
======

Acct.#:   ${ data.keeper_secret_membership.my_membership.account_number }
Password: ${ data.keeper_secret_membership.my_membership.password }

Name:
-----
%{ for n in data.keeper_secret_membership.my_membership.name ~}
First Name:   ${ n.first }
Midlle Name:  ${ n.middle }
Last Name:    ${ n.last }

%{ endfor ~}

FileRefs:
---------
%{ for fr in data.keeper_secret_membership.my_membership.file_ref ~}
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

output "account_number" {
  value = data.keeper_secret_membership.my_membership.account_number
}
