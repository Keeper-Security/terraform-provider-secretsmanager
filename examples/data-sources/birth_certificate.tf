terraform {
  required_providers {
    secretsmanager = {
      source  = "keeper-security/secretsmanager"
      version = ">= 1.1.1"
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

data "secretsmanager_birth_certificate" "my_birth_cert" {
  path        = "<record UID>"
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
UID:    ${ data.secretsmanager_birth_certificate.my_birth_cert.path }
Type:   ${ data.secretsmanager_birth_certificate.my_birth_cert.type }
Title:  ${ data.secretsmanager_birth_certificate.my_birth_cert.title }
Notes:  ${ data.secretsmanager_birth_certificate.my_birth_cert.notes }
======

Name:
-----
%{ for n in data.secretsmanager_birth_certificate.my_birth_cert.name ~}
First Name:   ${ n.first }
Midlle Name:  ${ n.middle }
Last Name:    ${ n.last }

%{ endfor ~}

Birth Date:    %{ if data.secretsmanager_birth_certificate.my_birth_cert.birth_date != null ~}
${ data.secretsmanager_birth_certificate.my_birth_cert.birth_date }
%{ endif ~}

FileRefs:
---------
%{ for fr in data.secretsmanager_birth_certificate.my_birth_cert.file_ref ~}
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

output "birth_date" {
  value = formatdate("YYYY-MM-DD", data.secretsmanager_birth_certificate.my_birth_cert.birth_date)
}
