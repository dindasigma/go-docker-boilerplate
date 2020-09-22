run: stop up

up:
	docker-compose -f docker-compose.yml up --build

stop:
	docker-compose -f docker-compose.yml stop

down:
	docker-compose -f docker-compose.yml down

test:
	docker-compose -f docker-compose.test.yml down --volumes
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit

test-down:
	docker-compose -f docker-compose.test.yml down --volumes

build-dev:
	swag init -g commands/runserver.go
	go build -v commands/runserver.go