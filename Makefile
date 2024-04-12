# Testing Vars
export TEST_CONTAINER_NAME=test_db
export TEST_DBSTRING=postgresql://postgres:postgres@localhost:5433/test?sslmode=disable
export TEST_GOOSE_DRIVER=postgres
export TEST_JWT_SECRET=test_secret


test.integration:
	docker run --rm -d -p 5433:5432 --name $$TEST_CONTAINER_NAME -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=test postgres


	sleep 10 # todo: bad practice, use go-migrate instead?

	goose -dir ./db/migrations postgres 'postgresql://postgres:postgres@localhost:5433/test?sslmode=disable' up # apply migrations
	go test -v ./tests/

	docker stop $$TEST_CONTAINER_NAME