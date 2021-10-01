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

data "keeper_secret_birth_certificate" "my_birth_cert" {
  path        = "<record UID>"
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
UID:    ${ data.keeper_secret_birth_certificate.my_birth_cert.path }
Type:   ${ data.keeper_secret_birth_certificate.my_birth_cert.type }
Title:  ${ data.keeper_secret_birth_certificate.my_birth_cert.title }
Notes:  ${ data.keeper_secret_birth_certificate.my_birth_cert.notes }
======

Name:
-----
%{ for n in data.keeper_secret_birth_certificate.my_birth_cert.name ~}
First Name:   ${ n.first }
Midlle Name:  ${ n.middle }
Last Name:    ${ n.last }

%{ endfor ~}

Birth Date:    %{ if data.keeper_secret_birth_certificate.my_birth_cert.birth_date != null ~}
${ data.keeper_secret_birth_certificate.my_birth_cert.birth_date }
%{ endif ~}

FileRefs:
---------
%{ for fr in data.keeper_secret_birth_certificate.my_birth_cert.file_ref ~}
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

output "birth_date" {
  value = formatdate("YYYY-MM-DD", data.keeper_secret_birth_certificate.my_birth_cert.birth_date)
}
