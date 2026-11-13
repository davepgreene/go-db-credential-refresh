module github.com/davepgreene/go-db-credential-refresh/store/awsrds

go 1.27.0

replace github.com/davepgreene/go-db-credential-refresh => ../../

require (
	bou.ke/monkey v1.0.2
	github.com/aws/aws-sdk-go-v2 v1.48.0
	github.com/aws/aws-sdk-go-v2/credentials v1.21.0
	github.com/aws/aws-sdk-go-v2/feature/rds/auth v1.7.11
	github.com/davepgreene/go-db-credential-refresh v1.2.1
	github.com/mitchellh/mapstructure v1.5.0
)

require (
	filippo.io/edwards25519 v1.2.0 // indirect
	github.com/aws/smithy-go v1.28.9 // indirect
	github.com/go-sql-driver/mysql v1.10.1 // indirect
	github.com/lib/pq v1.12.3 // indirect
)
