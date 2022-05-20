**A tool to automatically deploy https certificates to cloud services.**

This is not an ACME client and is recommended to be used with an ACME client that supports hooks. Of course, you can also use this tool alone.

*[Disclaimer for Mainland China / 针对中国大陆的免责声明](#disclaimer-for-mainland-china)*

## Supported Cloud Provider

- Tencent Cloud (`TencentCloud`)
  - CDN (`cdn`)

## Usage

### Configuration file

This tool will load `./cert-deployer.yaml` based on relative path as configuration file by default. You can specify any file with flag `--profile`. The configuration file cannot be omitted.

Here's a demo of configutation file:

> Please check `config/config.tmpl.yaml` for the latest profile format.

```yaml
log:
  # whether to export to a file
  enable-file: true
  # the directory where the log files are saved
  file-dir: ./
  # debug / info / warn / error
  level: debug
cloud-providers:
  # 'provider' must be in the support list
  - provider: TencentCloud
    secret-id: xxxxxxxxxxxxxxxxxxxxxx
    secret-key: xxxxxxxxxxxxxxxxxxxxxx
```

### Run

Run `cert-deployer help` for help.

A typical usage is:

```bash
cert-deployer \
  --profile /path/to/config.yaml \
  deploy --type cdn \
  --cert /path/to/ceat.pem \
  --key /path/to/private.pem
```

**The value of `type` must be in the support list.**

This command will deploy the certificate to all cloud providers corresponding to the type of asset as much as possible.

## Integration with ACME client

Theoretically, any ACME client with a hook interface can be integrated with this tool -- just run this tool after the certificate is updated.

### [acme.sh](https://github.com/acmesh-official/acme.sh)

Demo:

```bash
acme.sh --issue \
  -d www.example.com \
  -w /www/wwwroot/www.example.com/ \
  --post-hook "cert-deployer --profile /opt/cert-deployer.yaml deploy --type cdn --cert /root/.acme.sh/www.example.com/fullchain.cer --key /root/.acme.sh/www.example.com/www.example.com.key" --force
```

After that, hook command will be saved and apply to `--renew` or `--cron` commands as well. Try `acme.sh --renew -d www.example.com --force` to test.

## Add plugins

If you want to make some contributions to add more back-end support, in general, the steps are as follows：

1. Add a new package in `plugins/`.
2. Add necessary data structures. You may probably want to define a const called `Provider` as the name of the back-end and id.
3. Implement `search.AssetSearcher` and `deploy.Deployer`, and register them by calling `search.MustRegister()` or `deploy.MustRegister()`  in `init()` function.
4. Import your new plugin in `plugins/import.go`.

5. Update the support list in this file.

Congratulations 🥳

> In case you need a new asset type, please add it to `asset/asset_type.go` if it is a generic type (e.g. cdn), otherwise you may want to define them in your package.

## Disclaimer for Mainland China

This is a statement for Chinese mainland only.

坚持非必要不开源，推动大陆开源事业有序绿色发展，倒逼开源项目精品化。谨防恶意开源、恶意编程。本项目自愧不如，故仅在 GitHub 开源。任何人请勿 clone 到 Gitee 等大陆平台（无论是否公开）。我不保证代码里不会出现敏感字符串，例如 `freeSpace = make([]byte, 64, 128)`，以及一些有不良倾向的引导，例如 `go run helloworld.go`，故若您执意 clone 到大陆平台，您需要独自承担一切后果。

