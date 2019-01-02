GO               ?= go
TESTTIMEOUT      := 30s
TESTFLAGS        := -cover -failfast -timeout $(TESTTIMEOUT)
HEROKU_CLI       := $(shell which heroku)
GIT_CLI          := $(shell which git)

export CGO_ENABLED := 0

.PHONY: analyse
analyse:
	@echo [gocritic] Run ...
	@$(GOPATH)/bin/gocritic check -enable='#diagnostic' ./...
	@echo [gocritic] Done.
	@$(GOPATH)/bin/gosec -fmt text ./...

.PHONY: test
test: PKG ?= $(shell go list ./... | grep -v /vendor/)
test:
	@go test $(TESTFLAGS) -tags=integration $(PKG)

.PHONY: heroku/config
heroku/config:
	@$(HEROKU_CLI) config

.PHONY: heroku/info
heroku/info:
	@$(HEROKU_CLI) apps:info

.PHONY: heroku/deploy
heroku/deploy:
	@$(GIT_CLI) push heroku master
