package services

import (
	"fmt"
	"net/smtp"
	"os"
)

// SendEmail sends a 6-digit OTP to the specified recipient
func SendEmail(toEmail string, otp string) error {
	// 1. Get credentials from environment variables (configs)
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASS")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	// 2. Compose the email headers and body
	// Note: The empty line between headers and body is mandatory in SMTP
	subject := "Subject: Expense Tracker Verification Code\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(`
        <html>
            <body>
                <h2>Welcome to Expense Tracker!</h2>
                <p>Please use the following One-Time Password (OTP) to complete your registration:</p>
                <h1 style="color: #4A90E2; letter-spacing: 5px;">%s</h1>
                <p>This code is valid for 5 minutes. If you did not request this, please ignore this email.</p>
            </body>
        </html>`, otp)

	message := []byte(subject + mime + body)

	// 3. Authenticate with the SMTP server
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// 4. Send the email
	address := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	err := smtp.SendMail(address, auth, from, []string{toEmail}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
