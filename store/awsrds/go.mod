module github.com/davepgreene/go-db-credential-refresh/store/awsrds

go 1.23.0

toolchain go1.24.0

replace github.com/davepgreene/go-db-credential-refresh => ../../

require (
	bou.ke/monkey v1.0.2
	github.com/aws/aws-sdk-go-v2 v1.37.0
	github.com/aws/aws-sdk-go-v2/credentials v1.17.70
	github.com/aws/aws-sdk-go-v2/feature/rds/auth v1.5.13
	github.com/davepgreene/go-db-credential-refresh v1.0.0
	github.com/mitchellh/mapstructure v1.5.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/aws/smithy-go v1.22.5 // indirect
	github.com/go-sql-driver/mysql v1.9.3 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.3 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgtype v1.14.4 // indirect
	github.com/jackc/pgx/v4 v4.18.3 // indirect
	github.com/jackc/pgx/v5 v5.7.5 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/lib/pq v1.10.9 // indirect
	golang.org/x/crypto v0.40.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/text v0.27.0 // indirect
)
