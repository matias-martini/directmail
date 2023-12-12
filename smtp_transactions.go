package directmail

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"net/textproto"
	"strings"
	"time"
)

type SMTPServer struct {
	address string
	port    int
	reader  *textproto.Reader
	writer  *textproto.Writer
	conn    net.Conn
}

func (s *SMTPServer) Connect() error {
	// Establish a TCP connection
	fmt.Println("Connecting to: ", s.address)
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", s.address, s.port), time.Second*1)
	if err != nil {
		return err
	}

	// Save the connection
	s.conn = conn

	// Create a new reader and writer
	s.reader = textproto.NewReader(bufio.NewReader(conn))
	s.writer = textproto.NewWriter(bufio.NewWriter(conn))

	// Conduct the SMTP conversation
	if _, _, err := s.reader.ReadResponse(220); err != nil {
		return err
	}

	return nil
}

func (s *SMTPServer) GetCapabilities() ([]string, error) {
	if err := s.writer.PrintfLine("EHLO martini.monster"); err != nil {
		return nil, err
	}

	// Read and print the server's capabilities
	capabilities := []string{}
	for {
		line, err := s.reader.ReadLine()
		if err != nil {
			return nil, err
		}
		fmt.Println(line)
		if strings.HasPrefix(line, "250 ") { // Last line of response starts with 250 space
			capabilities = append(capabilities, line)
			break
		}
	}

	return capabilities, nil
}

func (s *SMTPServer) StartTLS() error {
	err := s.writer.PrintfLine("STARTTLS")
	if err != nil {
		return err
	}

	_, _, err = s.reader.ReadResponse(220)
	if err != nil {
		return err
	}

	tlsConfig := &tls.Config{
		ServerName:         s.address,
		InsecureSkipVerify: true, // Set this to true only for testing purposes
	}

	tlsConn := tls.Client(s.conn, tlsConfig)
	s.reader = textproto.NewReader(bufio.NewReader(tlsConn))
	s.writer = textproto.NewWriter(bufio.NewWriter(tlsConn))

	return nil
}

// Get domain from sender email
func GetDomainFromSenderEmail(senderEmail string) string {
	parts := strings.Split(senderEmail, "@")
	if len(parts) != 2 {
		return ""
	}
	domain := parts[1]
	return domain
}

func (s *SMTPServer) SendEmail(senderEmail string, recipientEmail string, subject string, body string) error {
	// MAIL FROM command
	fmt.Println("MAIL FROM: ", senderEmail)
	if err := s.writer.PrintfLine("MAIL FROM:<%s>", senderEmail); err != nil {
		return err
	}
	if _, _, err := s.reader.ReadResponse(250); err != nil {
		return err
	}

	// RCPT TO command
	fmt.Println("RCPT: ", recipientEmail)
	if err := s.writer.PrintfLine("RCPT TO:<%s>", recipientEmail); err != nil {
		return err
	}
	if _, _, err := s.reader.ReadResponse(250); err != nil {
		return err
	}

	// DATA command
	fmt.Println("DATA")
	if err := s.writer.PrintfLine("DATA"); err != nil {
		return err
	}
	if _, _, err := s.reader.ReadResponse(354); err != nil {
		return err
	}

	// Sending email content
	emailHeadersAndBody := "From: " + senderEmail + "\r\n" +
		"To: " + recipientEmail + "\r\n" +
		"Subject:" + subject + "\r\n" +
		"Date: " + time.Now().Format(time.RFC1123Z) + "\r\n" +
		"Message-ID: <" + generateMessageID() + ">\r\n" +
		"MIME-Version: 1.0\r\n" +
		"\r\n" +
		body + "\r\n"

	fmt.Println("Getting domain from sender email..")
	senderDomain := GetDomainFromSenderEmail(senderEmail)
	fmt.Println("Signing..")

	signedEmailContent, err := signWithDKIM(senderDomain, emailHeadersAndBody)

	if err != nil {
		return err
	}

	fmt.Println("Signed Email Content:\n", signedEmailContent)

	if err := s.writer.PrintfLine(emailHeadersAndBody + "\r\n."); err != nil {
		return err
	}
	if _, _, err := s.reader.ReadResponse(250); err != nil {
		return err
	}

	return nil
}
