// package services

// import (
// 	"log"
// 	"net/smtp"
// )

// func SendEmail(to, subject, body string) error {
// 	from := "your-email@example.com"
// 	password := "your-email-password"

// 	// Set up SMTP server
// 	smtpHost := "smtp.example.com"
// 	smtpPort := "587"

// 	msg := "From: " + from + "\n" +
// 		"To: " + to + "\n" +
// 		"Subject: " + subject + "\n\n" +
// 		body

// 	auth := smtp.PlainAuth("", from, password, smtpHost)
// 	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
// 	if err != nil {
// 		log.Printf("Failed to send email: %v", err)
// 		return err
// 	}

// 	log.Println("Email sent successfully")
// 	return nil
// }
