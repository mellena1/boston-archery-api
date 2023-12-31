.PHONY: build
build:
	sam build

.PHONY: local
local: build
	docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
	docker-compose up -d
	sam local start-api --env-vars .env.local.json --docker-network boston-archery

.PHONY: deploy
deploy: build
	sam deploy

.PHONY: build-SeasonsFunction
build-SeasonsFunction:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o $(ARTIFACTS_DIR)/bootstrap -ldflags="-s -w" handlers/seasons/*.go

.PHONY: build-AuthFunction
build-AuthFunction:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o $(ARTIFACTS_DIR)/bootstrap -ldflags="-s -w" handlers/auth/*.go

.PHONY: swagger
swagger:
	swagger generate spec -m -o docs/swagger.json

.PHONY: swagger-serve
swagger-serve: swagger
	swagger serve -F swagger docs/swagger.json

.PHONY: test
test:
	go test ./...
