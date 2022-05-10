# secretsmanager_database_credentials Resource

Use this resource to access secrets of type `databaseCredentials` stored in Keeper Vault

## Schema

### Optional

- **db_type** (Block List, Max: 1) Text field data. (see [below for nested schema](#nestedblock--db_type))
- **file_ref** (Block List, Max: 1) FileRef field data. (see [below for nested schema](#nestedblock--file_ref))
- **folder_uid** (String) The folder UID where the secret is stored. The shared folder must be non empty.
- **host** (Block List, Max: 1) Host field data. (see [below for nested schema](#nestedblock--host))
- **id** (String) The ID of this resource.
- **login** (Block List, Max: 1) Login field data. (see [below for nested schema](#nestedblock--login))
- **notes** (String) The secret notes.
- **password** (Block List, Max: 1) Password field data. (see [below for nested schema](#nestedblock--password))
- **title** (String) The secret title.
- **uid** (String) The UID of the new secret (using RFC4648 URL and Filename Safe Alphabet).

### Read-Only

- **type** (String) The secret type.

<a id="nestedblock--db_type"></a>
### Nested Schema for `db_type`

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

<a id="nestedblock--host"></a>
### Nested Schema for `host`

Optional:

- **label** (String) Field label.
- **privacy_screen** (Boolean) Privacy screen flag.
- **required** (Boolean) Required flag.
- **value** (Block List) Field value. (see [below for nested schema](#nestedblock--host--value))

Read-Only:

- **type** (String) Field type.

<a id="nestedblock--host--value"></a>
### Nested Schema for `host.value`

Optional:

- **host_name** (String) Hostname.
- **port** (String) Port.

<a id="nestedblock--login"></a>
### Nested Schema for `login`

Optional:

- **label** (String) Field label.
- **privacy_screen** (Boolean) Privacy screen flag.
- **required** (Boolean) Required flag.
- **value** (String) Field value.

Read-Only:

- **type** (String) Field type.

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
