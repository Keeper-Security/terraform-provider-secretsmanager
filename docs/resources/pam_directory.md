# secretsmanager_pam_directory Resource

Use this resource to create and manage PAM Directory records in Keeper Vault. Supports Active Directory, OpenLDAP, and other LDAP-compatible directories for privileged access management.

## Example Usage

### Active Directory

```terraform
resource "secretsmanager_pam_directory" "ad" {
  folder_uid = "<folder UID>"
  title      = "Corporate Active Directory"

  pam_hostname {
    value {
      hostname = "dc.corp.example.com"
      port     = "636"
    }
  }

  directory_type = "Active Directory"

  distinguished_name {
    label = "Base DN"
    value = "DC=corp,DC=example,DC=com"
  }

  use_ssl {
    value = true
  }
}
```

### Directory with Custom Fields

```terraform
resource "secretsmanager_pam_directory" "openldap" {
  folder_uid = "<folder UID>"
  title      = "OpenLDAP Staging"

  pam_hostname {
    value {
      hostname = "ldap.staging.example.com"
      port     = "389"
    }
  }

  directory_type = "OpenLDAP"

  custom {
    type  = "text"
    label = "Environment"
    value = "staging"
  }
}
```

## Argument Reference

* `folder_uid` - (Optional) The UID of the shared folder where the record will be created. At least one of `folder_uid` or `uid` must be set.
* `uid` - (Optional) The UID for the new record (RFC 4648 URL-safe base64). Auto-generated if not set. At least one of `folder_uid` or `uid` must be set.
* `title` - (Optional) The record title.
* `notes` - (Optional) The record notes.
* `pam_hostname` - (Optional) Directory server hostname and port.
* `directory_type` - (Optional) Directory type string (e.g. `Active Directory`, `OpenLDAP`).
* `distinguished_name` - (Optional) Base DN or search base. Block with `label` and `value` attributes.
* `domain_name` - (Optional) Domain name. Block with `value` attribute.
* `directory_id` - (Optional) Directory identifier. Block with `value` attribute.
* `user_match` - (Optional) User match filter. Block with `value` attribute.
* `use_ssl` - (Optional) SSL/TLS enabled flag. Block with `value` (boolean) attribute.
* `alternative_ips` - (Optional) Alternative IP addresses. Block with `value` attribute.
* `pam_settings` - (Optional) Connection and port-forward settings as a JSON string. Use `jsonencode()`.
* `rotation_scripts` - (Optional) Rotation script references.
* `provider_group` - (Optional) Cloud provider group. Block with `value` attribute.
* `provider_region` - (Optional) Cloud provider region. Block with `value` attribute.
* `file_ref` - (Optional) File references.
* `totp` - (Optional) One-time code (otpauth:// URI).
* `custom` - (Optional) User-defined custom fields. Each block requires `type` (Keeper field type) and `label` (display name), with optional `value` (plain string or `jsonencode()` for complex types), `required`, and `privacy_screen`. See [Nested Schema for `custom`](#nestedblock--custom) below.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `type` - The record type (`pamDirectory`).

## Import

PAM Directory records can be imported using their UID:

```
$ terraform import secretsmanager_pam_directory.example <record_UID>
```

<a id="nestedblock--custom"></a>
### Nested Schema for `custom`

Required:

- **label** (String) Display name for the field in Keeper UI.
- **type** (String) Keeper field type. Common values: `text`, `secret`, `url`, `email`, `phone`, `date`, `birthDate`, `expirationDate`, `name`, `address`, `paymentCard`, `bankAccount`, `host`, `keyPair`, `securityQuestion`, `checkbox`, `multiline`.

Optional:

- **privacy_screen** (Boolean) Whether this field is hidden behind a privacy screen in the Keeper UI.
- **required** (Boolean) Whether this field is required.
- **value** (String, Sensitive) Field value. Plain string for simple types. Use `jsonencode({...})` for structured types or `jsonencode([{...},{...}])` for multiple entries in one field. Format constraints: `checkbox` requires `"true"` or `"false"`; `date`, `birthDate`, and `expirationDate` require YYYY-MM-DD; `paymentCard` `jsonencode` keys use camelCase (`cardNumber`, `cardExpirationDate`, `cardSecurityCode`).
