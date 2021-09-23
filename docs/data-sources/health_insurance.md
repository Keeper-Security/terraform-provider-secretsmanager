# keeper_secret_health_insurance Data Source

Use this data source to read secrets of type `healthInsurance` stored in Keeper Vault

## Example Usage

```terraform
data "keeper_secret_health_insurance" "health_insurance" {
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
* `account_number` - Account Number.
* `name` - A list containing name information:
  - `first` - First name
  - `middle` - Middle name
  - `last` - Last name
* `login` - Account login.
* `password` - Account password.
* `url` - Account URL.
* `file_ref` - A list containing file reference information:
  - `uid` - File UID
  - `title` - File title
  - `name` - File name
  - `type` - File content type
  - `size` - File size
  - `last_modified` - File last modification timestamp
  - `url` - File download URL
  - `content_base64` - File content base64 encoded