# secretsmanager_ssh_keys Resource

Use this resource to create and manage secrets of type `sshKeys` in Keeper Vault. Supports automatic SSH key pair generation with optional passphrase encryption.

## Example Usage

### Manual SSH Keys

```terraform
resource "secretsmanager_ssh_keys" "my_ssh_keys" {
  folder_uid = "<folder UID>"
  title      = "My Title"
  notes      = "My Notes"

  login {
    label          = "My Login"
    required       = true
    privacy_screen = true
    value          = "MyLogin"
  }

  passphrase {
    label          = "My Pass"
    required       = true
    privacy_screen = true
    value          = "<SSH PASSPHRASE>"
  }

  host {
    label          = "My Host"
    required       = true
    privacy_screen = true
    value {
      host_name = "127.0.0.1"
      port      = "22"
    }
  }

  key_pair {
    label          = "My Keys"
    required       = true
    privacy_screen = true
    value {
      public_key  = "<PUBLIC KEY>"
      private_key = "<PRIVATE KEY>"
    }
  }
}
```

### Generated SSH Keys (ED25519)

```terraform
resource "secretsmanager_ssh_keys" "generated_ed25519" {
  folder_uid = "<folder UID>"
  title      = "Generated ED25519 Key"

  key_pair {
    generate = "yes"
    key_type = "ssh-ed25519"
  }
}
```

### Generated SSH Keys with Passphrase

When both `passphrase` and `key_pair` use `generate = "yes"`, the generated passphrase automatically encrypts the private key.

```terraform
resource "secretsmanager_ssh_keys" "generated_with_passphrase" {
  folder_uid = "<folder UID>"
  title      = "Generated RSA Key with Passphrase"

  passphrase {
    generate = "yes"
    complexity {
      length = 32
    }
  }

  key_pair {
    generate = "yes"
    key_type = "ssh-rsa"
    key_bits = 4096
  }
}
```

## Schema

### Optional

- **file_ref** (Block List, Max: 1) FileRef field data. (see [below for nested schema](#nestedblock--file_ref))
- **folder_uid** (String) The folder UID where the secret is stored. The parent shared folder must be non empty.
- **host** (Block List, Max: 1) Host field data. (see [below for nested schema](#nestedblock--host))
- **id** (String) The ID of this resource.
- **key_pair** (Block List, Max: 1) Key pair field data. (see [below for nested schema](#nestedblock--key_pair))
- **login** (Block List, Max: 1) Login field data. (see [below for nested schema](#nestedblock--login))
- **notes** (String) The secret notes.
- **passphrase** (Block List, Max: 1) Password field data. Used as the SSH key passphrase. When both passphrase and key pair generation are enabled, the passphrase encrypts the generated private key. (see [below for nested schema](#nestedblock--passphrase))
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

<a id="nestedblock--key_pair"></a>
### Nested Schema for `key_pair`

Optional:

- **generate** (String) Flag to force SSH key pair generation (when set to 'yes' or 'true'). When set, `public_key` and `private_key` are computed automatically.
- **key_type** (String) SSH key type. One of: `ssh-ed25519` (default), `ssh-rsa`, `ecdsa-sha2-nistp256`, `ecdsa-sha2-nistp384`, `ecdsa-sha2-nistp521`.
- **key_bits** (Number) Key size in bits. Only used for `ssh-rsa` key type. Valid values: 2048, 3072, 4096. Default: `4096`.
- **label** (String) Field label.
- **privacy_screen** (Boolean) Privacy screen flag.
- **required** (Boolean) Required flag.
- **value** (Block List) Field value. Computed when `generate` is set. (see [below for nested schema](#nestedblock--key_pair--value))

Read-Only:

- **type** (String) Field type.

<a id="nestedblock--key_pair--value"></a>
### Nested Schema for `key_pair.value`

Optional:

- **private_key** (String, Sensitive) Private key (PEM format). Computed when `generate` is set.
- **public_key** (String) Public key (OpenSSH authorized_keys format). Computed when `generate` is set.

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

When both passphrase and key pair generation are enabled, the passphrase automatically encrypts the generated private key.

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

<a id="nestedblock--custom"></a>
### Nested Schema for `custom`

Required:

- **label** (String) Display name for the field in Keeper UI.
- **type** (String) Keeper field type. Common values: `text`, `secret`, `url`, `email`, `phone`, `date`, `birthDate`, `expirationDate`, `name`, `address`, `paymentCard`, `bankAccount`, `host`, `keyPair`, `securityQuestion`, `checkbox`, `multiline`.

Optional:

- **privacy_screen** (Boolean) Whether this field is hidden behind a privacy screen in the Keeper UI.
- **required** (Boolean) Whether this field is required.
- **value** (String, Sensitive) Field value. Plain string for simple types. Use `jsonencode({...})` for structured types or `jsonencode([{...},{...}])` for multiple entries in one field. Format constraints: `checkbox` requires `"true"` or `"false"`; `date`, `birthDate`, and `expirationDate` require YYYY-MM-DD; `paymentCard` `jsonencode` keys use camelCase (`cardNumber`, `cardExpirationDate`, `cardSecurityCode`).
