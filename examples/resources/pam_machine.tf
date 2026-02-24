terraform {
  required_providers {
    secretsmanager = {
      source  = "keeper-security/secretsmanager"
      version = ">= 1.2.0"
    }
    local = {
      source  = "hashicorp/local"
      version = "2.1.0"
    }
  }
}

provider "local" {}
provider "secretsmanager" {
  credential = "<CREDENTIAL>"
  # credential = file("~/.keeper/credential")
}

# Example 1: PAM Machine with SSH protocol
resource "secretsmanager_pam_machine" "ssh_server" {
  folder_uid = "<folder UID>"
  title      = "Production SSH Server"
  notes      = "Main production SSH gateway"

  pam_hostname {
    value {
      hostname = "ssh.prod.example.com"
      port     = "22"
    }
  }

  # SSH-specific connection settings
  pam_settings = jsonencode([{
    connection = [{
      protocol             = "ssh"
      port                 = "22"
      recordingIncludeKeys = true
      colorScheme          = "green_black"
      allowSupplyUser      = false
      fontSize             = "14"
      command              = "/bin/bash"
    }]
    portForward = [{
      port      = "2222"
      reusePort = true
    }]
  }])
}

# Example 2: PAM Machine with RDP protocol
resource "secretsmanager_pam_machine" "windows_server" {
  folder_uid = "<folder UID>"
  title      = "Windows RDP Server"
  notes      = "Windows Server 2022 for development"

  pam_hostname {
    value {
      hostname = "win-dev.example.com"
      port     = "3389"
    }
  }

  # RDP-specific connection settings
  pam_settings = jsonencode([{
    connection = [{
      protocol             = "rdp"
      port                 = "3389"
      recordingIncludeKeys = false
      security             = "nla"
      ignoreCert           = true
      resizeMethod         = "display-update"
      enableFullWindowDrag = true
      enableWallpaper      = false
    }]
  }])

  # Optional: Operating system
  # operating_system {
  #   label = "OS"
  #   value = "Windows Server 2022"
  # }

}

# Example 3: PAM Machine with cloud instance metadata
resource "secretsmanager_pam_machine" "aws_instance" {
  folder_uid = "<folder UID>"
  title      = "AWS EC2 Instance"
  notes      = "Production EC2 web server"

  pam_hostname {
    value {
      hostname = "ec2-10-0-1-100.compute-1.amazonaws.com"
      port     = "22"
    }
  }

  pam_settings = jsonencode([{
    connection = [{
      protocol             = "ssh"
      port                 = "22"
      recordingIncludeKeys = true
    }]
  }])

  # Instance metadata
  instance_name {
    label = "Instance Name"
    value = "web-server-prod-01"
  }

  instance_id {
    label = "Instance ID"
    value = "i-0123456789abcdef0"
  }

  provider_group {
    label = "Provider"
    value = "AWS"
  }

  provider_region {
    label = "Region"
    value = "us-east-1"
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
