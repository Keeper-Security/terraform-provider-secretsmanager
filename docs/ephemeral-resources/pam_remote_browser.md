# secretsmanager_pam_remote_browser (Ephemeral Resource)

Use this ephemeral resource to read a PAM Remote Browser record from Keeper Secrets Manager. Ephemeral resources do not store secret values in Terraform state.

## Example Usage

```terraform
ephemeral "secretsmanager_pam_remote_browser" "browser" {
  path = "<record UID>"
}
```

## Argument Reference

* `path` - (Required) The UID of the PAM Remote Browser record.
* `title` - (Optional) The secret title. Used with `path = "*"` to search by title.

## Attributes Reference

* `type` - The type of the record (`pamRemoteBrowser`).
* `title` - Record title.
* `notes` - Record notes.
* `rbi_url` - The Remote Browser Interface URL.
* `pam_remote_browser_settings` - PAM Remote Browser connection settings as a JSON string.
* `traffic_encryption_seed` - The base64-encoded traffic encryption seed (sensitive).
* `file_ref` - A list containing file reference information:
  - `uid` - File UID.
  - `title` - File title.
  - `name` - File name.
  - `type` - File content type.
  - `size` - File size.
  - `last_modified` - File last modification timestamp.
  - `content_base64` - File content base64 encoded.
* `totp` - One-time password information:
  - `url` - TOTP URL.
  - `token` - Generated TOTP token (sensitive).
  - `ttl` - Time to live in seconds.
