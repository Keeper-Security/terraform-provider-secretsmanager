# secretsmanager_pam_remote_browser Data Source

Use this data source to read secrets of type `pamRemoteBrowser` stored in Keeper Vault.

## Example Usage

```terraform
data "secretsmanager_pam_remote_browser" "browser" {
  path = "<record UID>"
}
```

## Argument Reference

* `path` - (Optional) The UID of an existing record in Keeper Vault. Exactly one of `path` or `title` must be set.
* `title` - (Optional) The title of the record to search for. Exactly one of `path` or `title` must be set.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `type` - The type of the record (`pamRemoteBrowser`).
* `title` - Record title.
* `notes` - Record notes.
* `folder_uid` - The UID of the folder where the record is stored.
* `rbi_url` - Text field containing the Remote Browser Interface URL.
* `pam_remote_browser_settings` - PAM Remote Browser connection settings as a JSON string.
* `traffic_encryption_seed` - Text field containing the base64-encoded traffic encryption seed.
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
