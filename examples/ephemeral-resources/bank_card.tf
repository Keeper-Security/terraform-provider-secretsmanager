terraform {
  required_providers {
    secretsmanager = {
      source  = "keeper-security/secretsmanager"
      version = ">= 1.2.0"
    }
  }
}

provider "secretsmanager" {
  credential = "<CREDENTIAL>"
  # credential = file("~/.keeper/credential")
}

# Ephemeral resources do not store secret values in the Terraform state file.
# This makes them a more secure option for accessing sensitive credentials.

ephemeral "secretsmanager_bank_card" "my_card" {
  path = "<record UID>"
}

output "cardholder_name" {
  value     = ephemeral.secretsmanager_bank_card.my_card.cardholder_name
  ephemeral = true
}

output "pin_code" {
  value     = ephemeral.secretsmanager_bank_card.my_card.pin_code
  ephemeral = true
}

output "card_number" {
  value     = length(ephemeral.secretsmanager_bank_card.my_card.payment_card) < 1 ? "" : ephemeral.secretsmanager_bank_card.my_card.payment_card.0.card_number
  ephemeral = true
}
