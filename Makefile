.SILENT: clean

all:
	go build -o ./bin/api ./cmd/api

run:
	./bin/api

clean:
	rm ./bin/*