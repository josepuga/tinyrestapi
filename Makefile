# José Puga 2025. GPL3 License
#
BINARY_PATH=./app/tinyrestapi
VERSION=$(shell git describe --tags --always)
FLAGS="-w -s -X 'main.version=$(VERSION)'"

build:
# CGO_ENABLED= 0 for static compilation. 1 for dynamic.
	CGO_ENABLED=0 GOOS=linux go build -ldflags=$(FLAGS) -o $(BINARY_PATH) .
	CGO_ENABLED=0 GOOS=windows go build -ldflags=$(FLAGS) -o $(BINARY_PATH).exe .

run: build
	$(BINARY_PATH)

clean:
	go clean
	rm -f  $(BINARY_PATH)
	rm -f  $(BINARY_PATH).exe
