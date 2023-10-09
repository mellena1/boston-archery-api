.PHONY: build
build:
	sam build

.PHONY: local
local: build
	sam local start-api --env-vars .env.local.json

.PHONY: build-SeasonsFunction
build-SeasonsFunction:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o $(ARTIFACTS_DIR)/bootstrap -ldflags="-s -w" handlers/seasons/*.go

.PHONY: build-AuthFunction
build-AuthFunction:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o $(ARTIFACTS_DIR)/bootstrap -ldflags="-s -w" handlers/auth/*.go
