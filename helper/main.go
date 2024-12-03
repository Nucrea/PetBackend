package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"

	"github.com/segmentio/kafka-go"
	"gopkg.in/gomail.v2"
	"gopkg.in/yaml.v3"
)

const MSG_TEXT = `
<html>
	<head>
	</head>
	<body>
		<p>This message was sent because you forgot a password</p>
		<p>To change a password, use <a href="{{Link}}"/>this</a> link</p>
	</body>
</html>
`

func SendEmailForgotPassword(dialer *gomail.Dialer, from, to, token string) error {
	link := "localhost:8080/restore-password?token=" + token

	msgText := strings.ReplaceAll(MSG_TEXT, "{{Link}}", link)

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(from, "Pet Backend"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", msgText)

	return dialer.DialAndSend(m)
}

type Config struct {
	Kafka struct {
		Brokers []string `yaml:"brokers"`
		Topic   string   `yaml:"topic"`
	}

	SMTP struct {
		Server   string `yaml:"server"`
		Port     int    `yaml:"port"`
		Email    string `yaml:"email"`
		Password string `yaml:"password"`
	} `yaml:"smtp"`
}

func main() {
	ctx := context.Background()

	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err.Error())
	}

	config := &Config{}
	if err := yaml.Unmarshal(configFile, config); err != nil {
		log.Fatal(err.Error())
	}

	dialer := gomail.NewDialer(config.SMTP.Server, config.SMTP.Port, config.SMTP.Email, config.SMTP.Password)

	log.Println("starting reader...")

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: config.Kafka.Brokers,
		Topic:   config.Kafka.Topic,
		GroupID: "consumer-group-id",
	})

	log.Println("reader started")

	for {
		msg, err := r.FetchMessage(ctx)
		if err == io.EOF {
			log.Fatal("EOF")
			return
		}
		if err != nil {
			log.Fatal(err.Error())
			return
		}

		log.Printf("offset: %d, partition: %d, key: %s, value: %s\n", msg.Offset, msg.Partition, string(msg.Key), string(msg.Value))

		if err := r.CommitMessages(ctx, msg); err != nil {
			log.Fatalf("failed to commit: %s\n", err.Error())
			continue
		}

		value := struct {
			Email string `json:"email"`
			Token string `json:"token"`
		}{}

		if err := json.Unmarshal(msg.Value, &value); err != nil {
			log.Fatalf("failed to unmarshal: %s\n", err.Error())
			continue
		}

		if err := SendEmailForgotPassword(dialer, config.SMTP.Email, value.Email, value.Token); err != nil {
			log.Fatalf("failed to send email: %s\n", err.Error())
			continue
		}
	}
}
