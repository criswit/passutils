# passutils

Utilities for the [pass](https://www.passwordstore.org/) password manager.

## Prerequisites

- [Go](https://go.dev/) 1.24+
- [pass](https://www.passwordstore.org/) installed and configured with a GPG-encrypted password store at `~/.password-store`

## exporter

Decrypts and exports your entire password store to a single JSON file. The exported structure mirrors your password store's directory hierarchy.

Given a password store like:

```
~/.password-store/
├── personal/
│   ├── email.gpg
│   └── github.gpg
└── work/
    └── vpn.gpg
```

The exporter produces:

```json
{
   "personal": {
      "email": "my-email-password",
      "github": "my-github-token"
   },
   "work": {
      "vpn": "my-vpn-secret"
   }
}
```

### Build

```shell
make exporter
```

This compiles the binary to `./bin/exporter`.

### Usage

```shell
./bin/exporter --outdir <path>
```

| Flag | Default | Description |
|------|---------|-------------|
| `--outdir` | `password-export` | Output directory for the export file. Relative paths resolve from `$HOME`. Absolute paths are used as-is. |

The output file is always named `password-export.json` and is written with `0600` permissions (owner read/write only).

### Examples

```shell
# Export to ~/password-export/password-export.json (default)
./bin/exporter

# Export to ~/my-backup/password-export.json
./bin/exporter --outdir my-backup

# Export to an absolute path
./bin/exporter --outdir /tmp/pass-export
```
