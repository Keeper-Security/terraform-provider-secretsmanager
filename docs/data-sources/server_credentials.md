# keeper_secret_server_credentials Data Source

Use this data source to read secrets of type `serverCredentials` stored in Keeper Vault

## Example Usage

```terraform
data "keeper_secret_server_credentials" "server_credentials" {
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
* `host` - A list containing hostname and port information:
  - `host_name` - Database server hostname
  - `port` - Port number
* `login` - The secret login.
* `password` - The secret password.
* `file_ref` - A list containing file reference information:
  - `uid` - File UID
  - `title` - File title
  - `name` - File name
  - `type` - File content type
  - `size` - File size
  - `last_modified` - File last modification timestamp
  - `url` - File download URL
  - `content_base64` - File content base64 encoded
