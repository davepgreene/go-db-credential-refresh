module github.com/davepgreene/go-db-credential-refresh/store/awsrds

go 1.23.0

toolchain go1.24.0

replace github.com/davepgreene/go-db-credential-refresh => ../../

require (
	bou.ke/monkey v1.0.2
	github.com/aws/aws-sdk-go-v2 v1.36.3
	github.com/aws/aws-sdk-go-v2/credentials v1.17.63
	github.com/aws/aws-sdk-go-v2/feature/rds/auth v1.5.11
	github.com/davepgreene/go-db-credential-refresh v1.0.0
	github.com/mitchellh/mapstructure v1.5.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/aws/smithy-go v1.22.3 // indirect
	github.com/go-sql-driver/mysql v1.9.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.3 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgtype v1.14.4 // indirect
	github.com/jackc/pgx/v4 v4.18.3 // indirect
	github.com/jackc/pgx/v5 v5.7.2 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/lib/pq v1.10.9 // indirect
	golang.org/x/crypto v0.35.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/text v0.22.0 // indirect
)
