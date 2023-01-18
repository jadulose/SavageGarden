package main

import (
	"fmt"
	"github.com/wneessen/go-mail"
)

type MailConf struct {
	Addr     string
	Port     int
	Ssl      bool
	User     string
	Password string
	From     string
}

func (c *MailConf) Open() (*mail.Client, error) {
	client, err := mail.NewClient(c.Addr, mail.WithPort(c.Port), mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(c.User), mail.WithPassword(c.Password))
	if err != nil {
		return nil, err
	}
	client.SetSSL(c.Ssl)
	return client, nil
}

func (c *MailConf) CreateTestMsg(mailTo string) (*mail.Msg, error) {
	m := mail.NewMsg()
	if err := m.From(fmt.Sprintf("Savage Garden <%s>", c.From)); err != nil {
		return nil, err
	}
	if err := m.To(mailTo); err != nil {
		return nil, err
	}
	m.Subject("测试Savage Garden服务邮件发送功能")
	m.SetBodyString(mail.TypeTextPlain, "恭喜你，收到这封邮件，表示邮件配置成功！")
	return m, nil
}
