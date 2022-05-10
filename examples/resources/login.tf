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

resource "secretsmanager_login" "my_login" {
  folder_uid = "<folder UID>"
  title = "My Title"
  notes = "My Notes"
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
  totp {
    label = "My TOTP"
    required = true
    privacy_screen = true
    value = "otpauth://totp/Acme:Buster?secret=6I4PI5EUKS66GPRY5TMLJJP25MAYWAVL&issuer=Acme&algorithm=SHA1&digits=6&period=30"
  }
}

resource "local_file" "out" {
  filename        = "${path.module}/out.txt"
  file_permission = "0644"
  content         = <<EOT
FUID:   ${ secretsmanager_login.my_login.folder_uid }
UID:    ${ secretsmanager_login.my_login.uid }
Type:   ${ secretsmanager_login.my_login.type }
Title:  ${ secretsmanager_login.my_login.title }
Notes:  ${ secretsmanager_login.my_login.notes }
======

Login:
------
%{ for n in secretsmanager_login.my_login.login ~}
Type:   ${ n.type }
Label:   ${ n.label }
Required:   ${ n.required }
Privacy Screen:   ${ n.privacy_screen }
Value:   ${ n.value }
%{ endfor }

Password:
---------
%{ for n in secretsmanager_login.my_login.password ~}
Type:   ${ n.type }
Label:  ${ n.label }
Required:   ${ n.required }
Privacy Screen:     ${ n.privacy_screen }
Enforce Generation: ${ n.enforce_generation }
Generate:   %{ if n.generate != null }${n.generate}%{ endif }
Complexity: Length = ${ n.complexity.0.length }
Value:   ${ n.value }
%{ endfor }

URL:
----
%{ for n in secretsmanager_login.my_login.url ~}
Type:   ${ n.type }
Label:   ${ n.label }
Required:   ${ n.required }
Privacy Screen:   ${ n.privacy_screen }
Value:   ${ n.value }
%{ endfor }

TOTP:
-----
%{ for n in secretsmanager_login.my_login.totp ~}
Type:   ${ n.type }
Label:   ${ n.label }
Required:   ${ n.required }
Privacy Screen:   ${ n.privacy_screen }
Value:   ${ n.value }
%{ endfor }

EOT
}

output "secret_uid" {
  value = secretsmanager_login.my_login.uid
  sensitive = true
}
output "secret_title" {
  value = secretsmanager_login.my_login.title
  sensitive = true
}
output "secret_login" {
  value = secretsmanager_login.my_login.login
}
