# secretsmanager_ssn_card (Ephemeral Resource)

Use this ephemeral resource to read secrets of type `ssnCard` stored in Keeper Vault.

Unlike data sources, ephemeral resources do not store any secret values in the Terraform state file. The values are only available during the Terraform plan and apply phases, making this a more secure option for accessing sensitive credentials.

## Example Usage

```terraform
ephemeral "secretsmanager_ssn_card" "my_ssn" {
  path = "<record UID>"
}

output "identity_number" {
  value     = ephemeral.secretsmanager_ssn_card.my_ssn.identity_number
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
* `identity_number` - Identity Number.
* `name` - A list containing name information:
  - `first` - First name
  - `middle` - Middle name
  - `last` - Last name
* `file_ref` - A list containing file reference information:
  - `uid` - File UID
  - `title` - File title
  - `name` - File name
  - `type` - File content type
  - `size` - File size
  - `last_modified` - File last modification timestamp
  - `content_base64` - File content base64 encoded
