# secretsmanager_contact Data Source

Use this data source to read secrets of type `contact` stored in Keeper Vault

## Example Usage

```terraform
data "secretsmanager_contact" "contact" {
  path = "<record UID>"
}
```

## Argument Reference

* `path` - (Required) The UID of existing record in Keeper Vault.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `type` - The type of the record.
* `title` - Record title.
* `notes` - Record notes.
* `name` - A list containing name information:
  - `first` - First name
  - `middle` - Middle name
  - `last` - Last name
* `company` - Company name.
* `email` - Contact's e-mail.
* `phone` - A list containing phone information:
  - `region` - 2 letter country code (ISO 3166-1 alpha-2)
  - `number` - Phone number
  - `ext` - Phone extension
  - `type` - Phone type: Mobile, Home or Work
* `address_ref` - A list containing address information:
  - `uid` - The address reference record UID
  - `street1` - Street line 1
  - `street2` - Street line 2
  - `city` - City
  - `state` - State
  - `zip` - Zip
  - `country` - Country
* `file_ref` - A list containing file reference information:
  - `uid` - File UID
  - `title` - File title
  - `name` - File name
  - `type` - File content type
  - `size` - File size
  - `last_modified` - File last modification timestamp
  - `content_base64` - File content base64 encoded
