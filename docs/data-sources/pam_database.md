# secretsmanager_pam_database Data Source

Use this data source to read secrets of type `pamDatabase` stored in Keeper Vault.

## Example Usage

```terraform
data "secretsmanager_pam_database" "db" {
  path = "<record UID>"
}
```

## Argument Reference

* `path` - (Optional) The UID of an existing record in Keeper Vault. Exactly one of `path` or `title` must be set.
* `title` - (Optional) The title of the record to search for. Exactly one of `path` or `title` must be set.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `type` - The type of the record (`pamDatabase`).
* `title` - Record title.
* `notes` - Record notes.
* `folder_uid` - The UID of the folder where the record is stored.
* `pam_hostname` - PAM Hostname field. Contains a `value` block with:
  - `hostname` - Hostname or IP address.
  - `port` - Port number.
* `pam_settings` - PAM connection settings as a JSON string.
* `use_ssl` - Checkbox field indicating whether SSL is enabled.
* `rotation_scripts` - Script field for rotation scripts. Label: "Rotation Scripts".
* `database_id` - Text field containing the database identifier. Label: "Database Id".
* `database_type` - Database type (e.g. `postgresql`, `mysql`, `mongodb`).
* `provider_group` - Text field for provider group. Label: "Provider Group".
* `provider_region` - Text field for provider region. Label: "Provider Region".
* `file_ref` - A list containing file reference information:
  - `uid` - File UID.
  - `title` - File title.
  - `name` - File name.
  - `type` - File content type.
  - `size` - File size.
  - `last_modified` - File last modification timestamp.
  - `content_base64` - File content base64 encoded.
* `totp` - A list containing Time-based One-time password information:
  - `url` - TOTP URL.
  - `token` - Current TOTP password.
  - `ttl` - Time to live in seconds for current token.
