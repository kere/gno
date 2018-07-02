package gno

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"
)

// SMTPlainMail main
type SMTPlainMail struct {
	auth smtp.Auth
	host string
	port int
	from mail.Address
}

// NewSMTPlainMail plain
func NewSMTPlainMail(host string, from mail.Address, user, password string) *SMTPlainMail {
	return &SMTPlainMail{
		from: from,
		auth: smtp.PlainAuth("", user, password, host),
		host: host}
}

// SetPort port
func (m *SMTPlainMail) SetPort(p int) *SMTPlainMail {
	m.port = p
	return m
}

// SendHTML html
func (m *SMTPlainMail) SendHTML(to []*mail.Address, title string, body []byte) error {
	return m.Send(to, title, body, "html")
}

// SendPlain plain
func (m *SMTPlainMail) SendPlain(to []*mail.Address, title string, body []byte) error {
	return m.Send(to, title, body, "plain")
}

// Send mail
func (m *SMTPlainMail) Send(to []*mail.Address, title string, body []byte, mailType string) error {
	mailList := make([]string, len(to))
	mailAddrList := make([]string, len(to))
	for i, v := range to {
		mailList[i] = v.Address
		mailAddrList[i] = v.String()
	}

	header := make(map[string]string)
	header["From"] = m.from.String()
	header["To"] = strings.Join(mailAddrList, ",")
	header["Subject"] = title
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = fmt.Sprint("text/", mailType, "; charset=\"UTF-8\"")
	// header["Content-Type"] = "Multipart/Alternative;"
	header["Content-Transfer-Encoding"] = "base64"

	buf := bytes.NewBuffer([]byte(""))

	for k, v := range header {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	buf.WriteString("\r\n")
	src := body
	bufBody := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(bufBody, src)
	buf.Write(bufBody)
	port := 25
	if m.port != 0 {
		port = m.port
	}

	addr := fmt.Sprint(m.host, ":", port)

	return smtp.SendMail(
		addr,
		m.auth,
		m.from.Address,
		mailList,
		buf.Bytes())
}
