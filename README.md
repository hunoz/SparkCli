# SparkCli
## Description
The Spark CLI is used to perform some operations in AWS Cognito, such as password change, forgot password, performing first sign in, and registering a TOTP device. This CLI will compliment future CLIs such as Haze and Maroon which will utilize this CLI's data to perform their own operations. More information on those will be available when they are ready for release.

This repository should be forked and the GitHub secrets and workflows updated to match the correct Cognito pool ID, pool region, and client ID! Failure to do this will cause it to not run.

## Installation
1. Navigate to the [releases page](https://github.com/hunoz/SparkCli/releases) and download the binary for your operating system. If you do not see your operating system, please submit an issue with your OS and ARCH so that it can be added.
2. Place the binary in a location in your PATH (e.g. /usr/local/bin/spark)

## Usage
If you have not signed into the Cognito pool, please navigate to the `First Sign In` section.

### First Sign In
First Sign In is used if you have been added or created in a Cognito pool but have not performed a first sign in to verify your email and change your password.
1. Run `spark first-sign-in`
2. Follow the prompts as necessary

### Auth
Auth performs an auth to your Cognito pool and stores the token information for use with other tools that call Cognito-backed endpoints
1. Run `spark auth`. If you want to force a refresh of your tokens, also add `--force`
2. Follow the prompts to authenticate to the pool

### Change Password
Change Password will allow you to change your password. This operation uses the Cognito default password requirements.
1. Run `spark change-password`
2. Follow the prompts to change your password

### Register TOTP
Register TOTP will register a TOTP device to you. This does not currently support SMS, but may in the future.
1. Run `spark register-totp`
2. Follow the prompts to register a TOTP device

### Reset Password
Reset Password is for if you do not know your current password but have previously performed an initial sign in.
1. Run `spark reset-password`
2. Follow the prompts to reset your password

## Roadmap
* Add in an `init` subcommand that will place the pool ID, client ID, and region in the config, so that forking this repository isn't necessary and it can be used out of the box.
