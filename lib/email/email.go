package email

import (
	"crypto/tls"
	"fmt"
	"math/rand/v2"
	"net/smtp"
	"strconv"
	"time"
)

var Config struct {
	Host             string
	Port             int
	ToAddr, FromAddr string
	Password         string
}

const emailTemplate = `From: %s
To: %s
Subject: %s
Date: %s
Content-Type: text/plain; charset=UTF-8
MIME-Version: 1.0
Message-ID: <%s>

%s`

func ServerError(fault error) error {
	msgid := strconv.FormatInt(rand.Int64(), 10) + "@" + Config.Host
	msg := fmt.Sprintf(emailTemplate, Config.FromAddr, Config.ToAddr,
		"Error Report", time.Now().Format(time.RFC1123Z), msgid, fault)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         Config.Host,
	}

	hostWithPort := Config.Host + ":" + strconv.Itoa(Config.Port)
	conn, err := tls.Dial("tcp", hostWithPort, tlsConfig)
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(conn, Config.Host)
	if err != nil {
		return err
	}
	defer client.Close()

	auth := smtp.PlainAuth("", Config.FromAddr, Config.Password, Config.Host)
	if err := client.Auth(auth); err != nil {
		return err
	}

	if err := client.Mail(Config.FromAddr); err != nil {
		return err
	}

	if err := client.Rcpt(Config.ToAddr); err != nil {
		return err
	}

	wc, err := client.Data()
	if err != nil {
		return err
	}
	defer wc.Close()

	if _, err = wc.Write([]byte(msg)); err != nil {
		return err
	}
	return nil
}
