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
	yaml.Unmarshal(b, conf)

	if conf.User == "" {
		log.Fatalln("DAILY_EMAIL environment variable not provided, cannot authenticate.")
	} else if conf.Password == "" {
		log.Fatalln("DAILY_PASS environment variable not provided, cannot authenticate.")
	}

	auth := sasl.NewPlainClient(conf.User, conf.User, conf.Password)

	// We join the strings here for a neat 'To: <addresses>' formatted message, otherwise the recipients are BCC'd into the email.
	recipients := strings.Join(conf.To, ",")

	fmtMsg := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s\r\n", recipients, conf.Subject, conf.Content)
	msg := strings.NewReader(fmtMsg)

	serverAddr, ok := providerToSMTPServer[conf.Provider]
	if !ok {
		log.Fatalf("email provider %s does not exist in the mapping", conf.Provider)
	}

	err = smtp.SendMail(serverAddr, auth, conf.User, conf.To, msg)
	if err != nil {
		log.Fatal(err)
	}
}
