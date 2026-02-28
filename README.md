# Gochecker

A fast and modern version monitoring tool for software packages, inspired by nvchecker and powered by Go's concurrency model.

Gochecker allows you to track upstream versions from various sources like GitHub, AUR, PyPI, OCI registries, and Git repositories. It helps you stay updated with the latest releases of your favorite software.

## Features

- **Concurrent Scanning**: Leverages Go's goroutines for high-performance, parallel version checks.
- **Multiple Providers**: Support for GitHub, AUR, PyPI, OCI (Container Registries), and Git.
- **Recursive Configuration**: Automatically scan directories for multiple configuration files.
- **Flexible Output**: Supports colorized text, JSON, and YAML output formats for easy integration with other tools.
- **Version Persistence**: Tracks current versions in a local lock file (`.gochecker-lock.json`).

## Installation

### Using Go (recommended)

```bash
go install github.com/cilginc/gochecker@latest
```

### Using AUR (Arch Linux)

```bash
yay -S gochecker
```

## Quick Start

1. **Initialize your configuration:**

   ```bash
   gochecker init
   ```

   This creates a default `.gochecker.yaml` file.

2. **Check for updates:**

   ```bash
   gochecker check
   ```

3. **Update local version records:**
   ```bash
   gochecker update
   ```

## Configuration

The configuration is managed via a `.gochecker.yaml` file. Here is an example:

```yaml
packages:
  - name: gochecker
    github:
      repo: cilginc/gochecker
    prefix: "v"

  - name: pandas
    pypi:
      package: pandas

  - name: dualsensetester
    oci:
      image: ghcr.io/cilginc/dualsense-tester
    prefix: v

  - name: hexecute
    aur:
      package: hexecute
      strip_release: true

  - name: surge
    git:
      url: https://github.com/surge-downloader/Surge.git
    prefix: v
```

### Supported Providers (will add more in the future)

| Provider   | Description                                    | Required Fields     |
| :--------- | :--------------------------------------------- | :------------------ |
| **GitHub** | Track releases or tags from GitHub             | `repo` (owner/repo) |
| **AUR**    | Monitor packages from the Arch User Repository | `package`           |
| **PyPI**   | Check for Python package updates               | `package`           |
| **OCI**    | Monitor container image tags                   | `image`             |
| **Git**    | Track versions from any Git repository         | `url`               |

## CLI Usage

```text
Usage:
  gochecker [command]

Available Commands:
  check       Check for new versions of tracked packages
  init        Create a .gochecker.yaml file with some examples
  list        List all tracked packages
  test        Validate the configuration syntax
  update      Sync local version records

Flags:
  -c, --config string         Path to the configuration file (default ".gochecker.yaml")
  -r, --recursive             Recursively scan a directory for configuration files
  -d, --dir string            Directory to scan when using --recursive (default ".")
  -o, --output string         Output format: 'text', 'json', 'yaml' (default "text")
      --no-color              Disable colorized output
```

## Environment Variables

- `GITHUB_TOKEN`: Recommended for GitHub provider to avoid rate limiting.
- `NO_COLOR`: Disable colorized output.

## License

Distributed under the GPL-3.0 License. See `LICENSE` for more information.
