migrate-test:
	migrate -path database/postgres/migrations -database "postgres://root:root@host.docker.internal:5432/go_rabbitmq_test?sslmode=disable" -verbose ${verbose}

migrate-prod:
	migrate -path database/postgres/migrations -database "postgres://root:root@host.docker.internal:5432/go_rabbitmq?sslmode=disable" -verbose ${verbose}

create-table:
	migrate create -ext sql -dir database/postgres/migrations -seq ${table}

generate-pb:
	protoc --go_out=./ --go_opt=paths=source_relative --go-grpc_out=./ --go-grpc_opt=paths=source_relative ./user/pb/*.proto
	protoc --go_out=./ --go_opt=paths=source_relative --go-grpc_out=./ --go-grpc_opt=paths=source_relative ./category/pb/*.proto
	protoc --go_out=./ --go_opt=paths=source_relative --go-grpc_out=./ --go-grpc_opt=paths=source_relative ./todo/pb/*.proto
	protoc --go_out=./ --go_opt=paths=source_relative --go-grpc_out=./ --go-grpc_opt=paths=source_relative ./gateway/pb/*.proto

build:
	docker build ./gateway -t arvians/go-todo-gateway:latest
	docker build ./worker -t arvians/go-todo-worker:latest
	docker build ./todo -t arvians/go-todo-todo:latest
	docker build ./user -t arvians/go-todo-user:latest

test:
	ginkgo gateway/tests/integration