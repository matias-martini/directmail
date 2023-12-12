package directmail_test

import (
	"directmail"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
)

type Email struct {
	IsRelayed       bool   `json:"isRelayed"`
	ID              string `json:"id"`
	From            string `json:"from"`
	To              string `json:"to"`
	ReceivedDate    string `json:"receivedDate"`
	Subject         string `json:"subject"`
	AttachmentCount int    `json:"attachmentCount"`
	IsUnread        bool   `json:"isUnread"`
}

func TestSendEmail(t *testing.T) {
	// directmail.SendEmailToSMTPServer(
	// 	"matias@martini.monster",
	// 	"dawejip133@mcenb.com",
	// 	"Hola monstro",
	// 	"Aca el cuerpo del mensaje",
	// 	"192.168.0.39",
	// 	25,
	// )
	// mailServer, err := directmail.GetSMTPServer("dawejip133@mcenb.com")

	// fmt.Printf("Mail Server: '%+v'\n", mailServer)

	err := directmail.SendEmailToSMTPServer(
		"matias@martini.monster",
		"dawejip133@mcenb.com",
		"Hola monstro",
		"Aca el cuerpo del mensaje",
		"localhost",
		2525,
	)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	response, err := http.Get("http://localhost:8080/api/Messages")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(responseData))

	var emails []Email
	err = json.Unmarshal([]byte(responseData), &emails)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Parsed Emails: %+v\n", emails)
	emailCount := len(emails)

	if emailCount != 1 {
		t.Errorf("Expected only one email to be sent, but got %d", emailCount)
	}
}
