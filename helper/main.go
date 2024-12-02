package main

import (
	"context"
	"io"
	"log"

	"github.com/segmentio/kafka-go"
)

// type emailHelper struct {
// 	dialer *gomail.Dialer
// }

// func (e *emailHelper) SendEmailForgotPassword(email, token string) {
// 	link := "https://nucrea.ru?token=" + token

// 	const MSG_TEXT = `
// 	<html>
// 		<head>
// 		</head>
// 		<body>
// 			<p>This message was sent because you forgot a password</p>
// 			<p>To change a password, use <a href="{{Link}}"/>this</a> link</p>
// 		</body>
// 	</html>
// 	`
// 	msgText := strings.ReplaceAll(MSG_TEXT, "{{Link}}", link)

// 	m := gomail.NewMessage()
// 	m.SetHeader("From", "email")
// 	m.SetHeader("To", email)
// 	m.SetHeader("Subject", "Hello!")
// 	m.SetBody("text/html", msgText)

// 	if err := d.DialAndSend(m); err != nil {
// 		panic(err)
// 	}
// }

func main() {
	const (
		SMTP_SERVER   = ""
		SMTP_PORT     = 0
		SMTP_LOGIN    = ""
		SMTP_PASSWORD = ""
	)

	ctx := context.Background()

	// d := gomail.NewDialer(SMTP_SERVER, SMTP_PORT, SMTP_LOGIN, SMTP_PASSWORD)

	// conn, err := kafka.DialLeader(ctx, "", "")
	// if err != nil {
	// 	panic(err)
	// }
	// defer conn.Close()

	log.Println("starting reader...")

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "topic-A",
		// Partition: 0,
		GroupID: "consumer-group-id",
	})
	// r.SetOffset(kafka.FirstOffset)

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
		}
	}
}
