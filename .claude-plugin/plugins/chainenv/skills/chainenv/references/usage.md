# Usage Workflows

## Typical Read Workflow

1. Run `chainenv diag`.
2. Inspect available keys with `chainenv list` if config is present.
3. Read one key with `chainenv get KEY` or many keys with `chainenv get-env`.
4. Use `--shell bash|zsh|fish` when the output will be evaluated by a shell.

## Typical Write Workflow

1. Confirm the exact key name and backend.
2. Prefer checking existing state with `chainenv get KEY` or `chainenv ls`.
3. Use `chainenv set KEY VALUE` for new secrets.
4. Use `chainenv update KEY VALUE` when the secret should already exist.

## Backend Selection

- Prefer the default `keychain` backend for low-latency local access.
- Use `--backend 1password` only when the user explicitly wants 1Password-backed reads or writes.
- If the user wants faster day-to-day reads but stores canonical secrets in 1Password, suggest copying selected keys into keychain with `chainenv copy --from 1password --to keychain ...`.

## Common Failures

### `keychain backend unavailable`

- On macOS, confirm system keychain access is permitted.
- On Linux, confirm a Secret Service provider is installed and running.

### `failed to load <TOKEN_KEY> from keychain`

- The `["1password"].service_account_token_key` entry points at a missing keychain item.
- Fix the key name or store that token in keychain first.

### `No config found` or `No keys found`

- `chainenv get-env` without explicit keys depends on config discovery.
- Pass explicit keys instead, or create/fix `.chainenv.toml`.

### 1Password Is Slow

This is expected relative to keychain-backed reads. Prefer copying frequently used secrets into keychain when performance matters.
