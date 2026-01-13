package services

import (
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
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

func SendEmail(toEmail, otp, reason string) error {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("SENDGRID_API_KEY not set")
	}

	from := mail.NewEmail("Chat app", "harsh.santoshi07@gmail.com")

	to := mail.NewEmail("", toEmail)

	subject := "Your OTP for Chat app"

	plainText := fmt.Sprintf(
		"Your OTP for %s is %s. This code is valid for 5 minutes. If you did not request this, please ignore this email.",
		reason,
		otp,
	)

	htmlBody := fmt.Sprintf(`
	<html>
		<body style="font-family: Arial, sans-serif; text-align: center;">
			<h2>Expense Tracker</h2>
			<p>Your OTP for <strong>%s</strong> is:</p>
			<h1 style="color: #4A90E2; letter-spacing: 5px;">%s</h1>
			<p>This code is valid for 5 minutes.</p>
			<p>If you did not request this, please ignore this email.</p>
		</body>
	</html>
	`, reason, otp)

	message := mail.NewSingleEmail(
		from,
		subject,
		to,
		plainText,
		htmlBody,
	)

	client := sendgrid.NewSendClient(apiKey)
	_, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
