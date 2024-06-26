# secretsmanager_ssh_keys Resource

Use this resource to access secrets of type `sshKeys` stored in Keeper Vault

## Schema

### Optional

- **file_ref** (Block List, Max: 1) FileRef field data. (see [below for nested schema](#nestedblock--file_ref))
- **folder_uid** (String) The folder UID where the secret is stored. The parent shared folder must be non empty.
- **host** (Block List, Max: 1) Host field data. (see [below for nested schema](#nestedblock--host))
- **id** (String) The ID of this resource.
- **key_pair** (Block List, Max: 1) Key pair field data. (see [below for nested schema](#nestedblock--key_pair))
- **login** (Block List, Max: 1) Login field data. (see [below for nested schema](#nestedblock--login))
- **notes** (String) The secret notes.
- **passphrase** (Block List, Max: 1) Password field data. (see [below for nested schema](#nestedblock--passphrase))
- **title** (String) The secret title.
- **uid** (String) The UID of the new secret (using RFC4648 URL and Filename Safe Alphabet).

### Read-Only

- **type** (String) The secret type.

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

<a id="nestedblock--key_pair"></a>
### Nested Schema for `key_pair`

Optional:

- **label** (String) Field label.
- **privacy_screen** (Boolean) Privacy screen flag.
- **required** (Boolean) Required flag.
- **value** (Block List) Field value. (see [below for nested schema](#nestedblock--key_pair--value))

Read-Only:

- **type** (String) Field type.

<a id="nestedblock--key_pair--value"></a>
### Nested Schema for `key_pair.value`

Optional:

- **private_key** (String) Private key.
- **public_key** (String) Public key.

<a id="nestedblock--login"></a>
### Nested Schema for `login`

Optional:

- **label** (String) Field label.
- **privacy_screen** (Boolean) Privacy screen flag.
- **required** (Boolean) Required flag.
- **value** (String) Field value.

Read-Only:

- **type** (String) Field type.

<a id="nestedblock--passphrase"></a>
### Nested Schema for `passphrase`

Optional:

- **complexity** (Block List, Max: 1) Password complexity. (see [below for nested schema](#nestedblock--passphrase--complexity))
- **enforce_generation** (Boolean) Enforce generation flag.
- **generate** (String) Flag to force password generation (when set to 'yes' or 'true').
- **label** (String) Field label.
- **privacy_screen** (Boolean) Privacy screen flag.
- **required** (Boolean) Required flag.
- **value** (String, Sensitive) Field value.

Read-Only:

- **type** (String) Field type.

<a id="nestedblock--passphrase--complexity"></a>
### Nested Schema for `passphrase.complexity`

Optional:

- **caps** (Number) Number of uppercase characters.
- **digits** (Number) Number of digits.
- **length** (Number) Password length.
- **lowercase** (Number) Number of lowercase characters.
- **special** (Number) Number of special characters.
