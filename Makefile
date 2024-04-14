.PHONY: build run stop down delete test

build:
	docker-compose build

run:
	docker-compose up -d

stop:
	docker-compose stop

delete: stop
	docker-compose down

test:
	docker run --name testdb -e POSTGRES_PASSWORD=admin -e POSTGRES_DB=testdb -p 5432:5432 -d postgres:latest
	docker cp wait-for-postgres.sh testdb:/wait-for-postgres.sh
	docker exec testdb /bin/sh -c "chmod +x ./wait-for-postgres.sh && ./wait-for-postgres.sh localhost testdb"
	go test -v ./tests/
	docker rm -f testdb


restart: stop run