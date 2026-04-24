# secretsmanager_server_credentials Resource

Use this resource to create and manage secrets of type `serverCredentials` in Keeper Vault

## Example Usage

```terraform
resource "secretsmanager_server_credentials" "my_server_credentials" {
  folder_uid = "<folder UID>"
  title      = "My Title"
  notes      = "My Notes"

  host {
    label          = "My Host"
    required       = true
    privacy_screen = true
    value {
      host_name = "127.0.0.1"
      port      = "22"
    }
  }

  login {
    label          = "My Login"
    required       = true
    privacy_screen = true
    value          = "MyLogin"
  }

  password {
    label              = "My Pass"
    required           = true
    privacy_screen     = true
    enforce_generation = true
    generate           = "yes"
    complexity {
      length    = 20
      caps      = 5
      lowercase = 5
      digits    = 5
      special   = 5
    }
  }
}
```

## Schema

### Optional

- **file_ref** (Block List, Max: 1) FileRef field data. (see [below for nested schema](#nestedblock--file_ref))
- **folder_uid** (String) The folder UID where the secret is stored. The parent shared folder must be non empty.
- **host** (Block List, Max: 1) Host field data. (see [below for nested schema](#nestedblock--host))
- **id** (String) The ID of this resource.
- **login** (Block List, Max: 1) Login field data. (see [below for nested schema](#nestedblock--login))
- **notes** (String) The secret notes.
- **password** (Block List, Max: 1) Password field data. (see [below for nested schema](#nestedblock--password))
- **title** (String) The secret title.
- **uid** (String) The UID of the new secret (using RFC4648 URL and Filename Safe Alphabet).

- **custom** (Block List) User-defined custom fields. (see [below for nested schema](#nestedblock--custom))

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

<a id="nestedblock--custom"></a>
### Nested Schema for `custom`

Required:

- **label** (String) Display name for the field in Keeper UI.
- **type** (String) Keeper field type. Input is case-insensitive — any casing is accepted and normalized (e.g., `paymentcard` → `paymentCard`). Unknown types are rejected at plan time. Common values: `text`, `secret`, `url`, `email`, `phone`, `date`, `birthDate`, `expirationDate`, `name`, `address`, `paymentCard`, `bankAccount`, `host`, `keyPair`, `securityQuestion`, `checkbox`, `multiline`.

Optional:

- **privacy_screen** (Boolean) Whether this field is hidden behind a privacy screen in the Keeper UI.
- **required** (Boolean) Whether this field is required.
- **value** (String, Sensitive) Field value. Plain string for simple types. Use `jsonencode({...})` for structured types or `jsonencode([{...},{...}])` for multiple entries in one field. Format constraints: `checkbox` requires `"true"` or `"false"`; `date`, `birthDate`, and `expirationDate` require YYYY-MM-DD; `paymentCard` `jsonencode` keys use camelCase (`cardNumber`, `cardExpirationDate`, `cardSecurityCode`).
