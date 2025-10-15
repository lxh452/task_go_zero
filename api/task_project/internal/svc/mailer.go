package svc

import (
    "crypto/tls"
    "fmt"
    "net/smtp"
    "strings"

    "task_Project/api/task_project/internal/config"
)

type Mailer interface {
    Send(to []string, subject string, body string) error
}

type smtpMailer struct {
    host string
    port int
    username string
    password string
    from string
}

func MustNewMailer(c config.SMTPConf) Mailer {
    return &smtpMailer{host: c.Host, port: c.Port, username: c.Username, password: c.Password, from: c.From}
}

func (m *smtpMailer) Send(to []string, subject string, body string) error {
    addr := fmt.Sprintf("%s:%d", m.host, m.port)
    headers := []string{
        fmt.Sprintf("From: %s", m.from),
        fmt.Sprintf("To: %s", strings.Join(to, ",")),
        fmt.Sprintf("Subject: %s", subject),
        "MIME-Version: 1.0",
        "Content-Type: text/html; charset=\"UTF-8\"",
    }
    msg := []byte(strings.Join(headers, "\r\n") + "\r\n\r\n" + body)

    auth := smtp.PlainAuth("", m.username, m.password, m.host)
    tlsconfig := &tls.Config{ServerName: m.host}
    // Use STARTTLS if server supports; for simplicity, dial and StartTLS is omitted here
    // Many providers on 587 expect StartTLS; however net/smtp is deprecated. For production consider a modern SMTP client.
    return smtp.SendMail(addr, auth, m.from, to, msg)
}


