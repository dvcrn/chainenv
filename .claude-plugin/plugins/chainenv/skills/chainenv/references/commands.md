# Chainenv Commands

## Install

```bash
brew install dvcrn/formulas/chainenv
npm install -g @dvcrn/chainenv
go install github.com/dvcrn/chainenv@latest
```

## Diagnostics

```bash
chainenv diag
```

## Read Operations

```bash
chainenv ls
chainenv ls --backend 1password

chainenv list

chainenv get GITHUB_TOKEN
chainenv get GITHUB_TOKEN --backend 1password

chainenv get-env GITHUB_TOKEN,OPENAI_API_KEY --shell bash
chainenv get-env GITHUB_TOKEN,OPENAI_API_KEY --shell zsh
chainenv get-env GITHUB_TOKEN,OPENAI_API_KEY --shell fish

chainenv get-env --shell bash
```

## Write Operations

```bash
chainenv set GITHUB_TOKEN secret-value
chainenv set FEATURE_FLAG true --default true

chainenv update GITHUB_TOKEN new-secret-value
chainenv update GITHUB_TOKEN new-secret-value --backend 1password
```

`set` also upserts the key into `.chainenv.toml` or `chainenv.toml` when a config file is present or created by the command flow.

## Copy Between Backends

```bash
chainenv copy --from 1password --to keychain GITHUB_TOKEN,OPENAI_API_KEY
chainenv copy --from keychain --to 1password GITHUB_TOKEN
chainenv cp --from 1password --to keychain GITHUB_TOKEN --overwrite
```

## Shell Usage

```bash
eval "$(chainenv get-env --shell bash)"
eval "$(chainenv get-env GITHUB_USERNAME,GITHUB_PASSWORD --shell zsh)"
eval (chainenv get-env GITHUB_USERNAME,GITHUB_PASSWORD --shell fish)
```
