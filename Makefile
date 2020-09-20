mock:
	$(GOPATH)/bin/mockgen -source=./users/db.go -destination=./tests/users/mocks/db.go -package=mocks users/db && \
    $(GOPATH)/bin/mockgen -source=./wallets/db.go -destination=./tests/wallets/mocks/db.go -package=mocks wallets/db