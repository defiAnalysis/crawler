package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
	"gopkg.in/gomail.v2"
	"github.com/robfig/cron"
	_ "crawler/model"
)

const (
	HOST        = "smtp.163.com"
	SERVER_ADDR = "smtp.163.com:25"
	USER        = "duanchuanfu00@163.com" //发送邮件的邮箱
	PASSWORD    = "2866093520I@U"     //发送邮件邮箱的密码
)


type Email struct {
	to      string "to"
	subject string "subject"
	msg     string "msg"
}


func main() {
	cron2 := cron.New() //创建一个cron实例

	//执行定时任务（每5秒执行一次）
	err:= cron2.AddFunc("0/5 * 10,11,14,15,16,19,20 * * ?", getRequest)
	if err!=nil{
		fmt.Println(err)
	}

	//启动/关闭
	cron2.Start()
	defer cron2.Stop()
	select {
	//查询语句，保持程序运行，在这里等同于for{}
	}
}

func getRequest()  {
	res, err := http.Get("https://www.huobi.li/support/zh-cn/list/360000039481")
	if err != nil {
		// 错误处理
		fmt.Println("err:",err.Error())

		log.Fatal(err)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	doc.Find(".list-item").Each(func(i int, s *goquery.Selection) {
		data := s.Text()

		if strings.Contains(data, "上线")  && strings.Contains(data, "新币") && (strings.Contains(data,"全球观察区") || strings.Contains(data,"创新区")){
			fmt.Println(data)

			datas := strings.Split(data,"上线")
			usefulData := datas[1]


			hash := getMD5Hash(usefulData)



			fmt.Println("hash: ",hash)

		}
	})

	//content := fmt.Sprintf("<h4>%s<h4>","新币上线")
	//
	//sendMail("916140875@qq.com",content)
}

func sendMail(email string,content string) error {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", "duanchuanfu00@163.com","hayden")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "新币上线") //设置邮件主题
	m.SetBody("text/html", content)       //设置邮件正文
	// 第一个参数是host 第三个参数是发送邮箱，第四个参数 是邮箱密码
	d := gomail.NewDialer(HOST, 25, USER, PASSWORD)
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("错误：", err)
		return err
	}
	return nil
}

func getMD5Hash(text string) string {

	hasher := md5.New()

	hasher.Write([]byte(text))

	return hex.EncodeToString(hasher.Sum(nil))

}
