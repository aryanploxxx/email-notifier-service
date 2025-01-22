package main

import (
	"bytes"
	"email-notifier-service/models"
	"encoding/json"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func sendHtmlEmail(to []string, subject string, htmlBody string) error {
	auth := smtp.PlainAuth(
		"",
		os.Getenv("FROM_EMAIL"),
		os.Getenv("FROM_EMAIL_PASSWORD"),
		os.Getenv("FROM_EMAIL_SMTP"),
	)

	// These headers are required to be added in the message to make SMTP understand that the body contains HTML contain
	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"

	message := "Subject: " + subject + "\n" + headers + "\n\n" + htmlBody
	return smtp.SendMail(
		os.Getenv("SMTP_ADDR"),
		auth,
		os.Getenv("FROM_EMAIL"),
		to,
		[]byte(message),
	)
}

func HTMLTemplateEmailHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the JSON request body
	var reqBody models.EmailWithTemplateRequestBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Convert Param3 (comma-separated string) to a slice of strings
	to := strings.Split(reqBody.ToAddr, ",")

	// Parse the HTML template
	tmpl, err := template.ParseFiles("./templates/" + reqBody.Template + ".html")
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	// Render the template with the map data
	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, reqBody.Vars); err != nil {
		log.Fatalf("Failed to execute template: %v", err)
	}

	log.Println(rendered.String())

	err = sendHtmlEmail(to, reqBody.Subject, rendered.String())
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email sent successfully"))
}

// Make POST request at - http://localhost:9001/html_email_template

func main() {
	godotenv.Load()

	router := mux.NewRouter()
	router.HandleFunc("/html_email_template", HTMLTemplateEmailHandler)

	log.Printf("Server is listening at 9001")
	log.Fatal(http.ListenAndServe(":9001", router))
}

/*
	-> Text in Body
	{
		"to_addr": "test@gmail.com,mail@gmail.com",
		"subject": "HTML Mail Test",
		"body": "<h1>Woow</h1><p>This is a paragraph.</p>"
	}
*/

/*
	-> HTML in Body
	{
		"to_addr": "test@gmail.com,mail@gmail.com",
		"subject": "HTML Mail Test",
		"body": "<h1>woow</h1><p>this is a paragraph</p>"
	}
*/

/*
	-> HTML with Templates
	-> Email Clients are unable to resolve the relative path for images
	{
		"to_addr": "test@gmail.com",
		"subject": "HTML Mail Test",
		"template": "helloEmail",
		"vars": {
			"Name": "Aryan",
			"Image": "https://media.licdn.com/dms/image/v2/D5616AQFrjv-6DGEv-w/profile-displaybackgroundimage-shrink_350_1400/profile-displaybackgroundimage-shrink_350_1400/0/1721797487032?e=1743033600&v=beta&t=lmphfMyK457Yluc2zJtf-nY0OEJf8wG8VwygUrVLCJ0"
		}
	}
*/
