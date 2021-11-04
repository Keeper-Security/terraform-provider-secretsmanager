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

data "secretsmanager_field" "my_field" {
  path        = "<record UID>/field/login"
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
Path:   ${ data.secretsmanager_field.my_field.path }
Value:  ${ data.secretsmanager_field.my_field.value }
EOT
}

output "field_value" {
  value = data.secretsmanager_field.my_field.value
}
