package main

import (
	"fmt"
	"encoding/json"
	"github.com/smancke/mailck"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type emailStatus struct {
	Email string `json:"email"`
	Status string `json:"status"`
	Message string `json:"message"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	result, _ := mailck.Check("btmash@gmail.com", request.Body)
	returnData := new(emailStatus)
	returnData.Email = "btmash@gmail.com"
	switch {

	  case result.IsValid():
			returnData.Status = "valid"
			returnData.Message = "Perfectly valid email"

	  case result.IsError():
			returnData.Status = "error"
			returnData.Message = "SMTP Failed. May or may not be valid email"

	  case result.IsInvalid():
			returnData.Status = "invalid"
		  switch (result) {
		    case mailck.InvalidDomain:
					returnData.Message = "Invalid Domain provided"
		    case mailck.InvalidSyntax:
					returnData.Message = "Invalid email syntax"
				case mailck.MailboxUnavailable:
					returnData.Message = "Email address does not exist"
				case mailck.Disposable:
					returnData.Message = "Disposable email address. You may choose to test further and see if this is fine"
		  }
	}
	jsonResult, _ := json.Marshal(returnData)
	fmt.Println(string(jsonResult))
  return events.APIGatewayProxyResponse{Body: string(jsonResult), StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
