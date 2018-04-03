# Email verification via serverless

## Install

1. You first need to install [Go](https://golang.org/doc/install)
2. You need to install [Serverless](https://serverless.com/framework/docs/getting-started/)
3. You need [AWS Credentials](https://serverless.com/framework/docs/providers/aws/guide/credentials/)

## Set up task

1. Run `make`
2. Run `sls deploy`
3. You should receive a url showing where the application got deployed

## Testing the lambda

1. Open `test.sh`
2. Update LAMBDA_URL with the url received from the serverless steps above
3. Update the email address to test with your choice

### About this code

The brunt of the code is in `verify/main.go`. The code takes an email address
from the body of the POST request and checks to see if it is a valid email
address.
