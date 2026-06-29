# secretsmanager_metadata Data Source

Use this data source to read non-sensitive metadata about a record stored in Keeper Vault without exposing its field values.

This is useful for pairing with write-only attributes on other providers: the `revision` attribute changes whenever the Keeper record is modified, so it can be used as a version signal that triggers Terraform to re-apply a write-only value when the secret rotates.

## Example Usage

```terraform
data "secretsmanager_metadata" "by_uid" {
  path = "<record UID>"
}

data "secretsmanager_metadata" "by_title" {
  path  = "*"
  title = "<record title>"
}
```

## Argument Reference

* `path` - (Required) The UID of an existing record in Keeper Vault, or `*` to look up the record by title.
* `title` - (Optional) The record title. Used to look up the record when `path` is `*`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `uid` - The record UID.
* `type` - The record type (e.g. `login`).
* `title` - The record title.
* `notes` - The record notes.
* `revision` - Record revision counter. Increments on every modification (value, title, and notes changes all bump it).
* `folder_uid` - The UID of the folder containing this record.
* `is_editable` - Whether the KSM application credential can edit this record.
