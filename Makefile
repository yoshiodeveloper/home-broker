BINARY=homebroker
VERSION=1.0
BUILD=$(git rev-parse HEAD 2> /dev/null || echo "undefined")
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Build=$(BUILD)"

all:
	GO111MODULE=on go build -o $(BINARY) $(LDFLAGS)

clean:
	if [ -f $(BINARY) ] ; then rm $(BINARY) ; fi

docker:
	docker build \
		-t yoshiodeveloper/$(BINARY):latest \
		-t $(BINARY):$(VERSION) \
		--build-arg build=$(BUILD) --build-arg version=$(VERSION) \
		-f Dockerfile --no-cache .

test:
	GO111MODULE=on go test ./...

mock:
	$(GOPATH)/bin/mockgen -source=./users/db.go -destination=./tests/users/mocks/db.go -package=mocks users/db && \
    $(GOPATH)/bin/mockgen -source=./wallets/db.go -destination=./tests/wallets/mocks/db.go -package=mocks wallets/db