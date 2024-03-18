module github.com/davepgreene/go-db-credential-refresh/store/awsrds

go 1.18

replace github.com/davepgreene/go-db-credential-refresh => ../../

require (
	bou.ke/monkey v1.0.2
	github.com/aws/aws-sdk-go-v2 v1.26.0
	github.com/aws/aws-sdk-go-v2/credentials v1.17.7
	github.com/aws/aws-sdk-go-v2/feature/rds/auth v1.4.3
	github.com/davepgreene/go-db-credential-refresh v1.0.0
	github.com/mitchellh/mapstructure v1.5.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/aws/smithy-go v1.20.1 // indirect
	github.com/go-sql-driver/mysql v1.8.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.3 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d80cb9 // indirect
	github.com/jackc/pgtype v1.14.2 // indirect
	github.com/jackc/pgx/v4 v4.18.3 // indirect
	github.com/lib/pq v1.10.9 // indirect
	golang.org/x/crypto v0.21.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)
