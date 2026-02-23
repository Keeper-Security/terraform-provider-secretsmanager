# secretsmanager_pam_user Resource

Use this resource to create and manage PAM User records in Keeper Vault. Supports SSH key generation with optional passphrase encryption for privileged access management.

## Example Usage

### Basic PAM User

```terraform
resource "secretsmanager_pam_user" "basic" {
  folder_uid = "<folder UID>"
  title      = "PAM Admin User"

  login {
    value = "admin"
  }

  password {
    generate = "yes"
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

### PAM User with Generated SSH Key

```terraform
resource "secretsmanager_pam_user" "ssh_user" {
  folder_uid = "<folder UID>"
  title      = "PAM SSH User"

  login {
    value = "deploy"
  }

  password {
    generate = "yes"
    complexity { length = 20 }
  }

  private_key_passphrase {
    generate = "yes"
    complexity { length = 32 }
  }

  private_pem_key {
    generate = "yes"
    key_type = "ssh-ed25519"
  }
}
```

### PAM User with RSA Key

```terraform
resource "secretsmanager_pam_user" "rsa_user" {
  folder_uid = "<folder UID>"
  title      = "PAM RSA User"

  login {
    value = "service-account"
  }

  private_pem_key {
    generate = "yes"
    key_type = "ssh-rsa"
    key_bits = 4096
  }
}
```

## Schema

### Optional

- **connect_database** (Block List, Max: 1) Text field data. Label: "Connect Database".
- **distinguished_name** (Block List, Max: 1) Text field data. Label: "Distinguished Name".
- **file_ref** (Block List, Max: 1) FileRef field data.
- **folder_uid** (String) The folder UID where the secret is stored.
- **login** (Block List, Max: 1) Login field data.
- **managed** (Block List, Max: 1) Checkbox field data. Label: "Managed".
- **notes** (String) The secret notes.
- **password** (Block List, Max: 1) Password field data.
- **private_key_passphrase** (Block List, Max: 1) Private key passphrase. Stored as a custom field labeled "Private Key Passphrase". When used with key generation, the passphrase encrypts the generated private key. (see [below for nested schema](#nestedblock--private_key_passphrase))
- **private_pem_key** (Block List, Max: 1) Private PEM Key field data. Stored as a secret field labeled "Private PEM Key". Supports SSH key generation. (see [below for nested schema](#nestedblock--private_pem_key))
- **rotation_scripts** (Block List) Script field data. Label: "Rotation Scripts".
- **title** (String) The secret title.
- **totp** (Block List, Max: 1) One-time code field data.
- **uid** (String) The UID of the new secret (using RFC4648 URL and Filename Safe Alphabet).

### Read-Only

- **type** (String) The secret type.

<a id="nestedblock--private_pem_key"></a>
### Nested Schema for `private_pem_key`

Optional:

- **generate** (String) Flag to force SSH key generation (when set to 'yes' or 'true'). When set, `value` and `public_key` are computed automatically.
- **key_type** (String) SSH key type. One of: `ssh-ed25519` (default), `ssh-rsa`, `ecdsa-sha2-nistp256`, `ecdsa-sha2-nistp384`, `ecdsa-sha2-nistp521`.
- **key_bits** (Number) Key size in bits. Only used for `ssh-rsa`. Valid: 2048, 3072, 4096. Default: `4096`.
- **label** (String) Field label.
- **privacy_screen** (Boolean) Privacy screen flag.
- **required** (Boolean) Required flag.
- **value** (String, Sensitive) Private key in PEM format. Computed when `generate` is set.

Read-Only:

- **public_key** (String) Public key in OpenSSH format. Computed when `generate` is set.
- **type** (String) Field type.

<a id="nestedblock--private_key_passphrase"></a>
### Nested Schema for `private_key_passphrase`

Stored as a custom field on the record with type `secret` and label "Private Key Passphrase". Compatible with PAM/kdnrm managed records. When both passphrase and key generation are enabled, the passphrase encrypts the generated private key.

Optional:

- **complexity** (Block List, Max: 1) Passphrase complexity. (see [below for nested schema](#nestedblock--private_key_passphrase--complexity))
- **generate** (String) Flag to force passphrase generation (when set to 'yes' or 'true').
- **value** (String, Sensitive) Passphrase value. Computed when `generate` is set.

Read-Only:

- **type** (String) Field type.

<a id="nestedblock--private_key_passphrase--complexity"></a>
### Nested Schema for `private_key_passphrase.complexity`

Optional:

- **caps** (Number) Number of uppercase characters.
- **digits** (Number) Number of digits.
- **length** (Number) Passphrase length.
- **lowercase** (Number) Number of lowercase characters.
- **special** (Number) Number of special characters.
