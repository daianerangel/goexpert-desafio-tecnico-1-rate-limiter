docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

tests:
	go test