package main

import (
	"fmt"
	"encoding/json"
	"time"

	"github.com/smancke/mailck"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type emailStatus struct {
	Email            string `json:"email"`
	ValidSyntax      bool   `json:"syntax_valid"`
	DisposableStatus bool   `json:"is_disposable_address"`
	SmtpStatus       string `json:"smtp_valid"`
	Message          string `json:"message"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ch := make(chan *emailStatus, 1)
	returnData := new(emailStatus)

	// Create a goroutine that will attempt to connect to the SMTP server
	go func() {
		// You may need to change the email address to something you would use.
		result, _ := mailck.Check("test@example.com", request.Body)
		returnData := new(emailStatus)
		returnData.Email = request.Body
		switch {

		case result.IsValid():
			returnData.ValidSyntax = true
			returnData.DisposableStatus = false
			returnData.SmtpStatus = "true"
			returnData.Message = "Perfectly valid email"

		case result.IsError():
			returnData.ValidSyntax = true
			returnData.DisposableStatus = false
			returnData.SmtpStatus = "error"
			returnData.Message = "SMTP Failed. May or may not be valid email"

		case result.IsInvalid():
			returnData.SmtpStatus = "false"
			switch result {
			case mailck.InvalidDomain:
				returnData.ValidSyntax = false
				returnData.DisposableStatus = false
				returnData.Message = "Invalid Domain provided"
			case mailck.InvalidSyntax:
				returnData.ValidSyntax = false
				returnData.DisposableStatus = false
				returnData.Message = "Invalid email syntax"
			case mailck.MailboxUnavailable:
				returnData.ValidSyntax = true
				returnData.DisposableStatus = false
				returnData.Message = "Email address does not exist"
			case mailck.Disposable:
				returnData.ValidSyntax = true
				returnData.DisposableStatus = true
				returnData.Message = "Disposable email address."
			}
		}
		ch <- returnData
	}()

	select {
	case validate := <-ch:
		returnData = validate
	case <-time.After(10 * time.Second):
		returnData.Email = request.Body
		returnData.ValidSyntax = true
		returnData.DisposableStatus = true
		returnData.SmtpStatus = "error"
		returnData.Message = "SMTP Dial timed out. May or may not be valid email"
	}

	jsonResult, _ := json.Marshal(returnData)
	fmt.Println(string(jsonResult))
	return events.APIGatewayProxyResponse{Body: string(jsonResult), StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
