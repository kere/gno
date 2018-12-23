package gno

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"crypto/tls"
	"net"
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

//mailDial return a smtp client
func mailDial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

// MailUser mail
type MailUser struct {
	Name string
	Mail string
}

// MailTLS use tls
type MailTLS struct {
	Host string
	Port int
	From MailUser
	Auth smtp.Auth
}

// NewMailTLS new
func NewMailTLS(host string, port int, from MailUser, password string) MailTLS {
	auth := newMailAuth(from.Mail, password)
	return MailTLS{Host: host, Port: port, From: from, Auth: &auth}
}

//SendTo SendMail to one use TLS
func (m *MailTLS) SendTo(to MailUser, title string, body []byte) error {
	h := make(map[string]string)
	usr := ""
	if to.Name != "" {
		usr = to.Name + "<" + to.Mail + ">"
	} else {
		usr = to.Mail
	}
	h["To"] = usr
	toUsers := []string{usr}
	return m.send(h, toUsers, title, body)
}

//Send SendMail to uses TLS
func (m *MailTLS) Send(toName string, toAddrs []string, title string, body []byte) error {
	h := make(map[string]string)
	h["To"] = toName
	return m.send(h, toAddrs, title, body)
}

// //Send SendMail use TLS
// func (m *MailTLS) Send(toUserName string, to []string, title string, body []byte) (err error) {
//   header["To"] =
// }
func (m *MailTLS) send(header map[string]string, to []string, title string, body []byte) (err error) {
	header["From"] = m.From.Name + "<" + m.From.Mail + ">"
	header["Subject"] = title
	header["Content-Type"] = "text/html; charset=UTF-8"

	buf := bytes.NewBuffer(nil)
	for k, v := range header {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	buf.WriteString("\r\n")
	buf.Write(body)

	addr := fmt.Sprintf("%s:%d", m.Host, m.Port)
	//create smtp client
	c, err := mailDial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	if m.Auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(m.Auth); err != nil {
				return err
			}
		}
	}
	if err = c.Mail(m.From.Mail); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(buf.Bytes())
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

type loginAuth struct {
	username, password string
}

// newMailAuth func
func newMailAuth(username, password string) loginAuth {
	return loginAuth{username, password}
}

// Start Func
func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	// return "LOGIN", []byte{}, nil
	return "LOGIN", []byte(a.username), nil
}

// Next func
func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		}
	}
	return nil, nil
}
