<h1 align="center">Keeper Secrets Management For Terraform</h1>
<p align="center">
  <a href="https://docs.keeper.io/secrets-manager/secrets-manager/integrations/terraform">View docs</a>
</p>
<br/>

Keeper Secrets Manager provides your DevOps, IT Security and software development teams with a fully cloud-based, zero-knowledge platform for managing all of your infrastructure secrets such as API keys, database passwords, access keys, certificates and any type of confidential data. Essential tool for every engineer who wants to securely provision passwords and keys throughout entire development stack with just a few lines of code.

## Setup Secrets Manager

In order to set up Secrets Manager on a Keeper Enterprise Account follow the [Quick Start Guide](https://docs.keeper.io/secrets-manager/secrets-manager/quick-start-guide).

### Create Secrets Manager application
- Using Keeper **Commander** CLI
```bash
My Vault> sm app create [NAME]
My Vault> sm share add --app [NAME] --secret [UID] --editable
My Vault> sm client add --app [NAME] --unlock-ip --count 1
```
- Using Keeper **Secrets Manager** CLI and token generated while creating client (_use_ `sm client add` command above) generate local configuration
```bash
$ ksm profile init --token [TOKEN]
```

- Find record UID of a shared secret you want to use
```bash
$ ksm secret list
$ ksm secret get -u [UID]
```

### Plugin configuration
- Keeper credential could be generated with `ksm profile init` command, read from file, or sourced from the `KEEPER_CREDENTIAL` environment variable.  
Generate `credential` using Commander CLI
```
sm client add --app <APP_NAME> --unlock-ip --config-init=b64
```
`main.tf`
```
terraform {
  required_providers {
    # add keeper secrets manager plugin
    secretsmanager = {
      source  = "keeper-security/secretsmanager"
      version = ">= 1.1.7"
    }
  }
}

# Configure plugin
provider "secretsmanager" {
  credential = file("~/.keeper/credential")
}
```
- Data source usage - see working [examples](./examples) in this repo.

## Support

If you need help, send an e-mail to [sm@keepersecurity.com](mailto:sm@keepersecurity.com)

## Development

### Building

Get the source code:

```bash
git clone https://github.com/keeper-security/terraform-provider-secretsmanager
```

Build it using:

```bash
go build
```

### Testing

To run the [acceptance tests](https://www.terraform.io/docs/extend/testing/acceptance-tests/index.html), the following environment variables need to be set up.

* `KEEPER_CREDENTIAL` - Keeper Secrets Manager Credentials.

The acceptance tests expect to find certain records shared to your application - use the script below to create and populate shared folder named `tf_acc_test_dir` with the required records (_use_ [Keeper Commander CLI](https://docs.keeper.io/secrets-manager/commander-cli))

_Note:_ If you get **throttled** simply re-run the same command again (_and ignore any_ `'...already exists'` _messages on consecutive runs_)

`keeper tf_acc_test.cmd --batch-mode`

Contents of `tf_acc_test.cmd`:
```
@mkdir -sf -a /tf_acc_test_dir
@cd /tf_acc_test_dir
@add title=tf_acc_test_field notes=tf_acc_test_field type=login fields.login=tf_acc_test_field
@add title=tf_acc_test_login notes=tf_acc_test_login type=login
@add title=tf_acc_test_bank_account notes=tf_acc_test_bank_account type=bankAccount fields.bankAccount.accountNumber=1234
@add title=tf_acc_test_address notes=tf_acc_test_address type=address
@add title=tf_acc_test_bank_card notes=tf_acc_test_bank_card type=bankCard
@add title=tf_acc_test_birth_certificate notes=tf_acc_test_birth_certificate type=birthCertificate
@add title=tf_acc_test_contact notes=tf_acc_test_contact type=contact fields.name.first=John fields.name.last=Doe
@add title=tf_acc_test_driver_license notes=tf_acc_test_driver_license type=driverLicense
@add title=tf_acc_test_encrypted_notes notes=tf_acc_test_encrypted_notes type=encryptedNotes
@add title=tf_acc_test_file notes=tf_acc_test_file type=file
@add title=tf_acc_test_health_insurance notes=tf_acc_test_health_insurance type=healthInsurance
@add title=tf_acc_test_membership notes=tf_acc_test_membership type=membership
@add title=tf_acc_test_passport notes=tf_acc_test_passport type=passport
@add title=tf_acc_test_photo notes=tf_acc_test_photo type=photo
@add title=tf_acc_test_server_credentials notes=tf_acc_test_server_credentials type=serverCredentials
@add title=tf_acc_test_software_license notes=tf_acc_test_software_license type=softwareLicense
@add title=tf_acc_test_ssn_card notes=tf_acc_test_ssn_card type=ssnCard
@add title=tf_acc_test_ssh_keys notes=tf_acc_test_ssh_keys type=sshKeys
@add title=tf_acc_test_database_credentials notes=tf_acc_test_database_credentials type=databaseCredentials
```

With the environment variables properly set up, run:

```bash
export TF_ACC=1 ; go test ./...
```

or set all required environment variables and run tests with a single command line
```bash
export TF_ACC=1 ; export KEEPER_CREDENTIAL=<XXX> ; go test ./...
```
------
# Terraform Provider

The Keeper Secrets Manager Terraform Provider lets you manage your secrets using Terraform.
It is officially supported and actively maintained by Keeper Security.

## Usage
### Terraform v0.13 or above ([Terraform Registry](https://registry.terraform.io/))
```hcl
terraform {
  required_providers {
    secretsmanager = {
      source  = "keeper-security/secretsmanager"
      version = ">= 1.0.0"
    }
  }
}

provider "secretsmanager" {
  credential = "<CREDENTIAL>"
  # credential = file("~/.keeper/credential")
}

data "secretsmanager_database_credentials" "my_db_creds" {
  path  = "<UID>"
}

output "db_type" {
  value = data.secretsmanager_database_credentials.my_db_creds.db_type
}

output "login" {
  value = data.secretsmanager_database_credentials.my_db_creds.login
}
```

### Terraform v0.13 and above ([GitHub](https://github.com/keeper-security/terraform-provider-secretsmanager/) manual install)

Download archive with the [latest release](https://github.com/keeper-security/terraform-provider-secretsmanager/releases/latest) for your platform and copy it to the corresponding plugin folder (_Linux and MacOS:_ `~/.terraform.d/plugins/github.com/keeper-security/secretsmanager` _Windows:_ `%APPDATA%/terraform.d/plugins/github.com/keeper-security/secretsmanager`)  
Use the same config from above just remember to initialize `source` with the full URL `source  = "github.com/keeper-security/secretsmanager"`

MacOS:
```bash
mkdir -p ~/.terraform.d/plugins/github.com/keeper-security/secretsmanager && \
cd ~/.terraform.d/plugins/github.com/keeper-security/secretsmanager && \
curl -SfLOJ https://github.com/keeper-security/terraform-provider-secretsmanager/releases/latest/download/terraform-provider-secretsmanager_1.0.0_darwin_amd64.zip
```
Windows:
```bash
SETLOCAL EnableExtensions && ^
mkdir %APPDATA%\.terraform.d\plugins\github.com\keeper-security\secretsmanager && ^
cd %APPDATA%\.terraform.d\plugins\github.com\keeper-security\secretsmanager && ^
curl -SfLOJ https://github.com/keeper-security/terraform-provider-secretsmanager/releases/latest/download/terraform-provider-secretsmanager_1.0.0_windows_amd64.zip
```
Have a look at some working [examples](./examples) in this repo.

### Terraform v0.12 and below
Manually install the Keeper Secrets Manager provider by downloading the corresponding archive for your platform then extract the executable and move it to `~/.terraform/plugins` or `%APPDATA%\terraform.d\plugins` on Windows.

Afterwards you can run the following example with Terraform.
```hcl
terraform {
  required_providers {
    secretsmanager = {
      version = ">= 1.0.0"
    }
  }
}

provider "secretsmanager" {
  credential = "<CREDENTIAL>"
  # credential = file("~/.keeper/credential")
}

data "secretsmanager_database_credentials" "my_db_creds" {
  path  = "<UID>"
}

output "db_type" {
  value = data.secretsmanager_database_credentials.my_db_creds.db_type
}

output "login" {
  value = data.secretsmanager_database_credentials.my_db_creds.login
}
```
Have a look at some working [examples](./examples) in this repo.
