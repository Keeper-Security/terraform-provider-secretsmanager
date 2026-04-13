# secretsmanager_driver_license Resource

Use this resource to create and manage secrets of type `driverLicense` in Keeper Vault

## Example Usage

```terraform
resource "secretsmanager_driver_license" "my_driver_license" {
  folder_uid = "<folder UID>"
  title      = "My Title"
  notes      = "My Notes"

  driver_license_number {
    label          = "My Driver License"
    required       = true
    privacy_screen = true
    value          = "My Driver License# 1234"
  }

  name {
    label          = "John"
    required       = true
    privacy_screen = true
    value {
      first  = "John"
      middle = "D"
      last   = "Doe"
    }
  }

  birth_date {
    label          = "Birth Date"
    required       = true
    privacy_screen = true
    value          = 1651186276
  }

  expiration_date {
    label          = "Driver License Expiration Date"
    required       = true
    privacy_screen = true
    value          = 21651186276
  }

  address_ref {
    label          = "My Address Ref"
    required       = true
    privacy_screen = true
    value          = "<address ref UID>"
  }
}
```

## Schema

### Optional

- **address_ref** (Block List, Max: 1) AddressRef field data. (see [below for nested schema](#nestedblock--address_ref))
- **birth_date** (Block List, Max: 1) Birth date field data. (see [below for nested schema](#nestedblock--birth_date))
- **driver_license_number** (Block List, Max: 1) Account number field data. (see [below for nested schema](#nestedblock--driver_license_number))
- **expiration_date** (Block List, Max: 1) Expiration date field data. (see [below for nested schema](#nestedblock--expiration_date))
- **file_ref** (Block List, Max: 1) FileRef field data. (see [below for nested schema](#nestedblock--file_ref))
- **folder_uid** (String) The folder UID where the secret is stored. The parent shared folder must be non empty.
- **id** (String) The ID of this resource.
- **name** (Block List, Max: 1) Name field data. (see [below for nested schema](#nestedblock--name))
- **notes** (String) The secret notes.
- **title** (String) The secret title.
- **uid** (String) The UID of the new secret (using RFC4648 URL and Filename Safe Alphabet).

- **custom** (Block List) User-defined custom fields. (see [below for nested schema](#nestedblock--custom))

### Read-Only

- **type** (String) The secret type.

<a id="nestedblock--address_ref"></a>
### Nested Schema for `address_ref`

Optional:

- **label** (String) Field label.
- **privacy_screen** (Boolean) Privacy screen flag.
- **required** (Boolean) Required flag.
- **value** (String) Field value.

Read-Only:

- **type** (String) Field type.

<a id="nestedblock--birth_date"></a>
### Nested Schema for `birth_date`

Optional:

- **label** (String) Field label.
- **privacy_screen** (Boolean) Privacy screen flag.
- **required** (Boolean) Required flag.
- **value** (Number) Field value.

Read-Only:

- **type** (String) Field type.

<a id="nestedblock--driver_license_number"></a>
### Nested Schema for `driver_license_number`

Optional:

- **label** (String) Field label.
- **privacy_screen** (Boolean) Privacy screen flag.
- **required** (Boolean) Required flag.
- **value** (String, Sensitive) Field value.

Read-Only:

- **type** (String) Field type.

<a id="nestedblock--expiration_date"></a>
### Nested Schema for `expiration_date`

Optional:

- **label** (String) Field label.
- **privacy_screen** (Boolean) Privacy screen flag.
- **required** (Boolean) Required flag.
- **value** (Number) Field value.

Read-Only:

- **type** (String) Field type.

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

<a id="nestedblock--name"></a>
### Nested Schema for `name`

Optional:

- **label** (String) Field label.
- **privacy_screen** (Boolean) Privacy screen flag.
- **required** (Boolean) Required flag.
- **value** (Block List, Max: 1) Field value. (see [below for nested schema](#nestedblock--name--value))

Read-Only:

- **type** (String) Field type.

<a id="nestedblock--name--value"></a>
### Nested Schema for `name.value`

Optional:

- **first** (String) First name.
- **last** (String) Last name.
- **middle** (String) MIddle name.

<a id="nestedblock--custom"></a>
### Nested Schema for `custom`

Required:

- **label** (String) Display name for the field in Keeper UI.
- **type** (String) Keeper field type. Common values: `text`, `secret`, `url`, `email`, `phone`, `date`, `birthDate`, `expirationDate`, `name`, `address`, `paymentCard`, `bankAccount`, `host`, `keyPair`, `securityQuestion`, `checkbox`, `multiline`.

Optional:

- **privacy_screen** (Boolean) Whether this field is hidden behind a privacy screen in the Keeper UI.
- **required** (Boolean) Whether this field is required.
- **value** (String, Sensitive) Field value. Plain string for simple types. Use `jsonencode({...})` for structured types or `jsonencode([{...},{...}])` for multiple entries in one field. Format constraints: `checkbox` requires `"true"` or `"false"`; `date`, `birthDate`, and `expirationDate` require YYYY-MM-DD; `paymentCard` `jsonencode` keys use camelCase (`cardNumber`, `cardExpirationDate`, `cardSecurityCode`).
