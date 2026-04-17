# secretsmanager_address Resource

Use this resource to create and manage secrets of type `address` in Keeper Vault

## Example Usage

```terraform
resource "secretsmanager_address" "my_address" {
  folder_uid = "<folder UID>"
  title      = "My Title"
  notes      = "My Notes"

  address {
    label          = "My Address"
    required       = true
    privacy_screen = true
    value {
      street1 = "7422 Avalon"
      street2 = "Apt 21"
      city    = "Los Angeles"
      state   = "CA"
      country = "United States"
      zip     = "90003-2334"
    }
  }
}
```

## Schema

### Optional

- **address** (Block List, Max: 1) Address field data. (see [below for nested schema](#nestedblock--address))
- **file_ref** (Block List, Max: 1) FileRef field data. (see [below for nested schema](#nestedblock--file_ref))
- **folder_uid** (String) The folder UID where the secret is stored. The parent shared folder must be non empty.
- **notes** (String) The secret notes.
- **title** (String) The secret title.
- **uid** (String) The UID of the new secret (using RFC4648 URL and Filename Safe Alphabet).

- **custom** (Block List) User-defined custom fields. (see [below for nested schema](#nestedblock--custom))

### Read-Only

- **type** (String) The secret type.

<a id="nestedblock--address"></a>
### Nested Schema for `address`

Optional:

- **label** (String) Field label.
- **privacy_screen** (Boolean) Privacy screen flag.
- **required** (Boolean) Required flag.
- **value** (Block List) Field value. (see [below for nested schema](#nestedblock--address--value))

Read-Only:

- **type** (String) Field type.

<a id="nestedblock--address--value"></a>
### Nested Schema for `address.value`

Optional:

- **city** (String) City.
- **country** (String) Country.
- **state** (String) State.
- **street1** (String) Street line 1.
- **street2** (String) Street line 2.
- **zip** (String) ZIP code.

<a id="nestedblock--file_ref"></a>
### Nested Schema for `file_ref`

Optional:

- **label** (String) Field label.
- **required** (Boolean) Required flag.
- **value** (Block List) Field value (File UID list). (see [below for nested schema](#nestedblock--file_ref--value))

Read-Only:

- **type** (String) Field type.

<a id="nestedblock--file_ref--value"></a>
### Nested Schema for `file_ref.value`

Optional:

- **uid** (String) The file ref UID.

Read-Only:

- **content_base64** (String) The file content (base64).
- **last_modified** (String) The file last modified date.
- **name** (String) The file name.
- **size** (Number) The file size.
- **title** (String) The file title.
- **type** (String) The file type.

<a id="nestedblock--custom"></a>
### Nested Schema for `custom`

Required:

- **label** (String) Display name for the field in Keeper UI.
- **type** (String) Keeper field type. Input is case-insensitive — any casing is accepted and normalized (e.g., `paymentcard` → `paymentCard`). Unknown types are rejected at plan time. Common values: `text`, `secret`, `url`, `email`, `phone`, `date`, `birthDate`, `expirationDate`, `name`, `address`, `paymentCard`, `bankAccount`, `host`, `keyPair`, `securityQuestion`, `checkbox`, `multiline`.

Optional:

- **privacy_screen** (Boolean) Whether this field is hidden behind a privacy screen in the Keeper UI.
- **required** (Boolean) Whether this field is required.
- **value** (String, Sensitive) Field value. Plain string for simple types. Use `jsonencode({...})` for structured types or `jsonencode([{...},{...}])` for multiple entries in one field. Format constraints: `checkbox` requires `"true"` or `"false"`; `date`, `birthDate`, and `expirationDate` require YYYY-MM-DD; `paymentCard` `jsonencode` keys use camelCase (`cardNumber`, `cardExpirationDate`, `cardSecurityCode`).
