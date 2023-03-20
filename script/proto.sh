protoc proto/*.proto --go_out=plugins=grpc:./ --go_out=plugins=grpc:./adapter --go_out=plugins=grpc:./post --go_out=plugins=grpc:./user

protoc-go-inject-tag -input=./*/pb/*.pb.go