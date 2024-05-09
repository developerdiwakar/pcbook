# Run the main.go file
run:
	go run main.go

# Generate the go code for protobuf files  
gen: 
	protoc --proto_path=proto --go_out=pb --go-grpc_out=pb proto/*.proto

# Clean the pb directory
clean:
	rm pb/*.go
