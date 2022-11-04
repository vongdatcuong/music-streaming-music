start:
	docker-compose up --build
down:
	docker-compose down -v
remove:
	docker-compose rm -fsv
prune:
	docker image prune -f
gen-protos:
	protoc --go_out=.  --go-grpc_out=. protos/**/*.proto
export_go_path:
	export GO_PATH=~/go && export PATH=$PATH:/$GO_PATH/bin