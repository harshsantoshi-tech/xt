package services

import (
	"fmt"
	"github.com/resend/resend-go/v2"
	"log"
	"os"
)

// SendEmail sends a 6-digit OTP to the specified recipient
//func SendEmail(toEmail string, otp string, reason string) error {
//	// 1. Get credentials from environment variables (configs)
//	from := os.Getenv("SMTP_EMAIL")
//	password := os.Getenv("SMTP_PASS")
//	smtpHost := os.Getenv("SMTP_HOST")
//	smtpPort := os.Getenv("SMTP_PORT")
//
//	log.Info("SMTP_HOST:", smtpHost)
//	// 2. Compose the email headers and body
//	// Note: The empty line between headers and body is mandatory in SMTP
//	subject := "Subject: Expense Tracker Verification Code\n"
//	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
//	body := fmt.Sprintf(`
//        <html>
//            <body>
//                <h2>Welcome to Expense Tracker!</h2>
//                <p>Please use the following One-Time Password (OTP) %s</p>
//                <h1 style="color: #4A90E2; letter-spacing: 5px;">%s</h1>
//                <p>This code is valid for 5 minutes. If you did not request this, please ignore this email.</p>
//            </body>
//        </html>`, reason, otp)
//
//	message := []byte(subject + mime + body)
//
//	// 3. Authenticate with the SMTP server
//	auth := smtp.PlainAuth("", from, password, smtpHost)
//
//	// 4. Send the email
//	address := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
//	log.Info("Sending Email", toEmail)
//	err := smtp.SendMail(address, auth, from, []string{toEmail}, message)
//	if err != nil {
//		log.Error("Failed to send email", toEmail)
//		return fmt.Errorf("failed to send email: %w", err)
//	}
//
//	return nil
//}

func SendEmail(toEmail string, otp string, reason string) error {

	apiKey := os.Getenv("SENDGRID_API_KEY")
	client := resend.NewClient(apiKey)
	htmlBody := fmt.Sprintf(`
       <html>
           <body style="font-family: Arial, sans-serif; text-align: center;">
               <h2>Welcome to Expense Tracker!</h2>
               <p>Please use the following One-Time Password (OTP) %s</p>
               <h1 style="color: #4A90E2; letter-spacing: 5px;">%s</h1>
               <p>This code is valid for 5 minutes. If you did not request this, please ignore this email.</p>
           </body>
       </html>`, reason, otp)
	params := &resend.SendEmailRequest{
		From:    "onboarding@resend.dev",
		To:      []string{toEmail},
		Subject: "Hello World",
		Html:    htmlBody,
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		return err
	}
	log.Println(sent)
	return nil
}
