.PHONY: run build clean

run:
	go run main.go

build:
	mkdir -p build
	@echo "Building the synchronizer..."
	go build -o build/synchronizer main.go

clean:
	rm -f build