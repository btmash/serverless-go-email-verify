build:
	go get github.com/aws/aws-lambda-go/lambda
	go get github.com/aws/aws-lambda-go/events
	go get github.com/smancke/mailck
	env GOOS=linux go build -ldflags="-s -w" -o bin/verify verify/main.go
