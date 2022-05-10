# secretsmanager_bank_account Data Source

Use this data source to read secrets of type `bankAccount` stored in Keeper Vault

## Example Usage

```terraform
data "secretsmanager_bank_account" "bank_account" {
  path = "<record UID>"
}
```

## Argument Reference

* `path` - (Required) The UID of existing record in Keeper Vault.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `type` - The type of the record.
* `title` - Record title.
* `notes` - Record notes.
* `bank_account` - A list containing bank account information:
  - `account_type` - Account type
  - `other_type` - Other type (if specified)
  - `routing_number` - City
  - `account_number` - State
* `name` - A list containing name information:
  - `first` - First name
  - `middle` - Middle name
  - `last` - Last name
* `login` - Account login.
* `password` - Account password.
* `url` - Account URL.
* `card_ref` - A list containing card reference information.
  - `uid` - Card reference UID
  - `payment_card` - A list containing payment card information:
    - `card_number` - Card number
    - `card_expiration_date` - Card expiration date
    - `card_security_code` - Card security code
  - `cardholder_name` - Cardholder name
  - `pin_code` - PIN code
* `file_ref` - A list containing file reference information:
  - `uid` - File UID
  - `title` - File title
  - `name` - File name
  - `type` - File content type
  - `size` - File size
  - `last_modified` - File last modification timestamp
  - `content_base64` - File content base64 encoded
* `totp` - A list containing Time-based One-time password information:
  - `url` - TOTP URL
  - `token` - Current TOTP password
  - `ttl` - Time to live in seconds for current token
