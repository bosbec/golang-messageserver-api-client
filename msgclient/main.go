package main

import (
	"github.com/bosbec/golang-messageserver-api-client/client"
)

const (
	apiAccessId = "your access id"
	apiSecret   = "your api secret"
	apiUrl      = "https://apimessageserver.mobileresponse.se/api/v1/sms"
)

func main() {
	recipients := []string{"phoneNumber"}

	message := "your message"
	senderName := "your sender name"

	request := &client.SendSmsRequest{recipients, message, senderName}

	c := client.New(apiUrl, apiAccessId, apiSecret)

	c.SendSms(request)
}
