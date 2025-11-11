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

# Example 1: PAM Machine with SSH protocol
resource "secretsmanager_pam_machine" "ssh_server" {
  folder_uid = "<folder UID>"
  title = "Production SSH Server"
  notes = "Main production SSH gateway"

  pam_hostname {
    value {
      hostname = "ssh.prod.example.com"
      port = "22"
    }
  }

  # SSH-specific connection settings
  pam_settings = jsonencode([{
    connection = [{
      protocol = "ssh"
      port = "22"
      recordingIncludeKeys = true
      colorScheme = "green_black"
      allowSupplyUser = false
      fontSize = "14"
      command = "/bin/bash"
    }]
    portForward = [{
      port = "2222"
      reusePort = true
    }]
  }])

  login {
    label = "Admin Username"
    required = true
    value = "admin"
  }

  password {
    label = "Admin Password"
    required = true
    privacy_screen = true
    enforce_generation = true
    generate = "yes"
    complexity {
      length = 32
      caps = 8
      lowercase = 8
      digits = 8
      special = 8
    }
  }
}

# Example 2: PAM Machine with RDP protocol
resource "secretsmanager_pam_machine" "windows_server" {
  folder_uid = "<folder UID>"
  title = "Windows RDP Server"
  notes = "Windows Server 2022 for development"

  pam_hostname {
    value {
      hostname = "win-dev.example.com"
      port = "3389"
    }
  }

  # RDP-specific connection settings
  pam_settings = jsonencode([{
    connection = [{
      protocol = "rdp"
      port = "3389"
      recordingIncludeKeys = false
      security = "nla"
      ignoreCert = true
      resizeMethod = "display-update"
      enableFullWindowDrag = true
      enableWallpaper = false
    }]
  }])

  login {
    value = "Administrator"
  }

  password {
    value = "Str0ng!P@ssw0rd"
    privacy_screen = true
  }

  # Optional: Operating system
  # operating_system {
  #   label = "OS"
  #   value = ["Windows Server 2022"]
  # }

  # Optional: SSL verification
  # ssl_verification {
  #   value = [true]
  # }
}

# Example 3: PAM Machine with cloud instance metadata
resource "secretsmanager_pam_machine" "aws_instance" {
  folder_uid = "<folder UID>"
  title = "AWS EC2 Instance"
  notes = "Production EC2 web server"

  pam_hostname {
    value {
      hostname = "ec2-10-0-1-100.compute-1.amazonaws.com"
      port = "22"
    }
  }

  pam_settings = jsonencode([{
    connection = [{
      protocol = "ssh"
      port = "22"
      recordingIncludeKeys = true
    }]
  }])

  login {
    value = "ec2-user"
  }

  # Instance metadata
  instance_name {
    label = "Instance Name"
    value = ["web-server-prod-01"]
  }

  instance_id {
    label = "Instance ID"
    value = ["i-0123456789abcdef0"]
  }

  provider_group {
    label = "Provider"
    value = ["AWS"]
  }

  provider_region {
    label = "Region"
    value = ["us-east-1"]
  }
}

# Output the created machine records
output "ssh_server_uid" {
  value = secretsmanager_pam_machine.ssh_server.uid
}

output "ssh_server_hostname" {
  value = secretsmanager_pam_machine.ssh_server.pam_hostname[0].hostname
}

output "windows_server_uid" {
  value = secretsmanager_pam_machine.windows_server.uid
}
