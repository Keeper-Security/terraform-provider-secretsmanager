# secretsmanager_contact Resource

Use this resource to access secrets of type `contact` stored in Keeper Vault

## Schema

### Optional

- **address_ref** (Block List, Max: 1) AddressRef field data. (see [below for nested schema](#nestedblock--address_ref))
- **company** (Block List, Max: 1) Text field data. (see [below for nested schema](#nestedblock--company))
- **email** (Block List, Max: 1) Email field data. (see [below for nested schema](#nestedblock--email))
- **file_ref** (Block List, Max: 1) FileRef field data. (see [below for nested schema](#nestedblock--file_ref))
- **folder_uid** (String) The folder UID where the secret is stored. The shared folder must be non empty.
- **id** (String) The ID of this resource.
- **name** (Block List, Max: 1) Name field data. (see [below for nested schema](#nestedblock--name))
- **notes** (String) The secret notes.
- **phone** (Block List, Max: 1) Phone field data. (see [below for nested schema](#nestedblock--phone))
- **title** (String) The secret title.
- **uid** (String) The UID of the new secret (using RFC4648 URL and Filename Safe Alphabet).

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

<a id="nestedblock--company"></a>
### Nested Schema for `company`

Optional:

- **label** (String) Field label.
- **privacy_screen** (Boolean) Privacy screen flag.
- **required** (Boolean) Required flag.
- **value** (String) Field value.

Read-Only:

- **type** (String) Field type.

<a id="nestedblock--email"></a>
### Nested Schema for `email`

Optional:

- **label** (String) Field label.
- **privacy_screen** (Boolean) Privacy screen flag.
- **required** (Boolean) Required flag.
- **value** (String) Field value.

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

<a id="nestedblock--phone"></a>
### Nested Schema for `phone`

Optional:

- **label** (String) Field label.
- **privacy_screen** (Boolean) Privacy screen flag.
- **required** (Boolean) Required flag.
- **value** (Block List) Field value. (see [below for nested schema](#nestedblock--phone--value))

Read-Only:

- **type** (String) Field type.

<a id="nestedblock--phone--value"></a>
### Nested Schema for `phone.value`

Optional:

- **ext** (String) Extension number - ex. 9987
- **number** (String) Phone number - ex. 510-222-5555
- **region** (String) Region code - ex. US
- **type** (String) Phone number type - ex. Home, Work, or Mobile
