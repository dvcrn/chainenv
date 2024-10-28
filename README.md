
# chainenv CLI

A simple macOS Keychain CLI wrapper for managing passwords.

## Installation

```
pip install git+https://github.com/dvcrn/chainenv.git
```

## Usage

chainenv <command> [options]


### Commands

#### Get Password
Retrieves a password from the keychain for a specified account.

chainenv get <account>


#### Set Password
Stores a new password in the keychain for a specified account.

chainenv set <account> <password>


#### Update Password
Updates an existing password in the keychain for a specified account.

chainenv update <account> <password>


## Examples


# Get a password
chainenv get myaccount

# Set a new password
chainenv set myaccount mypassword123

# Update an existing password
chainenv update myaccount newpassword123


## Security

This tool uses the macOS Keychain for secure password storage. Passwords are stored using the `security` command-line tool with the following format:
- Service name: `chainenv-<account>`
- Account name: `<account>`
