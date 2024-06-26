# secretsmanager_health_insurance Resource

Use this resource to access secrets of type `healthInsurance` stored in Keeper Vault

## Schema

### Optional

- **account_number** (Block List, Max: 1) Account number field data. (see [below for nested schema](#nestedblock--account_number))
- **file_ref** (Block List, Max: 1) FileRef field data. (see [below for nested schema](#nestedblock--file_ref))
- **folder_uid** (String) The folder UID where the secret is stored. The parent shared folder must be non empty.
- **id** (String) The ID of this resource.
- **login** (Block List, Max: 1) Login field data. (see [below for nested schema](#nestedblock--login))
- **name** (Block List, Max: 1) Name field data. (see [below for nested schema](#nestedblock--name))
- **notes** (String) The secret notes.
- **password** (Block List, Max: 1) Password field data. (see [below for nested schema](#nestedblock--password))
- **title** (String) The secret title.
- **uid** (String) The UID of the new secret (using RFC4648 URL and Filename Safe Alphabet).
- **url** (Block List, Max: 1) URL field data. (see [below for nested schema](#nestedblock--url))

### Read-Only

- **type** (String) The secret type.

<a id="nestedblock--account_number"></a>
### Nested Schema for `account_number`

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

<a id="nestedblock--login"></a>
### Nested Schema for `login`

Optional:

- **label** (String) Field label.
- **privacy_screen** (Boolean) Privacy screen flag.
- **required** (Boolean) Required flag.
- **value** (String) Field value.

Read-Only:

- **type** (String) Field type.

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

<a id="nestedblock--password"></a>
### Nested Schema for `password`

Optional:

- **complexity** (Block List, Max: 1) Password complexity. (see [below for nested schema](#nestedblock--password--complexity))
- **enforce_generation** (Boolean) Enforce generation flag.
- **generate** (String) Flag to force password generation (when set to 'yes' or 'true').
- **label** (String) Field label.
- **privacy_screen** (Boolean) Privacy screen flag.
- **required** (Boolean) Required flag.
- **value** (String, Sensitive) Field value.

Read-Only:

- **type** (String) Field type.

<a id="nestedblock--password--complexity"></a>
### Nested Schema for `password.complexity`

Optional:

- **caps** (Number) Number of uppercase characters.
- **digits** (Number) Number of digits.
- **length** (Number) Password length.
- **lowercase** (Number) Number of lowercase characters.
- **special** (Number) Number of special characters.

<a id="nestedblock--url"></a>
### Nested Schema for `url`

Optional:

- **label** (String) Field label.
- **privacy_screen** (Boolean) Privacy screen flag.
- **required** (Boolean) Required flag.
- **value** (String) Field value.

Read-Only:

- **type** (String) Field type.
