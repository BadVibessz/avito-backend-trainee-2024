# testing vars
export TEST_CONTAINER_NAME=test_db
export TEST_DBSTRING=postgresql://postgres:postgres@localhost:5433/test?sslmode=disable
export TEST_GOOSE_DRIVER=postgres
export TEST_JWT_SECRET=test_secret
export TEST_DOCKER_PORT=5433

export DOCKER_IMAGE_NAME=avito_trainee_image


docker.up:
	docker run --rm -d -p 5432:5432 --name avito_trainee_postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=avito-trainee postgres
	docker compose up -d


test.integration:
	docker run --rm -d -p $$TEST_DOCKER_PORT:5432 --name $$TEST_CONTAINER_NAME -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=test postgres

	sleep 2 # wait for postgres to run in docker container, todo: bad practice, use go-migrate instead?

	# [command] || true prevents the script to stop even if error occurred executing command, so newly created docker container will be deleted anyway

	goose -dir ./db/migrations $$TEST_GOOSE_DRIVER $$TEST_DBSTRING up || true # apply migrations
	go test -v ./tests/ || true

	docker stop $$TEST_CONTAINER_NAME