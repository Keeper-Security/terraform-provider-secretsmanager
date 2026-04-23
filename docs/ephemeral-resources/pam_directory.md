# secretsmanager_pam_directory (Ephemeral Resource)

Use this ephemeral resource to read secrets of type `pamDirectory` stored in Keeper Vault.

Unlike data sources, ephemeral resources do not store any secret values in the Terraform state file. The values are only available during the Terraform plan and apply phases, making this a more secure option for accessing sensitive credentials.

## Example Usage

```terraform
ephemeral "secretsmanager_pam_directory" "ad" {
  path = "<record UID>"
}

output "ad_hostname" {
  value     = ephemeral.secretsmanager_pam_directory.ad.pam_hostname[0].host_name
  ephemeral = true
}

output "ad_directory_type" {
  value     = ephemeral.secretsmanager_pam_directory.ad.directory_type
  ephemeral = true
}
```

## Argument Reference

* `path` - (Optional) The UID of an existing record in Keeper Vault. Exactly one of `path` or `title` must be set.
* `title` - (Optional) The title of the record to search for. Exactly one of `path` or `title` must be set.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `type` - The type of the record (`pamDirectory`).
* `title` - Record title.
* `notes` - Record notes.
* `folder_uid` - The UID of the folder where the record is stored.
* `pam_hostname` - PAM Hostname field. Each entry contains:
  - `host_name` - Hostname or IP address.
  - `port` - Port number.
* `pam_settings` - PAM connection settings as a JSON string.
* `directory_type` - Directory type (e.g. `Active Directory`, `OpenLDAP`).
* `use_ssl` - Checkbox field indicating whether SSL is enabled.
* `rotation_scripts` - Script field for rotation scripts. Label: "Rotation Scripts".
* `distinguished_name` - Text field for the distinguished name. Label: "Distinguished Name".
* `domain_name` - Text field for the domain name. Label: "domainName".
* `directory_id` - Text field for the directory identifier. Label: "directoryId".
* `user_match` - Text field for user match attribute. Label: "userMatch".
* `provider_group` - Text field for provider group. Label: "providerGroup".
* `provider_region` - Text field for provider region. Label: "providerRegion".
* `alternative_ips` - Multiline field for alternative IP addresses. Label: "alternativeIPs".
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

* `custom` - A list of custom fields defined on the record. Each entry contains:
  - `type` - Field type (e.g. `text`, `secret`, `url`, `email`, `phone`, `multiline`, `checkbox`, `date`, `birthDate`, `expirationDate`, `name`, `address`, `paymentCard`, `bankAccount`, `host`, `keyPair`, `securityQuestion`).
  - `label` - Display name for the field in Keeper UI.
  - `value` - Field value. Complex types (e.g. `name`, `address`, `paymentCard`) are returned as a JSON-encoded string.
  - `required` - Whether this field is required.
  - `privacy_screen` - Whether this field is hidden behind a privacy screen in Keeper UI.