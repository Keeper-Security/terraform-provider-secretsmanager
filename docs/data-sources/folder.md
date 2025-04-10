# secretsmanager_folder Data Source

Use this data source to access individual folders in Keeper Vault

## Example Usage

```terraform
data "secretsmanager_folder" "folder" {
  name = "<Folder Name>"
}
```

## Schema

### Optional

- **id** (String) The ID of this resource.
- **name** (String) The folder name.
- **uid** (String) The folder uid.

### Read-Only

- **parent_uid** (String) The parent folder uid.
- **shared** (Boolean) Shared folder flag.


