package directmail

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"sort"
	"strings"
	"time"

	dkim "github.com/toorop/go-dkim"
)

func GetSMTPServer(email string) (string, error) {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid email address")
	}
	domain := parts[1]

	// Perform MX record lookup
	mxRecords, err := lookupMXRecords(domain)
	if err != nil {
		return "", fmt.Errorf("MX record lookup failed: %s", err)
	}

	// Attempt to connect to each MX record host
	for _, mx := range mxRecords {
		host := mx.Host
		conn, err := net.DialTimeout("tcp", host+":25", 1*time.Second)
		if err == nil {
			conn.Close() // Close the connection if successful
			return host, nil
		}
	}

	return "", fmt.Errorf("no SMTP server found for domain: %s", domain)
}

// lookupMXRecords performs DNS MX record lookup for the given domain
func lookupMXRecords(domain string) ([]*net.MX, error) {
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return nil, err
	}

	// Sorting MX records by preference
	sort.Slice(mxRecords, func(i, j int) bool {
		return mxRecords[i].Pref < mxRecords[j].Pref
	})

	return mxRecords, nil
}

func generateMessageID() string {
	return fmt.Sprintf("%d.%s@%s", time.Now().UnixNano(), randString(8), "martini.monster")
}

func randString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func signWithDKIM(domain, emailHeadersAndBody string) (string, error) {
	// Read the private key from file
	fmt.Printf("Looking for private key in env var: DIRECTEMAIL_PRIVATE_KEY\n")

	privateKey, privateKeyExists := os.LookupEnv("DIRECTEMAIL_PRIVATE_KEY")

	if !privateKeyExists {
		return "", fmt.Errorf("DIRECTEMAIL_PRIVATE_KEY env var not found")
	}
	fmt.Printf("Signing with private key: %s\n", privateKey)

	dkimOptions := dkim.NewSigOptions()
	dkimOptions.PrivateKey = []byte(privateKey)
	dkimOptions.Domain = domain
	dkimOptions.Selector = "default"
	dkimOptions.Headers = []string{"from", "date", "mime-version", "received", "received"}
	dkimOptions.Canonicalization = "relaxed/relaxed"

	// Prepare email for signing
	emailToSign := []byte(emailHeadersAndBody)

	// Sign the email
	err := dkim.Sign(&emailToSign, dkimOptions)
	if err != nil {
		return "", err
	}

	// Convert signed email back to string
	signedEmail := string(emailToSign)
	return signedEmail, nil
}
