make-migrate-up:
	migrate -path database/postgres/migrations -database "postgres://postgres:root@localhost:5432/go_rabbitmq?sslmode=disable" -verbose up

make-migrate-down:
	migrate -path database/postgres/migrations -database "postgres://postgres:root@localhost:5432/go_rabbitmq?sslmode=disable" -verbose down

make-migrate-table:
	migrate create -ext sql -dir database/postgres/migrations -seq ${table}
make-pb:
	protoc --go_out=./ --go_opt=paths=source_relative --go-grpc_out=./ --go-grpc_opt=paths=source_relative ./user/pb/*.proto
	protoc-go-inject-tag -input=./user/pb/*.pb.go
make-pb-todo:
	protoc --go_out=./todo/pb/ --go_opt=paths=source_relative --go-grpc_out=./todo/pb --go-grpc_opt=paths=source_relative ./todo/pb/*.proto ./user/pb/user.proto
	protoc-go-inject-tag -input=./todo/pb/*.pb.go
