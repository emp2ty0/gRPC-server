include .env
export

migrations-up:
	go run cmd/migrator/main.go --storage-path=./db/sso.db --migrations-path=./migrations

run-app:
	go run cmd/sso/main.go