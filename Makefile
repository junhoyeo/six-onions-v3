# Set an output name for the application
OUTFILE := six-onions-v3-app

# List the Go source files
SOURCES := $(wildcard six-onions-v3/*.go)

# Build the application
all: $(OUTFILE)

$(OUTFILE): $(SOURCES)
	@go build -o $(OUTFILE) ./six-onions-v3

# Clean up
clean:
	@rm -f $(OUTFILE)

# Run the application
run: $(OUTFILE)
	@./$(OUTFILE)

.PHONY: all clean run
