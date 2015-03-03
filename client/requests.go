package client

type SendSmsRequest struct {
	Recipients []string `json:"recipients"`
	Body       string   `json:"body"`
	SenderName string   `json:"senderName"`
}
