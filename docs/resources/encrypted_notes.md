# secretsmanager_encrypted_notes Resource

Use this resource to create and manage secrets of type `encryptedNotes` in Keeper Vault

## Example Usage

```terraform
resource "secretsmanager_encrypted_notes" "my_encrypted_notes" {
  folder_uid = "<folder UID>"
  title      = "My Title"
  notes      = "My Notes"

  note {
    label          = "My Note"
    required       = true
    privacy_screen = true
    value          = "My Note"
  }

  date {
    label          = "Date"
    required       = true
    privacy_screen = true
    value          = 1651186276
  }
}
```

## Schema

### Optional

- **date** (Block List, Max: 1) Date field data. (see [below for nested schema](#nestedblock--date))
- **file_ref** (Block List, Max: 1) FileRef field data. (see [below for nested schema](#nestedblock--file_ref))
- **folder_uid** (String) The folder UID where the secret is stored. The parent shared folder must be non empty.
- **id** (String) The ID of this resource.
- **note** (Block List, Max: 1) Secure note field data. (see [below for nested schema](#nestedblock--note))
- **notes** (String) The secret notes.
- **title** (String) The secret title.
- **uid** (String) The UID of the new secret (using RFC4648 URL and Filename Safe Alphabet).

- **custom** (Block List) User-defined custom fields. (see [below for nested schema](#nestedblock--custom))

### Read-Only

- **type** (String) The secret type.

<a id="nestedblock--date"></a>
### Nested Schema for `date`

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

<a id="nestedblock--note"></a>
### Nested Schema for `note`

Optional:

- **label** (String) Field label.
- **privacy_screen** (Boolean) Privacy screen flag.
- **required** (Boolean) Required flag.
- **value** (String, Sensitive) Field value.

Read-Only:

- **type** (String) Field type.

<a id="nestedblock--custom"></a>
### Nested Schema for `custom`

Required:

- **label** (String) Display name for the field in Keeper UI.
- **type** (String) Keeper field type. Common values: `text`, `secret`, `url`, `email`, `phone`, `date`, `birthDate`, `expirationDate`, `name`, `address`, `paymentCard`, `bankAccount`, `host`, `keyPair`, `securityQuestion`, `checkbox`, `multiline`.

Optional:

- **privacy_screen** (Boolean) Whether this field is hidden behind a privacy screen in the Keeper UI.
- **required** (Boolean) Whether this field is required.
- **value** (String, Sensitive) Field value. Plain string for simple types. Use `jsonencode({...})` for structured types or `jsonencode([{...},{...}])` for multiple entries in one field.
