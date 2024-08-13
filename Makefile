.PHONY: build
build:
	sam build

.PHONY: local
local: build
	docker run --privileged --rm tonistiigi/binfmt --install all
	docker-compose up -d
	sam local start-api --env-vars .env.local.json --docker-network boston-archery

.PHONY: deploy
deploy: build
	sam deploy

.PHONY: build-ApiFunction
build-ApiFunction:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o $(ARTIFACTS_DIR)/bootstrap -ldflags="-s -w" cmd/main.go

.PHONY: swagger
swagger:
	swagger generate spec -m -o docs/swagger.json

.PHONY: swagger-serve
swagger-serve: swagger
	swagger serve -F swagger docs/swagger.json

.PHONY: test
test:
	go test -cover -tags=integration ./...

.PHONY: test-in-ci
test-in-ci:
	TEST_IN_CI=true make test
