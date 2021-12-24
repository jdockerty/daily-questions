package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Provider string
	User     string
	Password string
	Content  string
	From     string
	Subject  string
	To       []string
}

var providerToSMTPServer = map[string]string{
	"gmail":   "smtp.gmail.com:25",
	"yahoo":   "smtp.mail.yahoo.com:465",
	"outlook": "smtp-mail.outlook.com:25",
}

// getEnv is used to retrieve an environment variale with a default if it not available.
func getEnv(key, defaultValue string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return val
}

func main() {

	configPath := getEnv("DAILY_CONF", "daily-conf.yaml")
	b, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("cannot read configuration file at %s: %s", configPath, err)
	}

	conf := &Config{}
	err = yaml.Unmarshal(b, conf)
	if err != nil {
		log.Fatalf("Unable to unmarshal configuration at %s: %s\n", configPath, err)
	}
	
	log.Println("Configuration unmarshaled successfully.")

	if conf.User == "" {
		log.Fatalln("DAILY_EMAIL environment variable not provided, cannot authenticate.")
	} else if conf.Password == "" {
		log.Fatalln("DAILY_PASS environment variable not provided, cannot authenticate.")
	}

	auth := sasl.NewPlainClient(conf.User, conf.User, conf.Password)
	log.Println("Created security and authentication layer.")

	// We join the strings here for a neat 'To: <addresses>' formatted message, otherwise the recipients are BCC'd into the email.
	recipients := strings.Join(conf.To, ",")
	log.Printf("Contents will be sent to: %s\n", recipients)

	fmtMsg := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s\r\n", recipients, conf.Subject, conf.Content)
	msg := strings.NewReader(fmtMsg)
	log.Printf("Message content is:\n%s\n", conf.Content)

	serverAddr, ok := providerToSMTPServer[conf.Provider]
	if !ok {
		log.Fatalf("email provider %s does not exist in the mapping", conf.Provider)
	}

	log.Println("Using server address:", serverAddr)

	err = smtp.SendMail(serverAddr, auth, conf.User, conf.To, msg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Email sent successfully.")
	
}
