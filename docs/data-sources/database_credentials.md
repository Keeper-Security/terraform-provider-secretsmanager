# secretsmanager_database_credentials Data Source

Use this data source to read secrets of type `databaseCredentials` stored in Keeper Vault

## Example Usage

```terraform
data "secretsmanager_database_credentials" "database_credentials" {
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
* `db_type` - Database type.
* `login` - Login name.
* `password` - Login password.
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
  - `content_base64` - File content base64 encoded
