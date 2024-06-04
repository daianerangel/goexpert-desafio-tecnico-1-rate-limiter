docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

delete-redis:
	chmod +x delete-redis.sh
	./delete-redis.sh

tests: delete-redis
	go test