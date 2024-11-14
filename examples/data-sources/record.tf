terraform {
  required_providers {
    secretsmanager = {
      source  = "keeper-security/secretsmanager"
      version = ">= 1.1.5"
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

data "secretsmanager_record" "my_record" {
  path        = "<record UID>"
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
UID:    ${ data.secretsmanager_record.my_record.path }
Type:   ${ data.secretsmanager_record.my_record.type }
Title:  ${ data.secretsmanager_record.my_record.title }
Notes:  ${ data.secretsmanager_record.my_record.notes }
======

Fields:
--------
%{ for f in data.secretsmanager_record.my_record.fields ~}
Type:               ${ f.type }
Label:              ${ f.label }
Required:           ${ f.required }
Privacy Screen:     ${ f.privacy_screen }
Enforce Generation: ${ f.enforce_generation }
Value:              ${ f.value }
${ length(f.complexity) < 1 ? "" : format("Complexity: %#v", f.complexity.0) }
%{ endfor ~}

Custom Fields:
--------
%{ for c in data.secretsmanager_record.my_record.custom ~}
Type:               ${ c.type }
Label:              ${ c.label }
Required:           ${ c.required }
Privacy Screen:     ${ c.privacy_screen }
Enforce Generation: ${ c.enforce_generation }
Value:              ${ c.value }
${ length(c.complexity) < 1 ? "" : format("Complexity: %#v", c.complexity.0) }
%{ endfor ~}

FileRefs:
---------
%{ for fr in data.secretsmanager_record.my_record.file_ref ~}
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

output "files" {
  value = data.secretsmanager_record.my_record.file_ref
  sensitive = true
}

output "first_field_type" {
  value = length(data.secretsmanager_record.my_record.fields) < 1 ? "" : data.secretsmanager_record.my_record.fields.0.type
}
output "first_field_value" {
  value = length(data.secretsmanager_record.my_record.fields) < 1 ? "" : data.secretsmanager_record.my_record.fields.0.value
}
