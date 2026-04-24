# secretsmanager_pam_database Resource

Use this resource to create and manage PAM Database records in Keeper Vault. Supports PostgreSQL, MySQL, MongoDB, and other database types for privileged access management.

## Example Usage

### PostgreSQL Database

```terraform
resource "secretsmanager_pam_database" "postgres" {
  folder_uid = "<folder UID>"
  title      = "Production PostgreSQL"

  pam_hostname {
    value {
      hostname = "postgres.prod.example.com"
      port     = "5432"
    }
  }

  database_type = "postgresql"

  pam_settings = jsonencode([{
    connection = [{
      protocol = "postgresql"
      port     = "5432"
    }]
  }])
}
```

### Database with Custom Fields

```terraform
resource "secretsmanager_pam_database" "mysql_staging" {
  folder_uid = "<folder UID>"
  title      = "Staging MySQL"

  pam_hostname {
    value {
      hostname = "mysql.staging.example.com"
      port     = "3306"
    }
  }

  database_type = "mysql"

  custom {
    type  = "text"
    label = "Environment"
    value = "staging"
  }

  custom {
    type  = "text"
    label = "Team"
    value = "platform"
  }
}
```

## Argument Reference

* `folder_uid` - (Optional) The UID of the shared folder where the record will be created. At least one of `folder_uid` or `uid` must be set.
* `uid` - (Optional) The UID for the new record (RFC 4648 URL-safe base64). Auto-generated if not set. At least one of `folder_uid` or `uid` must be set.
* `title` - (Optional) The record title.
* `notes` - (Optional) The record notes.
* `pam_hostname` - (Optional) Database hostname and port.
* `database_type` - (Optional) Database type string (e.g. `postgresql`, `mysql`, `mongodb`).
* `database_id` - (Optional) Database identifier (e.g. RDS instance ID). Block with `value` attribute.
* `use_ssl` - (Optional) SSL enabled flag. Block with `value` (boolean) attribute.
* `pam_settings` - (Optional) Connection and port-forward settings as a JSON string. Use `jsonencode()`.
* `rotation_scripts` - (Optional) Rotation script references.
* `provider_group` - (Optional) Cloud provider group. Block with `value` attribute.
* `provider_region` - (Optional) Cloud provider region. Block with `value` attribute.
* `file_ref` - (Optional) File references.
* `totp` - (Optional) One-time code (otpauth:// URI).
* `custom` - (Optional) User-defined custom fields. Each block requires `type` (Keeper field type) and `label` (display name), with optional `value` (plain string or `jsonencode()` for complex types), `required`, and `privacy_screen`. See [Nested Schema for `custom`](#nestedblock--custom) below.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `type` - The record type (`pamDatabase`).

## Import

PAM Database records can be imported using their UID:

```
$ terraform import secretsmanager_pam_database.example <record_UID>
```

<a id="nestedblock--custom"></a>
### Nested Schema for `custom`

Required:

- **label** (String) Display name for the field in Keeper UI.
- **type** (String) Keeper field type. Input is case-insensitive — any casing is accepted and normalized (e.g., `paymentcard` → `paymentCard`). Unknown types are rejected at plan time. Common values: `text`, `secret`, `url`, `email`, `phone`, `date`, `birthDate`, `expirationDate`, `name`, `address`, `paymentCard`, `bankAccount`, `host`, `keyPair`, `securityQuestion`, `checkbox`, `multiline`.

Optional:

- **privacy_screen** (Boolean) Whether this field is hidden behind a privacy screen in the Keeper UI.
- **required** (Boolean) Whether this field is required.
- **value** (String, Sensitive) Field value. Plain string for simple types. Use `jsonencode({...})` for structured types or `jsonencode([{...},{...}])` for multiple entries in one field. Format constraints: `checkbox` requires `"true"` or `"false"`; `date`, `birthDate`, and `expirationDate` require YYYY-MM-DD; `paymentCard` `jsonencode` keys use camelCase (`cardNumber`, `cardExpirationDate`, `cardSecurityCode`).
