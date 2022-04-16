module github.com/davepgreene/go-db-credential-refresh

go 1.15

require (
	bou.ke/monkey v1.0.2
	github.com/aws/aws-sdk-go-v2 v1.16.2
	github.com/aws/aws-sdk-go-v2/credentials v1.11.2
	github.com/aws/aws-sdk-go-v2/feature/rds/auth v1.1.19
	github.com/containerd/containerd v1.5.10 // indirect
	github.com/go-sql-driver/mysql v1.6.0
	github.com/go-test/deep v1.0.8
	github.com/hashicorp/go-hclog v1.2.0
	github.com/hashicorp/vault v1.10.0
	github.com/hashicorp/vault-plugin-auth-kubernetes v0.12.0
	github.com/hashicorp/vault/api v1.5.0
	github.com/hashicorp/vault/sdk v0.4.2-0.20220321211954-d7083ad326db
	github.com/jackc/pgx/v4 v4.15.0
	github.com/lib/pq v1.10.5
	github.com/mitchellh/mapstructure v1.4.3
	k8s.io/api v0.23.5
)
