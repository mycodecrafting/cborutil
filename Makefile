.PHONY: build clean

build:
	@echo "Building..."
	go build -o ./bin/cborutil

clean:
	@echo "Cleaning up..."
	rm -fr bin
