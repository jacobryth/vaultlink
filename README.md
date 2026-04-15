# vaultlink

> A CLI tool for syncing HashiCorp Vault secrets to local `.env` files with role-based filtering

---

## Installation

```bash
go install github.com/yourusername/vaultlink@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/vaultlink/releases).

---

## Usage

Set your Vault address and token, then run `vaultlink` to sync secrets to a local `.env` file.

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.your-vault-token"

# Sync secrets from a Vault path to a .env file
vaultlink sync --path secret/myapp --role backend --output .env

# Preview secrets without writing to disk
vaultlink sync --path secret/myapp --role frontend --dry-run
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--path` | Vault secret path | _(required)_ |
| `--role` | Filter secrets by role | `default` |
| `--output` | Output `.env` file path | `.env` |
| `--dry-run` | Print secrets without writing | `false` |

---

## Configuration

You can also use a `vaultlink.yaml` config file in your project root:

```yaml
vault_addr: https://vault.example.com
path: secret/myapp
role: backend
output: .env
```

---

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

---

## License

[MIT](LICENSE)