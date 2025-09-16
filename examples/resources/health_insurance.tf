terraform {
  required_providers {
    secretsmanager = {
      source  = "keeper-security/secretsmanager"
      version = ">= 1.1.7"
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

resource "secretsmanager_health_insurance" "my_health_insurance" {
  folder_uid = "<folder UID>"
  title = "My Title"
  notes = "My Notes"
  account_number {
    label = "My Account"
    required = true
    privacy_screen = true
    value = "My Account# 1234"
  }
  name {
    label = "John"
    required = true
    privacy_screen = true
    value {
      first = "John"
      middle = "D"
      last = "Doe"
    }
  }
  login {
    label = "My Login"
    required = true
    privacy_screen = true
    value = "MyLogin"
  }
  password {
    label = "My Pass"
    required = true
    privacy_screen = true
    enforce_generation = true
    generate = "yes"
    complexity {
      length = 20
      caps = 5
      lowercase = 5
      digits = 5
      special = 5
    }
    #value = "to_be_generated"
  }
  url {
    label = "My Url"
    required = true
    privacy_screen = true
    value = "https://192.168.1.1/"
  }
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
FUID:   ${ secretsmanager_health_insurance.my_health_insurance.folder_uid }
UID:    ${ secretsmanager_health_insurance.my_health_insurance.uid }
Type:   ${ secretsmanager_health_insurance.my_health_insurance.type }
Title:  ${ secretsmanager_health_insurance.my_health_insurance.title }
Notes:  ${ secretsmanager_health_insurance.my_health_insurance.notes }
======

Account Number:
---------------
%{ for n in secretsmanager_health_insurance.my_health_insurance.account_number ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Value:    ${ n.value }
%{ endfor }

Name:
-----
%{ for n in secretsmanager_health_insurance.my_health_insurance.name ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
First Name:     ${ n.value.0.first }
Middle Name:    ${ n.value.0.middle }
Last Name:      ${ n.value.0.last }
%{ endfor }

Login:
------
%{ for n in secretsmanager_health_insurance.my_health_insurance.login ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Value:    ${ n.value }
%{ endfor }

Password:
---------
%{ for n in secretsmanager_health_insurance.my_health_insurance.password ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Enforce Generation: ${ n.enforce_generation }
Generate: %{ if n.generate != null }${n.generate}%{ endif }
Complexity: Length = ${ n.complexity.0.length }
Value:    ${ n.value }
%{ endfor }

URL:
----
%{ for n in secretsmanager_health_insurance.my_health_insurance.url ~}
Type:     ${ n.type }
Label:    ${ n.label }
Required: ${ n.required }
Privacy Screen: ${ n.privacy_screen }
Value:    ${ n.value }
%{ endfor }

EOT
}

output "record_uid" {
  value = secretsmanager_health_insurance.my_health_insurance.uid
}
output "record_title" {
  value = secretsmanager_health_insurance.my_health_insurance.title
}
