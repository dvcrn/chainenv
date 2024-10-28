# chainenv CLI

A simple macOS Keychain CLI wrapper for getting/setting secrets from keychain, and using them as environment variables.

## Installation


pip install git+https://github.com/dvcrn/chainenv.git


## Usage


```
chainenv <command> [options]
```


### Commands

#### Get Password
Retrieves a password from the keychain for a specified account.


```
chainenv get <account>
```


#### Set Password
Stores a new password in the keychain for a specified account.


```
chainenv set <account> <password>
```


#### Update Password
Updates an existing password in the keychain for a specified account.


```
chainenv update <account> <password>
```


#### Get Multiple Passwords as Environment Variables
Retrieves multiple passwords and outputs them as shell exports.


```
chainenv get-env <account1,account2,...> [--fish|--bash|--zsh]
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
