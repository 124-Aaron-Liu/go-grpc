protoc calculator.proto --go_out=plugins=grpc:. --go_opt=paths=source_relative --go-grpc_out=. calculator.proto

protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/calculator/calculator.proto