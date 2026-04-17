# secretsmanager_pam_machine Resource

Use this resource to create and manage PAM Machine records in Keeper Vault. Supports SSH key generation with optional passphrase encryption for privileged access management.

## Example Usage

### Basic PAM Machine

```terraform
resource "secretsmanager_pam_machine" "basic" {
  folder_uid = "<folder UID>"
  title      = "Production Server"

  pam_hostname {
    value {
      hostname = "10.0.1.50"
      port     = "22"
    }
  }

  login {
    value = "root"
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

### PAM Machine with Generated SSH Key

```terraform
resource "secretsmanager_pam_machine" "ssh_machine" {
  folder_uid = "<folder UID>"
  title      = "SSH-Managed Server"

  pam_hostname {
    value {
      hostname = "10.0.1.50"
      port     = "22"
    }
  }

  login {
    value = "deploy"
  }

  private_key_passphrase {
    generate = "yes"
    complexity { length = 32 }
  }

  private_pem_key {
    generate = "yes"
    key_type = "ssh-ed25519"
  }

  operating_system {
    value = ["Linux"]
  }
}
```

## Schema

### Optional

- **file_ref** (Block List, Max: 1) FileRef field data.
- **folder_uid** (String) The folder UID where the secret is stored.
- **instance_id** (Block List, Max: 1) Text field data. Label: "Instance Id".
- **instance_name** (Block List, Max: 1) Text field data. Label: "Instance Name".
- **login** (Block List, Max: 1) Login field data.
- **notes** (String) The secret notes.
- **operating_system** (Block List, Max: 1) Text field data. Label: "Operating System".
- **pam_hostname** (Block List, Max: 1) PAM Hostname field data.
- **password** (Block List, Max: 1) Password field data.
- **private_key_passphrase** (Block List, Max: 1) Private key passphrase. Stored as a custom field labeled "Private Key Passphrase". When used with key generation, the passphrase encrypts the generated private key. (see [below for nested schema](#nestedblock--private_key_passphrase))
- **private_pem_key** (Block List, Max: 1) Private PEM Key field data. Stored as a secret field labeled "Private PEM Key". Supports SSH key generation. (see [below for nested schema](#nestedblock--private_pem_key))
- **provider_group** (Block List, Max: 1) Text field data. Label: "Provider Group".
- **provider_region** (Block List, Max: 1) Text field data. Label: "Provider Region".
- **rotation_scripts** (Block List) Script field data. Label: "Rotation Scripts".
- **ssl_verification** (Block List, Max: 1) Checkbox field data. Label: "SSL Verification".
- **title** (String) The secret title.
- **totp** (Block List, Max: 1) One-time code field data.
- **uid** (String) The UID of the new secret (using RFC4648 URL and Filename Safe Alphabet).

- **custom** (Block List) User-defined custom fields. (see [below for nested schema](#nestedblock--custom))

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

<a id="nestedblock--custom"></a>
### Nested Schema for `custom`

Required:

- **label** (String) Display name for the field in Keeper UI.
- **type** (String) Keeper field type. Input is case-insensitive — any casing is accepted and normalized (e.g., `paymentcard` → `paymentCard`). Unknown types are rejected at plan time. Common values: `text`, `secret`, `url`, `email`, `phone`, `date`, `birthDate`, `expirationDate`, `name`, `address`, `paymentCard`, `bankAccount`, `host`, `keyPair`, `securityQuestion`, `checkbox`, `multiline`.

Optional:

- **privacy_screen** (Boolean) Whether this field is hidden behind a privacy screen in the Keeper UI.
- **required** (Boolean) Whether this field is required.
- **value** (String, Sensitive) Field value. Plain string for simple types. Use `jsonencode({...})` for structured types or `jsonencode([{...},{...}])` for multiple entries in one field. Format constraints: `checkbox` requires `"true"` or `"false"`; `date`, `birthDate`, and `expirationDate` require YYYY-MM-DD; `paymentCard` `jsonencode` keys use camelCase (`cardNumber`, `cardExpirationDate`, `cardSecurityCode`).
