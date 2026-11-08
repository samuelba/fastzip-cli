# fastzip-cli

This project adds a command line interface (cli) for the [fastzip project](https://github.com/saracen/fastzip).

## Install

Download the Linux binary for your architecture (`amd64` or `arm64`) from the [latest release](https://github.com/samuelba/fastzip-cli/releases/latest):

```bash
curl -fsSLo fastzip-cli https://github.com/samuelba/fastzip-cli/releases/latest/download/fastzip-cli-linux-amd64
chmod +x fastzip-cli
```

## Example

### Compress

```bash
fastzip-cli compress /path/to/folder archive.zip
```

### Extract

```bash
fastzip-cli extract archive.zip /path/to/folder
```
