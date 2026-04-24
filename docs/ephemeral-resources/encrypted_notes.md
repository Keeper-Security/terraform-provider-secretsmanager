# secretsmanager_encrypted_notes (Ephemeral Resource)

Use this ephemeral resource to read secrets of type `encryptedNotes` stored in Keeper Vault.

Unlike data sources, ephemeral resources do not store any secret values in the Terraform state file. The values are only available during the Terraform plan and apply phases, making this a more secure option for accessing sensitive credentials.

## Example Usage

```terraform
ephemeral "secretsmanager_encrypted_notes" "my_notes" {
  path = "<record UID>"
}

output "notes" {
  value     = ephemeral.secretsmanager_encrypted_notes.my_notes.note
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
* `note` - Encrypted note.
* `date` - Date of the note.
* `file_ref` - A list containing file reference information:
  - `uid` - File UID
  - `title` - File title
  - `name` - File name
  - `type` - File content type
  - `size` - File size
  - `last_modified` - File last modification timestamp
  - `content_base64` - File content base64 encoded

* `custom` - A list of custom fields defined on the record. Each entry contains:
  - `type` - Field type (e.g. `text`, `secret`, `url`, `email`, `phone`, `multiline`, `checkbox`, `date`, `birthDate`, `expirationDate`, `name`, `address`, `paymentCard`, `bankAccount`, `host`, `keyPair`, `securityQuestion`).
  - `label` - Display name for the field in Keeper UI.
  - `value` - Field value. Complex types (e.g. `name`, `address`, `paymentCard`) are returned as a JSON-encoded string.
  - `required` - Whether this field is required.
  - `privacy_screen` - Whether this field is hidden behind a privacy screen in Keeper UI.