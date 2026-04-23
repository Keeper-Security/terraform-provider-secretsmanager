# secretsmanager_pam_machine (Ephemeral Resource)

Use this ephemeral resource to read secrets of type `pamMachine` stored in Keeper Vault.

Unlike data sources, ephemeral resources do not store any secret values in the Terraform state file. The values are only available during the Terraform plan and apply phases, making this a more secure option for accessing sensitive credentials.

## Example Usage

```terraform
ephemeral "secretsmanager_pam_machine" "server" {
  path = "<record UID>"
}

output "ssh_hostname" {
  value     = ephemeral.secretsmanager_pam_machine.server.pam_hostname[0].host_name
  ephemeral = true
}

output "machine_login" {
  value     = ephemeral.secretsmanager_pam_machine.server.login
  ephemeral = true
}
```

## Argument Reference

* `path` - (Optional) The UID of an existing record in Keeper Vault. Exactly one of `path` or `title` must be set.
* `title` - (Optional) The title of the record to search for. Exactly one of `path` or `title` must be set.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `type` - The type of the record (`pamMachine`).
* `title` - Record title.
* `notes` - Record notes.
* `folder_uid` - The UID of the folder where the record is stored.
* `pam_hostname` - PAM Hostname field. Each entry contains:
  - `host_name` - Hostname or IP address.
  - `port` - Port number.
* `pam_settings` - PAM connection settings as a JSON string.
* `login` - Login field containing the username.
* `password` - Password field (sensitive).
* `private_pem_key` - Secret field containing the private PEM key. Label: "Private PEM Key".
* `private_key_passphrase` - Secret field containing the passphrase used for private key protection. Label: "Private Key Passphrase".
* `rotation_scripts` - Script field for rotation scripts. Label: "Rotation Scripts".
* `operating_system` - Text field for the operating system. Label: "Operating System".
* `ssl_verification` - Checkbox field indicating whether SSL verification is enabled. Label: "SSL Verification".
* `instance_name` - Text field for the instance name. Label: "Instance Name".
* `instance_id` - Text field for the instance identifier. Label: "Instance Id".
* `provider_group` - Text field for provider group. Label: "Provider Group".
* `provider_region` - Text field for provider region. Label: "Provider Region".
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