##Setup
Start project with docker desktop
Init the Yaml file with
docker run --rm -v "%cd%:/src" -w /src kjconroy/sqlc init

Generate SQLC with...yaml file needs config from SQLC 1.4 tag setup
docker run --rm -v "%cd%:/src" -w /src kjconroy/sqlc generate

Can also run jazz from the MakeFile check that for terminal commands and prefix with `make`

##SQLC Reference
https://github.com/kyleconroy/sqlc/tree/v1.4.0

Note: SQLC written in `account.sql` generate to account.sql.go via CLI operations with the `make` file