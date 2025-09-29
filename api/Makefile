# Define the default target
.PHONY: default
default: help

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  seeder FILENAME     Process the specified seeder file name"

# application repository and binary file name
NAME=club

# application repository path
REPOSITORY=github.com/grundigdev/${NAME}


install:
	go mod download

run-dev:
	echo "Starting Application In Development Mode"
	go run ./cmd/api

reset:
	echo "Running migrations down"
	go run ./database/migrate_down.go
	echo "Running migrations up"
	go run ./database/migrate_up.go
	echo "Starting Application In Development Mode"
	go run ./cmd/api

seed:
	echo "Running migrations down"
	go run ./database/migrate_down.go
	echo "Running migrations up"
	go run ./database/migrate_up.go
	echo "Seedinng : $(FILENAME).go"
	go run "./database/seeders/absence_seeder.go"


#E.G make seeder FILENAME=category
.PHONY: seeder
seeder:
ifdef FILENAME
	echo "Seedinng : $(FILENAME).go"
	go run "./database/seeders/$(FILENAME)_seeder.go"
else
	echo "Error: FILENAME is not specified. Please provide the filename using 'make seeder FILENAME=<filename>'"
	exit 1
endif

migrate_fresh:
	echo "Running fresh migrations"
	go run ./database/migrate_fresh.go

migrate_down:
	go run ./database/drop_migrations.go

test:
	 gotest ./tests/... -v