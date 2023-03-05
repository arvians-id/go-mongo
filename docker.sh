docker build -t arvians/go-mongo/user -f ./user/Dockerfile .
docker build -t arvians/go-mongo/post -f ./post/Dockerfile .
docker build -t arvians/go-mongo/adapter -f ./adapter/Dockerfile .