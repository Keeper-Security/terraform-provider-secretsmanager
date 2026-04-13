# secretsmanager_folder Resource

Use this resource to manage folders in Keeper Vault

## Schema

### Required

- **name** (String) The folder name.
- **parent_uid** (String) The parent folder UID where the folder is created.

### Optional

- **force_delete** (Boolean) Force deletion of non empty folders.
- **id** (String) The ID of this resource.

- **custom** (Block List) User-defined custom fields. (see [below for nested schema](#nestedblock--custom))

### Read-Only

- **uid** (String) The folder UID (using RFC4648 URL and Filename Safe Alphabet).

<a id="nestedblock--custom"></a>
### Nested Schema for `custom`

Required:

- **label** (String) Display name for the field in Keeper UI.
- **type** (String) Keeper field type. Common values: `text`, `secret`, `url`, `email`, `phone`, `date`, `birthDate`, `expirationDate`, `name`, `address`, `paymentCard`, `bankAccount`, `host`, `keyPair`, `securityQuestion`, `checkbox`, `multiline`.

Optional:

- **privacy_screen** (Boolean) Whether this field is hidden behind a privacy screen in the Keeper UI.
- **required** (Boolean) Whether this field is required.
- **value** (String, Sensitive) Field value. Plain string for simple types. Use `jsonencode({...})` for structured types or `jsonencode([{...},{...}])` for multiple entries in one field. Format constraints: `checkbox` requires `"true"` or `"false"`; `date`, `birthDate`, and `expirationDate` require YYYY-MM-DD; `paymentCard` `jsonencode` keys use camelCase (`cardNumber`, `cardExpirationDate`, `cardSecurityCode`).
