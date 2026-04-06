# secretsmanager_pam_user (Ephemeral Resource)

Use this ephemeral resource to read secrets of type `pamUser` stored in Keeper Vault.

Unlike data sources, ephemeral resources do not store any secret values in the Terraform state file. The values are only available during the Terraform plan and apply phases, making this a more secure option for accessing sensitive credentials.

## Example Usage

```terraform
ephemeral "secretsmanager_pam_user" "admin" {
  path = "<record UID>"
}

output "db_login" {
  value     = ephemeral.secretsmanager_pam_user.admin.login[0].value
  ephemeral = true
}

output "db_password" {
  value     = ephemeral.secretsmanager_pam_user.admin.password[0].value
  ephemeral = true
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
* `private_key_passphrase` - Secret field containing the passphrase used for private key protection. Label: "Private Key Passphrase".
* `distinguished_name` - Text field for the distinguished name. Label: "Distinguished Name".
* `connect_database` - Text field for the database to connect to. Label: "Connect Database".
* `managed` - Checkbox field indicating whether the user is managed. Label: "Managed".
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
