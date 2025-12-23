# chainenv CLI

A simple macOS Keychain and Linux Secret Service wrapper for getting/setting secrets from keychain or 1Password, and using them as environment variables.

Also available as a [mise plugin](https://github.com/dvcrn/mise-chainenv)

## Installation

brew

```
brew install dvcrn/formulas/chainenv
```

npm

```
npm install -g @dvcrn/chainenv
```

(this tool is not written in JavaScript, just distributed through NPM)

... or with Go:

```
go install github.com/dvcrn/chainenv@latest
```

To use 1Password functionality, you need to have the 1Password CLI installed as well:

```
brew install 1password-cli
```

## Usage

```
Usage:
  chainenv [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  copy        Copy passwords between backends
  generate-env Generate environment exports from config
  get         Get a password for an account
  get-env     Get passwords as environment variables
  help        Help about any command
  list        List keys declared in config
  ls          List all stored accounts
  set         Set a password for an account
  update      Update a password for an existing account
  diag        Diagnose available backends

Flags:
      --backend string   Backend to use (keychain or 1password) (default "keychain")
      --debug            Enable debug logging
  -h, --help             help for chainenv
      --vault string     1Password vault to use (default "chainenv")
```

### Note on 1Password

Caveats: 1Password mode is very very slow. This is sped-up somewhat by using goroutines to parallelize the requests, but it's still slow.

My recommendation: Use `chainenv cp --from 1password --to keychain` to copy passwords from 1Password to the keychain, then use keychain for fast access.

```
❯ time ./chainenv get-env TEST,TEST2,Test3
TEST='123'
TEST2='123'

________________________________________________________
Executed in   36.67 millis    fish           external
   usr time   19.28 millis   30.00 micros   19.25 millis
   sys time   19.04 millis  502.00 micros   18.54 millis


❯ time ./chainenv get-env TEST,TEST2,Test3 --backend 1password
TEST='123'
TEST2='123'

________________________________________________________
Executed in    3.42 secs      fish           external
   usr time  256.23 millis   48.00 micros  256.18 millis
   sys time  129.46 millis  587.00 micros  128.87 millis
```

### Commands

#### List Accounts

Lists all accounts that have passwords stored in the configured backend.

```
chainenv ls
chainenv ls --backend 1password
```

#### List Config Keys

Lists keys declared in `.chainenv.toml` or `chainenv.toml`.

```
chainenv list
```

#### Get Password

Retrieves a password from the keychain for a specified account.

```
chainenv get <account>
chainenv get <account> --backend 1password
```

#### Set Password

Stores a new password in the keychain for a specified account.
Also registers the key in `.chainenv.toml` (or `chainenv.toml` if found).

```
chainenv set <account> <password>
chainenv set <account> <password> --backend 1password
chainenv set <account> <password> --default <value>
```

#### Update Password

Updates an existing password in the keychain for a specified account.

```
chainenv update <account> <password>
chainenv update <account> <password> --backend 1password
```

#### Get Multiple Passwords as Environment Variables

Retrieves multiple passwords and outputs them as shell exports.

```
chainenv get-env <account1,account2,...> [--fish|--bash|--zsh]
chainenv get-env <account1,account2,...> [--fish|--bash|--zsh] --backend 1password
```

... will output

```
account1='foo'
account2='bar'
```

or with --zsh

```
export account1='foo'
export account2='bar'
```

### Copy Passwords

```
chainenv cp --from <backend> --to <backend> ITEM1,ITEM2
```

### Generate Environment Exports

Outputs shell exports for all keys declared in config.

```
chainenv generate-env
chainenv generate-env --shell fish
```

## Project Config

If `.chainenv.toml` or `chainenv.toml` exists, `chainenv` will read it and use it to:
- List configured keys (`chainenv list`)
- Provide default fallbacks when a secret is missing
- Generate shell exports (`chainenv generate-env`)

Lookup order:
1. `.chainenv.toml`
2. `chainenv.toml`

Example config:

```
[[keys]]
name = "GITHUB_TOKEN"
provider = "keychain"
default = ""

[[keys]]
name = "SOME_FLAG"
provider = "keychain"
default = "true"
```

Notes:
- `default` values are stored in plaintext.
- `provider` can be `keychain` or `1password`.
- If a key has a `default` and the secret is missing, `chainenv get` and `chainenv get-env` will use the default.

## Examples

### List all stored accounts

```
chainenv ls
```

### List configured keys

```
chainenv list
```

### Get a password

```
chainenv get myaccount
```

### Set a new password

```
chainenv set myaccount mypassword123
```

### Update an existing password

```
chainenv update myaccount newpassword123
```

### Get multiple passwords as environment variables

```
chainenv get-env GITHUB_USERNAME,GITHUB_PASSWORD,AWS_KEY
```

### Generate environment exports

```
chainenv generate-env
```

## Usage in Shell Environments

### Bash/Zsh

#### Individual password retrieval

```
export GITHUB_USERNAME=$(chainenv get GITHUB_USERNAME)
export GITHUB_PASSWORD=$(chainenv get GITHUB_PASSWORD)
```

#### Multiple passwords at once

```
eval $(chainenv get-env GITHUB_USERNAME,GITHUB_PASSWORD --bash)
```

### Fish

#### Individual password retrieval

```
set -gx GITHUB_USERNAME (chainenv get GITHUB_USERNAME)
set -gx GITHUB_PASSWORD (chainenv get GITHUB_PASSWORD)
```

#### Multiple passwords at once

```
eval (chainenv get-env GITHUB_USERNAME,GITHUB_PASSWORD --fish)
```

### Direnv (.envrc)

Generate exports from config (write to `.envrc`):

```
chainenv generate-env > .envrc
```

#### Individual password retrieval

```
export GITHUB_USERNAME="$(chainenv get GITHUB_USERNAME)"
export GITHUB_PASSWORD="$(chainenv get GITHUB_PASSWORD)"
```

#### Multiple passwords at once

```
eval "$(chainenv get-env GITHUB_USERNAME,GITHUB_PASSWORD,AWS_KEY --bash)"
```

## Security

This tool uses the macOS Keychain for secure password storage. Passwords are stored using the `security` command-line tool with the following format:

- Service name: `chainenv-<account>`
- Account name: `<account>`

### Linux Keychain Support

On Linux, the `keychain` backend uses the Secret Service API via the system keyring (e.g., GNOME Keyring or KWallet). A running Secret Service provider is required (typically present on desktop distributions). On minimal/server installs you may need to install and start a keyring daemon.

- GNOME-based distros: GNOME Keyring usually preinstalled; `secret-tool` available via `libsecret-tools`.
- KDE Plasma: KWallet is typically installed; enable Secret Service integration in system settings.
- If no provider is available, commands using the `keychain` backend will return an error indicating how to enable one.

Stored items are grouped under the service `chainenv` with each account stored under its own key (the account name).

### 1Password

When using the 1Password backend, the `1password` CLI is used to retrieve the password. Secrets are stored in the _chainenv_ vault by default.

#### Diagnose Backends

Checks which backends are available on the current system.

```
chainenv diag
```

Sample output:

```
Backend diagnostics:
- macOS Keychain: available
- Linux Keyring (Secret Service/KWallet): unavailable (not Linux)
- 1Password CLI: available (signed in)
```
