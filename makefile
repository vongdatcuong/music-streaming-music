start:
	docker-compose up --build
start-service:
	docker-compose up --build service
start-db:
	docker-compose up db
down:
	docker-compose down -v
remove:
	docker-compose rm -fsv
prune:
	docker image prune -f
gen-protos:
	protoc --go_out=. --grpc-gateway_out=. --go-grpc_out=. protos/**/*.proto 
export_go_path:
	export GO_PATH=~/go && export PATH=$PATH:/$GO_PATH/bin
connect_network:
	docker network connect api-gateway_fullstack music-streaming-music-service && docker inspect -f '{{ $network := index .NetworkSettings.Networks "api-gateway_fullstack" }}{{ $network.IPAddress}}' music-streaming-music-service