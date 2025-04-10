# secretsmanager_folder Resource

Use this resource to manage folders in Keeper Vault

## Schema

### Required

- **name** (String) The folder name.
- **parent_uid** (String) The parent folder UID where the folder is created.

### Optional

- **force_delete** (Boolean) Force deletion of non empty folders.
- **id** (String) The ID of this resource.

### Read-Only

- **uid** (String) The folder UID (using RFC4648 URL and Filename Safe Alphabet).
