package directmail

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func SendEmail(senderEmail string, recipientEmail string, subject string, body string) error {
	//Create instance of SMTPServer
	mailServer, err := GetSMTPServer(recipientEmail)

	if err != nil {
		panic(err)
	}

	err = SendEmailToSMTPServer(senderEmail, recipientEmail, subject, body, mailServer, 25)

	return err
}

func SendEmailToSMTPServer(senderEmail string, recipientEmail string, subject string, body string, SMTPServerAddress string, SMTPServerPort int) error {
	smtpServer := SMTPServer{address: SMTPServerAddress, port: SMTPServerPort}

	//Connect to the remote SMTP server
	smtpServer.Connect()

	//Get capabilities of the remote SMTP server
	capabilities, err := smtpServer.GetCapabilities()

	if err != nil {
		return err
	}

	// If the server supports STARTTLS, switch to TLS
	if stringInSlice("STARTTLS", capabilities) {
		err := smtpServer.StartTLS()

		if err != nil {
			return err
		}
	}

	// Send email
	err = smtpServer.SendEmail(senderEmail, recipientEmail, subject, body)

	if err != nil {
		return err
	}

	return nil
}
