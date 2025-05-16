# Define the name of the Go application
BINARY_NAME = tagyou

# Define the source file
SOURCE_FILE = main.go

# Build target
build:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	cd sqlc && sqlc generate
	go mod tidy
	go test ./...
	go build -o $(BINARY_NAME) $(SOURCE_FILE)

# Clean target
clean:
	rm -f $(BINARY_NAME)
	rm -r sqlc/dbaccess
	find . -name *.db3 -type f -delete

# Init for tests
init:
	find . -name *.db3 -type f -delete
	find . -name *.csv -type f -delete
	INIT_DB=true \
	INIT_ADMIN_PASSWORD=password \
	go run main.go
