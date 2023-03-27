package mail

import (
	"gopkg.in/gomail.v2"
)

type EmailInfo struct {
	Username string // 账号
	Password string // 密码
	Host     string // 服务器地址
	Port     int    // 端口 默认465
}

type Emailer struct {
	EmailInfo EmailInfo
	Dialer    *gomail.Diale
}

func NewEmailer(emailInfo EmailInfo) *Emailer {
	dialer := gomail.NewDialer(emailInfo.Host, emailInfo.Port, emailInfo.Username, emailInfo.Password)
	return &Emailer{
		Dialer:    dialer,
		EmailInfo: emailInfo,
	}
}

func (e *Emailer) Send(mailTo []string, send_name, title, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.EmailInfo.Username)
	//这种方式可以添加别名，即“XX官方”
	//说明：如果是用网易邮箱账号发送，以下方法别名可以是中文，如果是qq企业邮箱，以下方法用中文别名，会报错，需要用上面此方法转码
	//m.SetHeader("From", "FB Sample"+"<"+mailConn["user"]+">") //这种方式可以添加别名，即“FB Sample”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
	//m.SetHeader("From", mailConn["user"])
	m.SetHeader("To", mailTo...)  //发送给多个用户
	m.SetHeader("Subject", title) //设置邮件主题
	m.SetBody("text/html", body)  //设置邮件正文

	return e.Dialer.DialAndSend(m)
}
