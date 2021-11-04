# secretsmanager_ssh_keys Data Source

Use this data source to read secrets of type `sshKeys` stored in Keeper Vault

## Example Usage

```terraform
data "secretsmanager_ssh_keys" "ssh_keys" {
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
* `login` - The secret login.
* `key_pair` - A list containing public and private key pair information:
  - `public_key` - The public key.
  - `private_key` - The private key.
* `passphrase` - The passphrase to unlock the key.
* `host` - A list containing hostname and port information:
  - `host_name` - Database server hostname
  - `port` - Port number
* `file_ref` - A list containing file reference information:
  - `uid` - File UID
  - `title` - File title
  - `name` - File name
  - `type` - File content type
  - `size` - File size
  - `last_modified` - File last modification timestamp
  - `url` - File download URL
  - `content_base64` - File content base64 encoded
