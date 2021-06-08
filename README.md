# get-aws-secret-value

[![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/kitos9112/get-aws-secret-value.svg)](https://hub.docker.com/r/kitos9112/aws_get_secret_value/tags)

Retrieves an AWS secret value as-is after given its secret name. The script can read an AWS Profile name as well as a credentials file in the CLI, otherwise it will automatically load its settings following AWS SDK standards

* Environment Variables
* Shared Credentials file
* Shared Configuration file (if SharedConfig is enabled)
* EC2 Instance Metadata (credentials only)

```bash
> get-aws-secret-value
Help:
  -aws-region string
        AWS Region where to send requests to (default "default")
  -credentials-file string
        Full path to credentials file (default "/root/.aws/credentials")
  -profile string
        AWS Profile to use (default "default")
  -secret-name string
        Secret To Fetch (default "secret")
  -secret-version string
        Version of secret To Fetch (default "version")
```
<!-- TOC -->

- [get-aws-secret-value](#app)
  - [Get it](#get-it)
  - [Use it](#use-it)
    - [Examples](#examples)

<!-- /TOC -->

## Get it

Using go get:

```bash
go get -u github.com/kitos9112/get-aws-secret-value
```

Or [download the binary](https://github.com/kitos9112/get-aws-secret-value/releases/latest) from the releases page.

```bash
# Linux
curl -L https://github.com/kitos9112/get-aws-secret-value/releases/download/0.1.98/get-aws-secret-value_0.1.98_linux_x86_64.tar.gz | tar xz

# OS X
curl -L https://github.com/kitos9112/get-aws-secret-value/releases/download/0.1.98/get-aws-secret-value_0.1.98_osx_x86_64.tar.gz | tar xz

# Windows
curl -LO https://github.com/kitos9112/get-aws-secret-value/releases/download/0.1.98/get-aws-secret-value_0.1.98_windows_x86_64.zip
unzip get-aws-secret-value_0.1.98_windows_x86_64.zip
```

## Use it

```text

get-aws-secret-value [OPTIONS] [COMMAND [ARGS...]]

/bin/sh: get-aws-secret-value: not found
```

### Examples

The simplest example that could easily be integrated into a CICD pipeline:

```shell
> export AWS_PROFILE=myAwsProfile
> export AWS_REGION=eu-west-1
> get-aws-secret-value --secret-name my_secret_name
mySecretValue

```

Or in case you leverage IaC within your favourite public cloud using Terragrunt, you could retrieve the value of an AWS secret previously created and pre-populated by more complext data structures (e.g. JSON)

``` hcl
# terragrunt.hcl
inputs = {
my_secret_var1 = lookup(jsondecode(run_cmd("--terragrunt-quiet", "/usr/local/bin/aws-get-secret-value", "--secret-name", "my_secret", "--aws-region", "eu-west-1")), "secretKey1")
my_secret_var2 = lookup(jsondecode(run_cmd("--terragrunt-quiet", "/usr/local/bin/aws-get-secret-value", "--secret-name", "my_secret", "--aws-region", "eu-west-1")), "secretKey2")
}
```

As you can see, a simple cross-platform binary file could be utilised in many scenarios that aid when retrieving an AWS secret value.
