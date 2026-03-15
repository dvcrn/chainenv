---
name: chainenv
description: Operate the `chainenv` CLI for local secret workflows across macOS Keychain, Linux keyring, and optional 1Password integration. Use when requests mention `chainenv`, `.chainenv.toml`, `chainenv.toml`, keychain vs 1Password, shell export generation, copying secrets between backends, or troubleshooting backend availability and `op` token loading.
---

# Chainenv

Use this skill when the user wants to work with the installed `chainenv` CLI. This skill is for operating the product, not for changing its source code.

## Prerequisites

- `chainenv` must be installed and available on `PATH`.
- Start with `chainenv diag` to confirm backend availability on the current machine.
- For 1Password flows, the `op` CLI must be installed. Authentication can come from an active `op` session or from `OP_SERVICE_ACCOUNT_TOKEN`.
- If config-driven behavior matters, read `references/auth.md` for `.chainenv.toml` layout and token-loading behavior.

## When To Use It

Use this skill for:

- backend diagnostics
- reading stored secrets with `get` or `get-env`
- generating shell exports for bash, zsh, fish, or plain output
- writing or updating secrets with `set` and `update`
- copying secrets between keychain and 1Password
- explaining how `.chainenv.toml` affects provider selection and default fallbacks

Read `references/commands.md` when you need exact command forms.
Read `references/usage.md` when you need workflow guidance or troubleshooting.

## Operating Guidance

1. Confirm the CLI is installed and run `chainenv diag`.
2. Prefer read-only commands first: `ls`, `list`, `get`, or `get-env`.
3. Check whether config should drive provider selection or default fallbacks.
4. Use write commands only after the target key and backend are explicit.
5. Prefer keychain for day-to-day reads when 1Password latency is a concern.

## Command Selection

- Use `chainenv ls` to list all stored accounts in the selected backend.
- Use `chainenv list` to list keys declared in config.
- Use `chainenv get <KEY>` for one secret.
- Use `chainenv get-env ... --shell <shell>` for multi-key shell exports.
- Use `chainenv get-env --shell <shell>` with no keys only when config files define the keys to load.
- Use `chainenv set` to create a secret and register it in config.
- Use `chainenv update` to change an existing secret.
- Use `chainenv copy` or `chainenv cp` to move secrets between backends.

Prefer `--shell fish|bash|zsh` for new examples. The legacy `--fish`, `--bash`, and `--zsh` flags still work when matching older user setups.

## Failure Guidance

- If `chainenv diag` shows the keychain backend unavailable, the host may not support the selected backend or the keyring service is missing.
- If 1Password commands fail, check for `op` on `PATH`, signed-in state, and whether `OP_SERVICE_ACCOUNT_TOKEN` is set or can be loaded from config.
- If `get-env` without explicit keys prints `No config found` or `No keys found`, switch to explicit key arguments or fix the config file.
- `default` values in config are plaintext fallbacks, not encrypted secrets.
- Prefer live `chainenv --help` and subcommand help over stale documentation when they differ.
