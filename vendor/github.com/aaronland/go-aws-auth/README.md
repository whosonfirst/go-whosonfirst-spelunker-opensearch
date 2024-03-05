# go-aws-auth

Go package providing methods and tools for determining or assigning AWS credentials.

This package targets [aws-sdk-go-v2](https://github.com/aws/aws-sdk-go-v2/). For similar functionality targeting `aws-sdk-go` please consult the [aaronland/go-aws-session](https://github.com/aaronland/go-aws-session) package.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/aaronland/go-aws-auth.svg)](https://pkg.go.dev/github.com/aaronland/go-aws-auth)

## Tools

```
$> make cli
go build -mod vendor -o bin/aws-mfa-session cmd/aws-mfa-session/main.go
go build -mod vendor -o bin/aws-get-credentials cmd/aws-get-credentials/main.go
go build -mod vendor -o bin/aws-set-env cmd/aws-set-env/main.go
```

### aws-get-credentials

`aws-get-credentials` is a command line tool to emit one or more keys from a given profile in an AWS .credentials file.

```
$> ./bin/aws-get-credentials -h
Usage of ./bin/aws-get-credentials:
  -profile string
    	A valid AWS credentials profile (default "default")
```

### aws-mfa-session

`aws-mfa-session` is a command line to create session-based authentication keys and secrets for a given profile and multi-factor authentication (MFA) token and then writing that key and secret back to a "credentials" file in a specific profile section.

```
$> ./bin/aws-mfa-session -h
Usage of ./bin/aws-mfa-session:
  -duration string
    	A valid ISO8601 duration string indicating how long the session should last (months are currently not supported) (default "PT1H")
  -profile string
    	A valid AWS credentials profile (default "default")
  -session-profile string
    	The name of the AWS credentials profile to update with session credentials (default "session")
```

For example:

```
$> ./bin/aws-mfa-session -profile {PROFILE} -duration PT8H
Enter your MFA token code: 123456
2018/07/26 09:47:09 Updated session credentials for 'session' profile, expires Jul 26 17:47:09 (2018-07-27 00:51:52 +0000 UTC)
```

### aws-set-env

`aws-set-env` is a command line tool to assign required AWS authentication environment variables for a given profile in a AWS .credentials file.

```
$> ./bin/aws-set-env -h
Usage of ./bin/aws-set-env:
  -profile string
    	A valid AWS credentials profile (default "default")
  -session-token
    	Require AWS_SESSION_TOKEN environment variable (default true)
```

## See also:

* https://github.com/aws/aws-sdk-go-v2/
