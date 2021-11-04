# secretsmanager_bank_card Data Source

Use this data source to read secrets of type `bankCard` stored in Keeper Vault

## Example Usage

```terraform
data "secretsmanager_bank_card" "bank_card" {
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
* `payment_card` - A list containing payment card information:
  - `card_number` - Card number
  - `card_expiration_date` - Card expiration date
  - `card_security_code` - Card security code
* `cardholder_name` - Cardholder name.
* `pin_code` - PIN code.
* `address_ref` - A list containing address information:
  - `uid` - The address reference record UID
  - `street1` - Street line 1
  - `street2` - Street line 2
  - `city` - City
  - `state` - State
  - `zip` - Zip
  - `country` - Country
* `file_ref` - A list containing file reference information:
  - `uid` - File UID
  - `title` - File title
  - `name` - File name
  - `type` - File content type
  - `size` - File size
  - `last_modified` - File last modification timestamp
  - `url` - File download URL
  - `content_base64` - File content base64 encoded
