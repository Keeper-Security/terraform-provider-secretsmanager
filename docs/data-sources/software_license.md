# keeper_secret_software_license Data Source

Use this data source to read secrets of type `softwareLicense` stored in Keeper Vault

## Example Usage

```terraform
data "keeper_secret_software_license" "software_license" {
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
* `license_number` - License Number.
* `activation_date` - Date of activation.
* `expiration_date` - Date of expiration.
* `file_ref` - A list containing file reference information:
  - `uid` - File UID
  - `title` - File title
  - `name` - File name
  - `type` - File content type
  - `size` - File size
  - `last_modified` - File last modification timestamp
  - `url` - File download URL
  - `content_base64` - File content base64 encoded
