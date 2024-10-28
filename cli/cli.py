
#!/usr/bin/env python3
import argparse
import subprocess
import sys

def get_password(account):
    try:
        cmd = ['security', 'find-generic-password', '-a', account, '-s', f'chainenv-{account}', '-w']
        result = subprocess.run(cmd, capture_output=True, text=True)
        if result.returncode == 0:
            return result.stdout.strip()
        return None
    except Exception as e:
        print(f"Error retrieving password: {e}", file=sys.stderr)
        return None

def set_password(account, password, update=False):
    try:
        cmd = ['security', 'add-generic-password', '-a', account, '-s', f'chainenv-{account}', '-w', password, '-j', 'Set by chainenv']
        if update:
            cmd.append('-U')
        result = subprocess.run(cmd, capture_output=True, text=True)
        if result.returncode != 0:
            raise Exception(result.stderr.strip())
        return True    
    except Exception as e:
        print(f"Error setting password: {e}", file=sys.stderr)
        return False

def main():
    parser = argparse.ArgumentParser(description='Simple macOS Keychain CLI wrapper')
    subparsers = parser.add_subparsers(dest='command', help='Commands')
    
    # Get command
    get_parser = subparsers.add_parser('get', help='Get password from keychain')
    get_parser.add_argument('account', help='Account name to retrieve')
    
    # Set command
    set_parser = subparsers.add_parser('set', help='Set password in keychain')
    set_parser.add_argument('account', help='Account name to set')
    set_parser.add_argument('password', help='Password to store')
    
    # Update command
    update_parser = subparsers.add_parser('update', help='Update password in keychain')
    update_parser.add_argument('account', help='Account name to update')
    update_parser.add_argument('password', help='Password to store')
    
    args = parser.parse_args()
    
    if args.command == 'get':
        password = get_password(args.account)
        if password:
            print(password)
        else:
            print("Password not found", file=sys.stderr)
            sys.exit(1)
    elif args.command == 'set':
        if set_password(args.account, args.password):
            print(f"Password set for {args.account}")
        else:
            print("Failed to set password", file=sys.stderr)
            sys.exit(1)
    elif args.command == 'update':
        if set_password(args.account, args.password, update=True):
            print(f"Password updated for {args.account}")
        else:
            print("Failed to update password", file=sys.stderr)
            sys.exit(1)
    else:
        parser.print_help()
        sys.exit(1)

if __name__ == '__main__':
    main()