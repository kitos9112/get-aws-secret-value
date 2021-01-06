# ${APP}

[![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/kitos9112/aws-get-secret-value.svg)](https://hub.docker.com/r/kitos9112/aws_get_secret_value/tags)

Retrieves an AWS secret value as-is after given its secret name. The script can read an AWS Profile name as well as a credentials file in the CLI, otherwise it will automatically load its settings following AWS SDK standards

* Environment Variables
* Shared Credentials file
* Shared Configuration file (if SharedConfig is enabled)
* EC2 Instance Metadata (credentials only)

```bash
> aws-get-secret-value
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
# Linux
curl -L https://github.com/kitos9112/${APP}/releases/download/${VERSION}/${APP}_${VERSION}_linux_x86_64.tar.gz | tar xz

# OS X
curl -L https://github.com/kitos9112/${APP}/releases/download/${VERSION}/${APP}_${VERSION}_osx_x86_64.tar.gz | tar xz

# Windows
curl -LO https://github.com/kitos9112/${APP}/releases/download/${VERSION}/${APP}_${VERSION}_windows_x86_64.zip
unzip ${APP}_${VERSION}_windows_x86_64.zip
```

## Use it

```text

${APP} [OPTIONS] [COMMAND [ARGS...]]

${USAGE}
```

### Examples

The simplest example that could easily be integrated into a CICD pipeline:

```shell
> export AWS_PROFILE=myAwsProfile
> export AWS_REGION=eu-west-1
> aws-get-secret-value --secret-name my_secret_name
mySecretValue
```
