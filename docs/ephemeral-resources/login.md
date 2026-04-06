# secretsmanager_login (Ephemeral Resource)

Use this ephemeral resource to read secrets of type `login` stored in Keeper Vault.

Unlike data sources, ephemeral resources do not store any secret values in the Terraform state file. The values are only available during the Terraform plan and apply phases, making this a more secure option for accessing sensitive credentials.

## Example Usage

```terraform
ephemeral "secretsmanager_login" "db_server" {
  path = "<record UID>"
}

output "login" {
  value     = ephemeral.secretsmanager_login.db_server.login
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
* `login` - The secret login.
* `password` - The secret password.
* `url` - The secret url.
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
