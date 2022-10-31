start:
	docker-compose up --build && docker-compose rm -f
down:
	docker-compose down -v
