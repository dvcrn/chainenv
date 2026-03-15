# Authentication And Config

## Required Runtime

- `chainenv` must be installed on `PATH`.
- `op` is only required for the `1password` backend.

## 1Password Expectations

`chainenv` supports 1Password in two ways:

- an active `op` CLI session
- `OP_SERVICE_ACCOUNT_TOKEN` in the environment

If `OP_SERVICE_ACCOUNT_TOKEN` is already set, `chainenv` uses it directly.

If the variable is not set and config contains a `["1password"].service_account_token_key`, `chainenv` tries to read that key from the keychain backend and export it into `OP_SERVICE_ACCOUNT_TOKEN` for the current process.

## Config Files

Config lookup prefers:

1. `.chainenv.toml`
2. `chainenv.toml`

Minimal example:

```toml
["1password"]
service_account_token_key = "CHAINENV_OP_TOKEN"

[[keys]]
name = "GITHUB_TOKEN"
provider = "keychain"
default = ""

[[keys]]
name = "OPENAI_API_KEY"
provider = "1password"
```

## Config Effects

- `chainenv list` reads declared keys from config.
- `chainenv get-env` with no key arguments loads key names from config.
- Per-key `provider` overrides the global `--backend` default.
- Per-key `default` is used when a secret is missing.

## Caveats

- `default` values are plaintext.
- If the configured token key is missing from keychain, 1Password commands fail before backend operations begin.
- The default 1Password vault is `chainenv`, but callers can override it with `--vault`.
