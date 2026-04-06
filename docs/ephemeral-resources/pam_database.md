# secretsmanager_pam_database (Ephemeral Resource)

Use this ephemeral resource to read secrets of type `pamDatabase` stored in Keeper Vault.

Unlike data sources, ephemeral resources do not store any secret values in the Terraform state file. The values are only available during the Terraform plan and apply phases, making this a more secure option for accessing sensitive credentials.

## Example Usage

```terraform
ephemeral "secretsmanager_pam_database" "db" {
  path = "<record UID>"
}

output "db_hostname" {
  value     = ephemeral.secretsmanager_pam_database.db.pam_hostname[0].value[0].hostname
  ephemeral = true
}

output "db_type" {
  value     = ephemeral.secretsmanager_pam_database.db.database_type
  ephemeral = true
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
* `totp` - One-time code field represented as a block list with:
  - `value` - TOTP URI/secret value (e.g. `otpauth://...`).
