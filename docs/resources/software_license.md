# secretsmanager_software_license Resource

Use this resource to access secrets of type `softwareLicense` stored in Keeper Vault

## Schema

### Optional

- **activation_date** (Block List, Max: 1) Date field data. (see [below for nested schema](#nestedblock--activation_date))
- **expiration_date** (Block List, Max: 1) Expiration date field data. (see [below for nested schema](#nestedblock--expiration_date))
- **file_ref** (Block List, Max: 1) FileRef field data. (see [below for nested schema](#nestedblock--file_ref))
- **folder_uid** (String) The folder UID where the secret is stored. The shared folder must be non empty.
- **id** (String) The ID of this resource.
- **license_number** (Block List, Max: 1) License number field data. (see [below for nested schema](#nestedblock--license_number))
- **notes** (String) The secret notes.
- **title** (String) The secret title.
- **uid** (String) The UID of the new secret (using RFC4648 URL and Filename Safe Alphabet).

### Read-Only

- **type** (String) The secret type.

<a id="nestedblock--activation_date"></a>
### Nested Schema for `activation_date`

Optional:

- **label** (String) Field label.
- **privacy_screen** (Boolean) Privacy screen flag.
- **required** (Boolean) Required flag.
- **value** (Number) Field value.

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

<a id="nestedblock--license_number"></a>
### Nested Schema for `license_number`

Optional:

- **label** (String) Field label.
- **privacy_screen** (Boolean) Privacy screen flag.
- **required** (Boolean) Required flag.
- **value** (String) Field value.

Read-Only:

- **type** (String) Field type.
