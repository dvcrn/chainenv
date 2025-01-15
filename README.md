# chainenv CLI

A simple macOS Keychain CLI wrapper for getting/setting secrets from keychain or 1Password, and using them as environment variables.

## Installation

brew

```
brew install dvcrn/formulas/chainenv
```

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
  get         Get a password for an account
  get-env     Get passwords as environment variables
  help        Help about any command
  set         Set a password for an account
  update      Update a password for an existing account

Flags:
      --backend string   Backend to use (keychain or 1password) (default "keychain")
      --debug            Enable debug logging
  -h, --help             help for chainenv
      --vault string     1Password vault to use (default "chainenv")

Use "chainenv [command] --help" for more information about a command.
```

Caveats: 1Password mode is quite slow.


### Commands

#### Get Password
Retrieves a password from the keychain for a specified account.


```
chainenv get <account>
chainenv get <account> --backend 1password
```

#### Set Password
Stores a new password in the keychain for a specified account.


```
chainenv set <account> <password>
chainenv set <account> <password> --backend 1password
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



## Examples


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


### 1Password

When using the 1Password backend, the `1password` CLI is used to retrieve the password. Secrets are stored in the *chainenv* vault by default.
