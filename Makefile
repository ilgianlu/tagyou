# Define the name of the Go application
BINARY_NAME = tagyou

# Define the source file
SOURCE_FILE = main.go

# Build target
build:
	cd sqlc && sqlc generate
	go test ./...
	go build -o $(BINARY_NAME) $(SOURCE_FILE)

# Clean target
clean:
	rm -f $(BINARY_NAME)
	rm -r sqlc/dbaccess
	find . -name *.db3 -type f -delete