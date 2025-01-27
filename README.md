<p align="center"><a href="#readme"><img src=".github/images/card.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/w/rds-cli-completion-generator/ci-push"><img src="https://kaos.sh/w/rds-cli-completion-generator/ci-push.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/w/rds-cli-completion-generator/codeql"><img src="https://kaos.sh/w/rds-cli-completion-generator/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="#license"><img src=".github/images/license.svg"/></a>
</p>

<p align="center"><a href="#installation">Installation</a> • <a href="#command-line-completion">Command-line completion</a> • <a href="#man-documentation">Man documentation</a> • <a href="#usage">Usage</a> • <a href="#ci-status">CI Status</a> • <a href="#contributing">Contributing</a> • <a href="#license">License</a></p>

<br/>

`rds-cli-completion-generator` is a a tool to generate completion for [RDS](https://kaos.sh/rds) CLI.

### Installation

#### From source

To build the `rds-cli-completion-generator` from scratch, make sure you have a working Go 1.22+ workspace (_[instructions](https://go.dev/doc/install)_), then:

```
go install github.com/essentialkaos/rds-cli-completion-generator@latest
```

#### Prebuilt binaries

You can download prebuilt binaries for Linux and macOS from [EK Apps Repository](https://apps.kaos.st/rds-cli-completion-generator/latest):

```bash
bash <(curl -fsSL https://apps.kaos.st/get) rds-cli-completion-generator
```

#### Container Image

The latest version of `rds-cli-completion-generator` also available as container image on [GitHub Container Registry](https://kaos.sh/p/rds-cli-completion-generator) and [Docker Hub](https://kaos.sh/d/rds-cli-completion-generator):

```bash
podman run --rm -it ghcr.io/essentialkaos/rds-cli-completion-generator:latest
# or
docker run --rm -it ghcr.io/essentialkaos/rds-cli-completion-generator:latest
```

### Command-line completion

You can generate completion for `bash`, `zsh` or `fish` shell.

Bash:
```bash
sudo rds-cli-completion-generator --completion=bash 1> /etc/bash_completion.d/rds-cli-completion-generator
```

ZSH:
```bash
sudo rds-cli-completion-generator --completion=zsh 1> /usr/share/zsh/site-functions/rds-cli-completion-generator
```

Fish:
```bash
sudo rds-cli-completion-generator --completion=fish 1> /usr/share/fish/vendor_completions.d/rds-cli-completion-generator.fish
```

### Man documentation

You can generate man page using next command:

```bash
rds-cli-completion-generator --generate-man | sudo gzip > /usr/share/man/man1/rds-cli-completion-generator.1.gz
```

### Usage

<p align="center"><img src=".github/images/usage.svg"/></p>

### CI Status

| Branch | Status |
|--------|----------|
| `master` | [![CI](https://kaos.sh/w/rds-cli-completion-generator/ci-push.svg?branch=master)](https://kaos.sh/w/rds-cli-completion-generator/ci-push?query=branch:master) |
| `develop` | [![CI](https://kaos.sh/w/rds-cli-completion-generator/ci-push.svg?branch=develop)](https://kaos.sh/w/rds-cli-completion-generator/ci-push?query=branch:develop) |

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/contributing-guidelines#contributing-guidelines).

### License

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
