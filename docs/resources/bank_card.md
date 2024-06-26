# secretsmanager_bank_card Resource

Use this resource to access secrets of type `bankCard` stored in Keeper Vault

## Schema

### Optional

- **address_ref** (Block List, Max: 1) AddressRef field data. (see [below for nested schema](#nestedblock--address_ref))
- **cardholder_name** (Block List, Max: 1) Text field data. (see [below for nested schema](#nestedblock--cardholder_name))
- **file_ref** (Block List, Max: 1) FileRef field data. (see [below for nested schema](#nestedblock--file_ref))
- **folder_uid** (String) The folder UID where the secret is stored. The parent shared folder must be non empty.
- **id** (String) The ID of this resource.
- **notes** (String) The secret notes.
- **payment_card** (Block List, Max: 1) Payment card field data. (see [below for nested schema](#nestedblock--payment_card))
- **pin_code** (Block List, Max: 1) PinCode field data. (see [below for nested schema](#nestedblock--pin_code))
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

<a id="nestedblock--cardholder_name"></a>
### Nested Schema for `cardholder_name`

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

<a id="nestedblock--payment_card"></a>
### Nested Schema for `payment_card`

Optional:

- **label** (String) Field label.
- **privacy_screen** (Boolean) Privacy screen flag.
- **required** (Boolean) Required flag.
- **value** (Block List) Field value. (see [below for nested schema](#nestedblock--payment_card--value))

Read-Only:

- **type** (String) Field type.

<a id="nestedblock--payment_card--value"></a>
### Nested Schema for `payment_card.value`

Optional:

- **card_expiration_date** (String) Card expiration date.
- **card_number** (String) Card number.
- **card_security_code** (String) Card security code.

<a id="nestedblock--pin_code"></a>
### Nested Schema for `pin_code`

Optional:

- **label** (String) Field label.
- **privacy_screen** (Boolean) Privacy screen flag.
- **required** (Boolean) Required flag.
- **value** (String) Field value.

Read-Only:

- **type** (String) Field type.
