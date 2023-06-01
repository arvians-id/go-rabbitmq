migrate:
	migrate -path database/postgres/migrations -database "postgres://root:root@host.docker.internal:5432/${table}?sslmode=disable" -verbose ${verbose}

table:
	migrate create -ext sql -dir database/postgres/migrations -seq ${table}

pb:
	protoc \
 	--go_out=services/user/pb --go-grpc_out=services/user/pb \
 	--go_out=services/category/pb --go-grpc_out=services/category/pb \
 	--go_out=services/todo/pb --go-grpc_out=services/todo/pb \
 	--go_out=gateway/pb --go-grpc_out=gateway/pb \
	--go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --proto_path=proto \
	./proto/*.proto

gql:
	cd gateway/api/gql/ && go run github.com/99designs/gqlgen@v0.17.31 generate

dataloaden:
	cd gateway/api/gql/model && \
	go run github.com/vektah/dataloaden UserLoader int64 *github.com/arvians-id/go-rabbitmq/gateway/api/gql/model.User && \
	go run github.com/vektah/dataloaden TodoLoader int64 []*github.com/arvians-id/go-rabbitmq/gateway/api/gql/model.Todo && \
	go run github.com/vektah/dataloaden CategoryLoader int64 []*github.com/arvians-id/go-rabbitmq/gateway/api/gql/model.Category

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