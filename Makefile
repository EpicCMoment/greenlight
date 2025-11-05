.SILENT: clean

all: format build run

build:
	go build -o ./bin/api ./cmd/api

run:
	./bin/api

format:
	@echo "Converting all go files to standard format..."
	go fmt ./...

clean:
	rm ./bin/*