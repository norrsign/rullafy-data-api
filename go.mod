module github.com/norrsign/rullafy-data-api

go 1.24.5

require (
	github.com/jackc/pgx/v5 v5.7.5
	github.com/sirupsen/logrus v1.9.3
	github.com/vanern/goapi v1.0.0
)

require (
	github.com/golang-jwt/jwt/v5 v5.2.2 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/spf13/cobra v1.9.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/sync v0.13.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
)

replace (

	github.com/vanern/goapi => ../../goapi
)
