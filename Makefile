migrate:
	migrate -path database/postgres/migrations -database "postgres://root:root@host.docker.internal:5432/${table}?sslmode=disable" -verbose ${verbose}

table:
	migrate create -ext sql -dir database/postgres/migrations -seq ${table}

pb:
	protoc --go_out=./ --go_opt=paths=source_relative --go-grpc_out=./ --go-grpc_opt=paths=source_relative ./user/pb/*.proto
	protoc --go_out=./ --go_opt=paths=source_relative --go-grpc_out=./ --go-grpc_opt=paths=source_relative ./category/pb/*.proto
	protoc --go_out=./ --go_opt=paths=source_relative --go-grpc_out=./ --go-grpc_opt=paths=source_relative ./todo/pb/*.proto
	protoc --go_out=./ --go_opt=paths=source_relative --go-grpc_out=./ --go-grpc_opt=paths=source_relative ./gateway/pb/*.proto

build:
	docker build ./gateway -t arvians/go-todo-gateway:latest
	docker build ./worker -t arvians/go-todo-worker:latest
	docker build ./todo -t arvians/go-todo-todo:latest
	docker build ./user -t arvians/go-todo-user:latest
	docker build ./category_todo -t arvians/go-todo-category-todo:latest
	docker build ./category -t arvians/go-todo-category:latest

# If you want to run the test, do not run with docker-compose.yml
# instead you have to run docker-compose.dev.yml and run all the services manually
test:
	ginkgo gateway/tests/integration