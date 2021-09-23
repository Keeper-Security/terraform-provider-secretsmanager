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

data "keeper_secret_field" "my_field" {
  path        = "<record UID>/field/login"
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
Path:   ${ data.keeper_secret_field.my_field.path }
Value:  ${ data.keeper_secret_field.my_field.value }
EOT
}

output "field_value" {
  value = data.keeper_secret_field.my_field.value
}
