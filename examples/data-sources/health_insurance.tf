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

data "secretsmanager_health_insurance" "my_insurance" {
  path        = "<record UID>"
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
UID:    ${ data.secretsmanager_health_insurance.my_insurance.path }
Type:   ${ data.secretsmanager_health_insurance.my_insurance.type }
Title:  ${ data.secretsmanager_health_insurance.my_insurance.title }
Notes:  ${ data.secretsmanager_health_insurance.my_insurance.notes }
======

Acct.#:   ${ data.secretsmanager_health_insurance.my_insurance.account_number }
Login:    ${ data.secretsmanager_health_insurance.my_insurance.login }
Password: ${ data.secretsmanager_health_insurance.my_insurance.password }
URL:      ${ data.secretsmanager_health_insurance.my_insurance.url }

Name:
-----
%{ for n in data.secretsmanager_health_insurance.my_insurance.name ~}
First Name:   ${ n.first }
Midlle Name:  ${ n.middle }
Last Name:    ${ n.last }

%{ endfor ~}

FileRefs:
---------
%{ for fr in data.secretsmanager_health_insurance.my_insurance.file_ref ~}
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

output "login" {
  value = data.secretsmanager_health_insurance.my_insurance.login
}
