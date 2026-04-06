# secretsmanager_record (Ephemeral Resource)

Use this ephemeral resource to read secrets of any type stored in Keeper Vault.

Unlike data sources, ephemeral resources do not store any secret values in the Terraform state file. The values are only available during the Terraform plan and apply phases, making this a more secure option for accessing sensitive credentials.

## Example Usage

```terraform
ephemeral "secretsmanager_record" "my_record" {
  path = "<record UID>"
}

output "first_field_value" {
  value     = length(ephemeral.secretsmanager_record.my_record.fields) < 1 ? "" : ephemeral.secretsmanager_record.my_record.fields.0.value
  ephemeral = true
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
