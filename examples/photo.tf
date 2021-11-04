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

data "secretsmanager_photo" "my_photos" {
  path        = "<record UID>"
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
UID:    ${ data.secretsmanager_photo.my_photos.path }
Type:   ${ data.secretsmanager_photo.my_photos.type }
Title:  ${ data.secretsmanager_photo.my_photos.title }
Notes:  ${ data.secretsmanager_photo.my_photos.notes }
======

FileRefs:
---------
%{ for fr in data.secretsmanager_photo.my_photos.file_ref ~}
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

output "photo_count" {
  value = length(data.secretsmanager_photo.my_photos.file_ref.*)
  sensitive = true
}
