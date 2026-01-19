# Terraform Provider for Rizhiyi

This is the Terraform provider for Rizhiyi (日志易). It allows you to manage Rizhiyi resources such as accounts, roles, indexes, dashboards, alerts, and parser rules via Terraform.

## Requirements

*   [Terraform](https://www.terraform.io/downloads.html) v0.13+
*   [Go](https://golang.org/doc/install) 1.20.4+ (to build the provider plugin)

## Building The Provider

1.  Clone the repository.
2.  Build the provider using Go:

    ```bash
    go build -o terraform-provider-rizhiyi
    ```

## Installation

We will be using the implicit local mirror method to install our custom provider.

### Linux System

Create the directory structure:

```bash
mkdir -p ~/.terraform.d/plugins/terraform-rizhiyi.com/rizhiyiprovider/rizhiyi/1.0.0/linux_amd64
```

Copy the binary:

```bash
cp terraform-provider-rizhiyi ~/.terraform.d/plugins/terraform-rizhiyi.com/rizhiyiprovider/rizhiyi/1.0.0/linux_amd64/
```

### Windows System

Create the directory structure:

```cmd
mkdir %APPDATA%\terraform.d\plugins\terraform-rizhiyi.com\rizhiyiprovider\rizhiyi\1.0.0\windows_amd64
```

Copy the binary to the created folder.

### CLI Configuration (`.terraformrc`)

Create or update `$HOME/.terraformrc` (or `%APPDATA%\terraform.rc` on Windows) with the following content to enable the local plugin:

```hcl
plugin_cache_dir   = "$HOME/.terraform.d/plugin-cache"
disable_checkpoint = true
```

## Provider Configuration

The Rizhiyi provider can be configured via Terraform configuration or environment variables.

```hcl
provider "rizhiyi" {
  host  = "192.168.1.224:8090"
  token = "cml6aGl5aToxMjM0NTY=" # Base64 encoded admin:password
}
```

Or using environment variables:

*   `RIZHIYI_HOST`: The endpoint of your Rizhiyi resource server.
*   `RIZHIYI_TOKEN`: The HTTP Basic Authentication token (Base64 encoded `username:password`).

## Supported Resources

*   `rizhiyi_account`: Manage user accounts.
*   `rizhiyi_role`: Manage user roles.
*   `rizhiyi_index`: Manage log indexes.
*   `rizhiyi_dashboard`: Manage dashboards.
*   `rizhiyi_alert`: Manage alerts.
*   `rizhiyi_parser_rule`: Manage parser rules.

## Examples

Check the `examples/` directory for usage examples.

```bash
cd examples
terraform init
terraform plan
terraform apply
```

**NOTE:** When developing and testing local provider builds, if terraform version `>= 0.13 +` you would have to replace the provider binaries in the `.terraform` folder with your local build. [Follow these guidelines](https://github.com/hashicorp/terraform/blob/master/website/upgrade-guides/0-13.html.markdown)
