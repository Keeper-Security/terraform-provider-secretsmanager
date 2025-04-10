# secretsmanager_record Data Source

Use this data source to read secrets of any type stored in Keeper Vault

## Example Usage

```terraform
data "secretsmanager_record" "record" {
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
* `fields` - A list containing fields information:
  - `type` - Field type
  - `label` - Field label
  - `required` - Required field flag
  - `privacy_screen` - Privacy screen flag
  - `enforce_generation` - Enforce generation flag (for password field)
  - `complexity` - A list containing password complexity information
    - `length` - Minimum Password length.
    - `caps` Number of uppercase characters.
    - `lowercase` Number of lowercase characters.
    - `digits` Number of digits.
    - `special` Number of special characters.
  - `value` - Field value
* `custom` - A list containing custom fields information:
  - `type` - Field type
  - `label` - Field label
  - `required` - Required field flag
  - `privacy_screen` - Privacy screen flag
  - `enforce_generation` - Enforce generation flag (for password field)
  - `complexity` - A list containing password complexity information
    - `length` - Minimum Password length.
    - `caps` Number of uppercase characters.
    - `lowercase` Number of lowercase characters.
    - `digits` Number of digits.
    - `special` Number of special characters.
  - `value` - Field value
* `file_ref` - A list containing file reference information:
  - `uid` - File UID
  - `title` - File title
  - `name` - File name
  - `type` - File content type
  - `size` - File size
  - `last_modified` - File last modification timestamp
  - `content_base64` - File content base64 encoded
