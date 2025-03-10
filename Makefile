include .env

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/api -db-dsn=${OMDB_API_DB_DSN_DEV} 

## live/server: run air
.PHONY: live/server
live/server:
	air

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${OMDB_API_DB_DSN_DEV}

## db/create-dev-database: create local database
.PHONY: db/create-dev-database
db/create-dev-database:
	@echo 'Creating dev database...'
	sudo -i -u postgres psql -c "CREATE DATABASE ${DB_NAME}"
	sudo -i -u postgres psql -d ${DB_NAME} -c "CREATE EXTENSION IF NOT EXISTS citext"
	sudo -i -u postgres psql -d ${DB_NAME} -c "CREATE ROLE ${DB_NAME} WITH LOGIN PASSWORD '${DB_PASSWORD_DEV}'"
	sudo -i -u postgres psql -d ${DB_NAME} -c "GRANT ALL ON SCHEMA public TO ${DB_NAME}"

## db/data-download: downloads new OMDB CSV files
.PHONY: db/data-download
db/data-download:
	@echo 'Downloading OMDB CSVs...'
	@./sql/data-import/download.sh
	@echo 'OMDB data downloaded and unziped...'


## db/import-data: imports OMDB dataset
.PHONY: db/import-data
db/import-data:
	@echo 'Importing OMDB data...'
	@psql -v ON_ERROR_STOP=1 -d "${OMDB_API_DB_DSN_DEV}" -c "\i sql/data-import/run.sql"
	@echo 'Import done...'

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	@cd sql/migrations && goose postgres ${OMDB_API_DB_DSN_DEV} up && cd ../..

## db/pg_dump/schema: slqc generates go types from database schema and queries
.PHONY: db/pg_dump/schema
db/pg_dump/schema:
	pg_dump "${OMDB_API_DB_DSN_DEV}" --schema-only > sql/schema/omdb_api_schema.sql

## db/sqlc/generate: slqc generates go types from database schema and queries
.PHONY: db/sqlc/generate
db/sqlc/generate:
	sqlc generate

## db/init: initialize a docker postress container
.PHONY: db/init
db/init:
	@docker run -e POSTGRES_PASSWORD=${DOCKER_POSTGRES_PW} --name=${DOCKER_POSTGRES_CONTAINER_NAME} --rm -d -p 5432:5432 postgres && sleep 3
	@docker exec -u postgres -it pg-omdb-api psql -c "CREATE DATABASE omdb_api;"


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy and vendor dependencies and format, vet and test all code
.PHONY: audit
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

# ==================================================================================== #
# BUILD
# ==================================================================================== #

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -ldflags='-s' -o=./bin/omdb-api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/omdb-api ./cmd/api


# ==================================================================================== #
# PRODUCTION
# ==================================================================================== #

## production/connect: connect to the production server
.PHONY: production/connect
production/connect:
	ssh omdb-api@${PRODUCTION_HOST_IP}

## production/deploy/api: deploy the api to production
.PHONY: production/deploy/app
production/deploy/app:
	ssh -t omdb-api@${PRODUCTION_HOST_IP} 'mkdir -p api/sql/{schema,migrations}'

	rsync -P ./bin/linux_amd64/omdb-api omdb-api@${PRODUCTION_HOST_IP}:~/api
	rsync -rP --delete ./sql/schema/ omdb-api@${PRODUCTION_HOST_IP}:~/api/sql/schema
	
	rsync -rP --delete ./sql/migrations/ omdb-api@${PRODUCTION_HOST_IP}:~/api/sql/migrations
	ssh -t omdb-api@${PRODUCTION_HOST_IP} 'cd api/sql/migrations && goose postgres ${OMDB_API_DB_DSN_PROD} up'
	
	rsync -P ./remote/production/omdb-api.service omdb-api@${PRODUCTION_HOST_IP}:~/api
	ssh -t omdb-api@${PRODUCTION_HOST_IP} '\
		sudo mv ~/api/omdb-api.service /etc/systemd/system/ \
		&& sudo systemctl enable omdb-api \
		&& sudo systemctl restart omdb-api \
		'
	@echo "deployment complete..."

## production/deploy/initial-setup: initial setup of api. Next: production/import-data/transfer
.PHONY: production/deploy/initial-setup
production/deploy/initial-setup:
	rsync -P ./remote/setup/01-initial-setup.sh omdb-api@${PRODUCTION_HOST_IP}:~/
	ssh -t omdb-api@${PRODUCTION_HOST_IP} 'sudo ./01-initial-setup.sh'
	@echo "Initial Setup complete..."

## production/import-data/transfer: transfer data to prod. Next: production/import-data/run
.PHONY: production/import-data/transfer
production/import-data/transfer:
	ssh -t omdb-api@${PRODUCTION_HOST_IP} 'mkdir -p sql/data-import/'
	rsync -rP --delete ./sql/data-import/ omdb-api@${PRODUCTION_HOST_IP}:~/sql/data-import

## production/import-data/run: run import to prod database. Next: production/deploy/app
.PHONY: production/import-data/run
production/import-data/run:
	ssh -t omdb-api@${PRODUCTION_HOST_IP} '\
		psql -v ON_ERROR_STOP=1 -d "${OMDB_API_DB_DSN_PROD}" -c "\i sql/data-import/run.sql"\
		'
	@echo "data import complete..."
 