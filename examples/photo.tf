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

data "keeper_secret_photo" "my_photos" {
  path        = "<record UID>"
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
UID:    ${ data.keeper_secret_photo.my_photos.path }
Type:   ${ data.keeper_secret_photo.my_photos.type }
Title:  ${ data.keeper_secret_photo.my_photos.title }
Notes:  ${ data.keeper_secret_photo.my_photos.notes }
======

FileRefs:
---------
%{ for fr in data.keeper_secret_photo.my_photos.file_ref ~}
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
  value = length(data.keeper_secret_photo.my_photos.file_ref.*)
  sensitive = true
}
