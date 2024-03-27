# Cert Deployer

**A tool to automatically deploy https certificates to cloud services.**

[![Release](https://img.shields.io/github/v/release/ichenhe/cert-deployer?style=flat-square)](https://github.com/ichenhe/cert-deployer/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/ichenhe/cert-deployer?style=flat-square)](https://goreportcard.com/report/github.com/ichenhe/cert-deployer)
[![Build State](https://img.shields.io/github/actions/workflow/status/ichenhe/cert-deployer/check.yml?style=flat-square)](https://github.com/ichenhe/cert-deployer/actions)

This is not an ACME client and is recommended to be used with an ACME client to fully automate the cert workflow. Of course, you can also use this tool alone.

## Supported Cloud Provider

- Tencent Cloud (`TencentCloud`)
  - CDN (`cdn`)

## Quick start

```bash
# deploy cert to all matched cdn in TencentCloud
./cert-deployer --provider TencentCloud \
--secret-id "xxxxxxxxx" \
--secret-key "yyyyyyyyyy" \
--cert "/path/to/fullchain.pem" \
--key "/path/to/privkey.pem" \
--type cdn
```

The value of `provider` / `type` must be in the support list.

## Usage

The global flag `--profile` can be used to specify the configuration file.

### Pre-defined deployment

Record all required parameters for deployment in the profile to reuse them or deploy many certs at once.

```yaml
# one provider indicates one account in a cloud platform
cloud-providers:
  example: # name (id) of the provider used as a reference
    provider: TencentCloud # 'provider' must be in the support list
    secret-id: "xxxxxxxxxxxxxxxxxxxxx"
    secret-key: "xxxxxxxxxxxxxxxxxxxxx"
  # more providers may follow...

deployments:
  tencent-cdn: # name (id) of the deployment used as a reference
    provider-id: example # must be in `cloud-providers`
    cert: "/fullchain.pem"
    key: "/private.pem"
    # what assets do you want to deploy to?
    assets:
      - type: cdn
        id: "xxx" # id can be omitted which means deploying to all applicable cdns
      # more assets may follow...
  # more deployments may follow...
```

Now you can execult defined deployment:

```bash
./cert-deployer --profile "/path/to/profile.yml" deploy --deployment tencent-cdn
```

Or provide many deployments with short flag `-d`:

```bash
./cert-deployer --profile "/path/to/profile.yml" deploy -d deploy1 -d deploy2
```

### Trigger

Trigger can execute deployments automatically based on event. Trigger must be defined in profile, no cli-only approach.

**To use trigger, must keep the program running.**

The following example will monitor the file `/path/to/cert.pem`. Once its content (hash) changed, deployments `deploy1` and `deploy2` will be executed. However, the program will wait for 1000ms before really execute the deployment. The timer will be reset if another event occurs during the wait.

Please be aware that:

- The file to be monitored doesn't have to be the same as certificate.
- Waitting time is recommended because the writing of the file may be completed multiple times.

```yaml
deployments:
  deploy1: # omit ...
  deploy2: # omit ...

triggers:
  cert: # name (id) of the trigger
    type: file_monitoring # only `file_monitoring` for now
    deployments: [ "deploy1", "deploy2" ] # must be in `deployments`
    options:
      file: "/path/to/cert.pem" # file to monitor, not a folder
      wait: 1000 # ms wait for before executing deployemnts
  # more triggersa may follow...
```

Run the following command to start all triggers:

```bash
./cert-deployer --profile "/path/to/profile.yml" run
```

### Logging

By default, cert-deployer write logs to stdout (console) with INFO level.

Logging behaviour can be specified in the profile:

```yaml
log-drivers:
  # write logs to console
  - driver: stdout
    level: debug    # debug / info / warn / error 
    format: fluent # json / fluent (human-readable)

  # write logs to a file
  - driver: file
    level: info
    format: json
    options:
      # where to write to, can not be a folder
      file: "/var/log/cert-deployer/log.json"
  # more drivers may follow...
```

Currently only driver `stdout` and `file` are supported. Many drivers can be defined, even two `file` driver with different target file.

Please note that default logger will be disabled if you specify an empty `log-drivers`, no log output in this case. Delete the `log-drivers` itself if you want to keep default behavior.

> Currently, no built-in log-rotation feature. Please use external program such as `logrotate` in linux.

### Profile

This program will try to find and load `cert-deployer.yaml` based on working dir. However, you are encouraged to specify the file with gloabl flag `--profile`.

Please check `config/config.tmpl.yaml` for the latest profile format. A JSON Schema file is also provided which describes the profile format very well.

## Integration with ACME client

cert-deployer can work with ACME client in two ways:

- [RECOMMAND] Use trigger feature to monitor the cert file generated by ACME client.
- Use hook feature of ACME client to execute cert-deployer.



Here's an example to work with [acme.sh](https://github.com/acmesh-official/acme.sh) via hook.

```bash
acme.sh --issue \
  -d www.example.com \
  -w /www/wwwroot/www.example.com/ \
  --post-hook "cert-deployer deploy --type cdn --cert /root/.acme.sh/www.example.com/fullchain.cer --key /root/.acme.sh/www.example.com/www.example.com.key --provider TencentCloud --secret-id xxxx --secret-key yyyyy" --force
```

After that, hook command will be saved and apply to `--renew` or `--cron` commands as well. Try `acme.sh --renew -d www.example.com --force` to test.

## Migrating from v0.1

Legacy usage is no longer supported, which means you shouldn't specify the cloud provider in profile while provide the target asset or cert file via cli.

Instead, you can either:

- Execute fully custom deployment as described in *quck start*.
- Define everything in profile as described in *pre-defined deployment*.

In addition this, you are encouraged to use *trigger* to integrate with ACME client instead of hook, which is more easier and clearer.

## Add plugins

If you want to make some contributions to add more back-end support, in general, the steps are as follows：

1. Add a new package in `plugins/`.
2. Add necessary data structures. You may probably want to define a const called `Provider` as the name of the back-end and id.
3. Implement `domain.Deployer`, and register it by calling `registry.MustRegister()` in `init()` function.
4. Import your new plugin in `plugins/import.go`.
5. Update the support list in this file.

Congratulations 🥳

> In case you need a new asset type, please add it to `asset/asset_type.go` if it is a generic type (e.g. cdn), otherwise you may want to define them in your package.

## Disclaimer for Mainland China

This is a statement for Chinese mainland only.

我们不会故意，但亦不能保证整个仓库中不包含（潜在的）敏感内容，因此不鼓励任何人将本仓库镜像到大陆平台。若您执意这么做，后果自行承担。

