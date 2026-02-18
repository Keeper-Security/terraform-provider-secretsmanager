# secretsmanager_pam_user Data Source

Use this data source to read secrets of type `pamUser` stored in Keeper Vault.

## Example Usage

```terraform
data "secretsmanager_pam_user" "admin" {
  path = "<record UID>"
}
```

## Argument Reference

* `path` - (Optional) The UID of an existing record in Keeper Vault. Exactly one of `path` or `title` must be set.
* `title` - (Optional) The title of the record to search for. Exactly one of `path` or `title` must be set.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `type` - The type of the record (`pamUser`).
* `title` - Record title.
* `notes` - Record notes.
* `folder_uid` - The UID of the folder where the record is stored.
* `login` - Login field containing the username.
* `password` - Password field (sensitive).
* `rotation_scripts` - Script field for rotation scripts. Label: "Rotation Scripts".
* `private_pem_key` - Secret field containing the private PEM key. Label: "Private PEM Key".
* `distinguished_name` - Text field for the distinguished name. Label: "Distinguished Name".
* `connect_database` - Text field for the database to connect to. Label: "connectDatabase".
* `managed` - Checkbox field indicating whether the user is managed. Label: "Managed".
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
