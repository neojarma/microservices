package entity

type MailPayload struct {
	To      string `json:"to"`
	Message string `json:"message"`
}
