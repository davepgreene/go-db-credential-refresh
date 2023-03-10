module github.com/davepgreene/go-db-credential-refresh/store/awsrds

go 1.15

replace github.com/davepgreene/go-db-credential-refresh => ../../

require (
	bou.ke/monkey v1.0.2
	github.com/aws/aws-sdk-go-v2 v1.17.6
	github.com/aws/aws-sdk-go-v2/credentials v1.13.16
	github.com/aws/aws-sdk-go-v2/feature/rds/auth v1.1.19
	github.com/davepgreene/go-db-credential-refresh v0.0.0-00010101000000-000000000000
	github.com/mitchellh/mapstructure v1.5.0
)
