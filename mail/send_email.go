package mail

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
	"jd-sec/logger"
	"os"
)

func SendEmail(toEmail string) error {
	d := gomail.NewDialer("smtp.qq.com", 25, "669484592@qq.com", "dwzlhixmpnhibdde")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	m := gomail.NewMessage()
	m.SetHeader("From", "669484592@qq.com")
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "sec-login")
	m.SetBody("text/html", "用京东扫一扫附件二维码扫码登陆")
	pa, _ := os.Getwd()
	qrCodePath := pa + "/qrcode.png"
	m.Attach(qrCodePath)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		logger.Errorf("发送邮件失败，异常:%v", err)
		return err
	}
	return nil
}
