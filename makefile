start:
	docker-compose up --build
down:
	docker-compose down -v
remove:
	docker-compose rm -fsv
prune:
	docker image prune -f