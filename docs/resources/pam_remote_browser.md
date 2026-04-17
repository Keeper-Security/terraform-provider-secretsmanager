# secretsmanager_pam_remote_browser Resource

Use this resource to create and manage PAM Remote Browser records in Keeper Vault. Supports remote browser isolation (RBI) for secure web access through privileged access management.

## Example Usage

### Basic PAM Remote Browser

```terraform
resource "secretsmanager_pam_remote_browser" "basic" {
  folder_uid = "<folder UID>"
  title      = "Internal Web App"

  rbi_url {
    value = "https://internal-app.example.com"
  }
}
```

### PAM Remote Browser with Settings

```terraform
resource "secretsmanager_pam_remote_browser" "with_settings" {
  folder_uid = "<folder UID>"
  title      = "Secure Browser Session"

  rbi_url {
    value = "https://admin.example.com"
  }

  pam_remote_browser_settings = jsonencode({
    "connection" = {
      "protocol"            = "http"
      "allowUrlManipulation" = false
    }
  })
}
```

## Argument Reference

* `folder_uid` - (Optional) The UID of the shared folder where the record will be created. At least one of `folder_uid` or `uid` must be set.
* `uid` - (Optional) The UID for the new record (RFC 4648 URL-safe base64). Auto-generated if not set. At least one of `folder_uid` or `uid` must be set.
* `title` - (Optional) The record title.
* `notes` - (Optional) The record notes.
* `rbi_url` - (Optional) The Remote Browser Interface URL. Block with `value` attribute.
* `pam_remote_browser_settings` - (Optional) Connection settings as a JSON string.
* `traffic_encryption_seed` - (Optional) Base64-encoded 256-bit encryption seed. Block with `value` attribute.
* `file_ref` - (Optional) File references.
* `totp` - (Optional) One-time code (otpauth:// URI).
* `custom` - (Optional) User-defined custom fields. Each block requires `type` (Keeper field type) and `label` (display name), with optional `value` (plain string or `jsonencode()` for complex types), `required`, and `privacy_screen`. See [Nested Schema for `custom`](#nestedblock--custom) below.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `type` - The record type (`pamRemoteBrowser`).

## Import

PAM Remote Browser records can be imported using their UID:

```
$ terraform import secretsmanager_pam_remote_browser.example <record_UID>
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
