package models

// Define a struct to hold the incoming JSON request body
type EmailRequestBody struct {
	ToAddr  string `json:"to_addr"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type EmailWithTemplateRequestBody struct {
	ToAddr   string            `json:"to_addr"`
	Subject  string            `json:"subject"`
	Template string            `json:"template"`
	Vars     map[string]string `json:"vars"`
}
