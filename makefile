start:
	docker-compose up

q2c:
	sqlc generate

test:
	go test -v -covermode=atomic -race -coverpkg=./... ./... \
	-coverprofile coverage.out.tmp && cat coverage.out.tmp | grep -v "_mock.go" | grep -v "injector.go" | grep -v "_gen.go" > coverage.out && rm coverage.out.tmp && \
	go tool cover -func coverage.out
	
generateInjector:
	wire ./...

migrateInit:
	migrate create -ext sql -dir db/migration -seq $(name)

mockRepo:
	mockgen -package mockrepo -destination db/repository/mock/repository_mock.go -source=db/repository/repository.go -aux_files github.com/gadhittana01/cosmos-validation-tracking/db/repository=db/repository/querier.go

mockConnect:
	mockgen -package mockrepo -destination db/repository/mock/connect_mock.go -source=./utils/connect.go

mockAuthMiddleware:
	mockgen -package mockutl -source=./utils/auth_middleware.go -destination=./utils/mock/auth_middleware_mock.go

mockValidatorSvc:
	mockgen -package mocksvc -source=./service/validator_service.go -destination=./service/mock/validator_service_mock.go

mockLogger:
	mockgen -package mockutl -source=./utils/logger.go -destination=./utils/mock/logger_mock.go

mockHTTPClient:
	mockgen -package mockutl -source=./utils/http_client.go -destination=./utils/mock/http_client_mock.go

checkLint:
	golangci-lint run ./... -v

fixLint:
	golangci-lint run --fix

swaggerLocal:
	swag init && rm -rf docs/swagger.json