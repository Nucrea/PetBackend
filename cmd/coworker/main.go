package main

import (
	"backend/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
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
		<p>To change a password, use <a href="{{.Link}}">this</a>link</p>
	</body>
</html>
`

type HtmlTemplate struct {
	Link string
}

func SendEmailForgotPassword(dialer *gomail.Dialer, from, to, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(from, "Pet Backend"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", body)

	return dialer.DialAndSend(m)
}

type Config struct {
	App struct {
		LogFile    string `yaml:"logFile"`
		ServiceUrl string `yaml:"serviceUrl"`
	}

	Kafka struct {
		Brokers         []string `yaml:"brokers"`
		Topic           string   `yaml:"topic"`
		ConsumerGroupId string   `yaml:"consumerGroupId"`
	} `yaml:"kafka"`

	SMTP struct {
		Server   string `yaml:"server"`
		Port     int    `yaml:"port"`
		Login    string `yaml:"login"`
		Password string `yaml:"password"`
		Email    string `yaml:"email"`
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

	dialer := gomail.NewDialer(config.SMTP.Server, config.SMTP.Port, config.SMTP.Login, config.SMTP.Password)

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: config.Kafka.Brokers,
		Topic:   config.Kafka.Topic,
		GroupID: config.Kafka.ConsumerGroupId,
	})

	logger, err := logger.New(
		ctx,
		logger.NewLoggerOpts{
			Debug:      true,
			OutputFile: config.App.LogFile,
		},
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	logger.Printf("coworker service started\n")

	template, err := template.New("verify-email").Parse(MSG_TEXT)
	if err != nil {
		log.Fatal(err)
	}

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

		link := fmt.Sprintf("%s/verify-user?token=%s", config.App.ServiceUrl, value.Token)

		builder := &strings.Builder{}
		if err := template.Execute(builder, HtmlTemplate{link}); err != nil {
			log.Printf("failed to execute html template: %s\n", err.Error())
			continue
		}

		if err := SendEmailForgotPassword(dialer, config.SMTP.Email, value.Email, builder.String()); err != nil {
			log.Printf("failed to send email: %s\n", err.Error())
			continue
		}
	}
}
