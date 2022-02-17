# ${APP}

[![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/kitos9112/get-aws-secret-value.svg)](https://hub.docker.com/r/kitos9112/aws_get_secret_value/tags)

Retrieves an AWS secret value as-is and throws its content to `stdout` in plain UTF-8 encoding.

Capable of reading an AWS Profile name as well as a credentials file from the `CLI`

Defaults to AWS SDK standards for order of precedence, most to least:

* Environment variables.
* Shared credentials file.
* Shared Configuration file (if SharedConfig is enabled)
* If your application uses an ECS task definition or RunTask API operation, IAM role for tasks.
* If your application is running on an Amazon EC2 instance, IAM role for Amazon EC2.
* If your application is running on an EKS cluster with IRSA enabled, IAM role for pods.

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

- [${APP}](#app)
  - [Get it](#get-it)
  - [Use it](#use-it)
    - [Examples](#examples)

<!-- /TOC -->

## Get it

Using go get:

```bash
go get -u github.com/kitos9112/${APP}
```

Or [download the binary](https://github.com/kitos9112/${APP}/releases/latest) from the releases page.

```bash
# Linux x86_64
curl -L https://github.com/kitos9112/${APP}/releases/download/${VERSION}/${APP}_${VERSION}_linux_x86_64.tar.gz | tar xz

# Linux arm64
curl -L https://github.com/kitos9112/${APP}/releases/download/${VERSION}/${APP}_${VERSION}_linux_arm64.tar.gz | tar xz

# OS X x86_64
curl -L https://github.com/kitos9112/${APP}/releases/download/${VERSION}/${APP}_${VERSION}_osx_x86_64.tar.gz | tar xz

# OS X arm64
curl -L https://github.com/kitos9112/${APP}/releases/download/${VERSION}/${APP}_${VERSION}_osx_arm64.tar.gz | tar xz

# Windows x86_64
curl -LO https://github.com/kitos9112/${APP}/releases/download/${VERSION}/${APP}_${VERSION}_windows_x86_64.zip
unzip ${APP}_${VERSION}_windows_x86_64.zip
```

## Use it

```text

${APP} [OPTIONS] [COMMAND [ARGS...]]

${USAGE}
```

### Examples

A very basic examples:

```shell
# 1- Rely on ENV variables
> export AWS_PROFILE=myAwsProfile
> export AWS_REGION=eu-west-1
> get-aws-secret-value --secret-name my_secret_name
mySecretValue

# 2- Override or set AWS-specific commands via CLI options
> get-aws-secret-value --profile myAwsPRofile --credFile /tmp/myAwsCredFile --secret-name my_secret_name
mySecretValue
```

Or in case you leverage IaC using Terragrunt, you could retrieve the value of an AWS secret previously created and pre-populated with more complex data structures (e.g. JSON)

``` hcl
# terragrunt.hcl

locals {
  my_secret_var1 = lookup(jsondecode(run_cmd("--terragrunt-quiet", "/usr/local/bin/aws-get-secret-value", "--secret-name", "my_secret", "--aws-region", "eu-west-1")), "secretKey1")
}

inputs = {
  my_password = local.my_secret_var1
}
```

As you can see, a simple cross-platform binary file could be utilised in many scenarios that aid at retrieving an AWS secret value presented at stdout.
