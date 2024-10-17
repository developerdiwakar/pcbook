# Run the main.go file
run:
	go run main.go

# Generate the go code for protobuf files  
gen: 
	protoc --proto_path=proto proto/*.proto --go_out=:pb --go-grpc_out=:pb 

# Clean the pb directory
clean:
	rm pb/*.go

# To run the test of all packages
test:
	go test -cover -race ./...

# Old Commands
# protoc --proto_path=proto proto/*.proto  --go_out=:pb --go-grpc_out=:pb --grpc-gateway_out=:pb --openapiv2_out=:swagger
# protoc --go_out=pb --go_opt=paths=source_relative \
# --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
# proto/*.proto
