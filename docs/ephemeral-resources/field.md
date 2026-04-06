# secretsmanager_field (Ephemeral Resource)

Use this ephemeral resource to read fields of secrets of any type stored in Keeper Vault.

Unlike data sources, ephemeral resources do not store any secret values in the Terraform state file. The values are only available during the Terraform plan and apply phases, making this a more secure option for accessing sensitive credentials.

## Example Usage

```terraform
ephemeral "secretsmanager_field" "my_field" {
  path = "<record UID>/field/login"
}

output "field_value" {
  value     = ephemeral.secretsmanager_field.my_field.value
  ephemeral = true
}
```

## Argument Reference

* `path` - (Required) The path to a field of a secret stored in existing record in Keeper Vault. Provide full path to the field - regular fields are accessible by field type and custom fields are accessible by field label: ex. `<record UID>/field/login`, ex. `<record UID>/custom_field/custom1`, ex. `<record UID>/custom_field/custom2`. Use `*` in place of `<record UID>` in combination with `title` argument (_see below_) - to find the record by title (which then expands `*` to the actual `<record UID>`) ex. `*/field/login`

* `title` - (Optional) The title of a secret stored in existing record in Keeper Vault. If there's a need to find record by title - use `*` in place of `<record UID>`. If a single record is found by the title then `*` is expanded to the actual `<record UID>`

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `value` - The value of the selected field.
