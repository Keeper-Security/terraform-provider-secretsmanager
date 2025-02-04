# secretsmanager_folders Data Source

Use this data source to list all folders shared to a KSM Application

## Example Usage

```terraform
data "secretsmanager_folders" "folders" { }
```

## Schema

### Optional

- **id** (String) The ID of this resource.

### Read-Only

- **folders** (List of Object) List of all folders shared to the KSM Application. (see [below for nested schema](#nestedatt--folders))

<a id="nestedatt--folders"></a>
### Nested Schema for `folders`

Read-Only:

- **name** (String)
- **parent_uid** (String)
- **shared** (Boolean)
- **uid** (String)
