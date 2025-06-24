package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
)

type FormData struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type Config struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	RecipientEmail string
	ServerPort   string
	APIKey       string
}

func loadConfig() (*Config, error) {
	config := &Config{
		SMTPHost:       os.Getenv("SMTP_HOST"),
		SMTPPort:       os.Getenv("SMTP_PORT"),
		SMTPUsername:   os.Getenv("SMTP_USERNAME"),
		SMTPPassword:   os.Getenv("SMTP_PASSWORD"),
		RecipientEmail: os.Getenv("RECIPIENT_EMAIL"),
		ServerPort:     os.Getenv("SERVER_PORT"),
		APIKey:         os.Getenv("API_KEY"),
	}

	if config.SMTPHost == "" {
		return nil, fmt.Errorf("SMTP_HOST environment variable is required")
	}
	if config.SMTPPort == "" {
		config.SMTPPort = "587"
	}
	if config.SMTPUsername == "" {
		return nil, fmt.Errorf("SMTP_USERNAME environment variable is required")
	}
	if config.SMTPPassword == "" {
		return nil, fmt.Errorf("SMTP_PASSWORD environment variable is required")
	}
	if config.RecipientEmail == "" {
		return nil, fmt.Errorf("RECIPIENT_EMAIL environment variable is required")
	}
	if config.ServerPort == "" {
		config.ServerPort = "8080"
	}
	if config.APIKey == "" {
		return nil, fmt.Errorf("API_KEY environment variable is required")
	}

	return config, nil
}

func sendEmail(config *Config, formData FormData) error {
	auth := smtp.PlainAuth("", config.SMTPUsername, config.SMTPPassword, config.SMTPHost)

	to := []string{config.RecipientEmail}
	
	subject := formData.Subject
	if subject == "" {
		subject = "New Form Submission"
	}

	headers := make(map[string]string)
	headers["From"] = config.SMTPUsername
	headers["To"] = config.RecipientEmail
	headers["Subject"] = subject
	headers["Reply-To"] = formData.Email

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n"
	message += fmt.Sprintf("Name: %s\r\n", formData.Name)
	message += fmt.Sprintf("Email: %s\r\n", formData.Email)
	message += fmt.Sprintf("\r\n%s", formData.Message)

	err := smtp.SendMail(
		config.SMTPHost+":"+config.SMTPPort,
		auth,
		config.SMTPUsername,
		to,
		[]byte(message),
	)

	return err
}

func handleFormSubmission(config *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-API-Key")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		apiKey := r.Header.Get("X-API-Key")
		if apiKey != config.APIKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var formData FormData
		err := json.NewDecoder(r.Body).Decode(&formData)
		if err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		if formData.Email == "" || formData.Message == "" {
			http.Error(w, "Email and message are required fields", http.StatusBadRequest)
			return
		}

		err = sendEmail(config, formData)
		if err != nil {
			log.Printf("Failed to send email: %v", err)
			http.Error(w, "Failed to send email", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "success",
			"message": "Email sent successfully",
		})
	}
}

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	http.HandleFunc("/send-email", handleFormSubmission(config))
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("Server starting on port %s", config.ServerPort)
	if err := http.ListenAndServe(":"+config.ServerPort, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}