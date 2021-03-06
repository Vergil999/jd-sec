package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"jd-sec/configs"
	"jd-sec/entity"
	"jd-sec/http_client"
	"jd-sec/logger"
	"jd-sec/mail"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var ticket string

var itemId string
var email string
var secTimeStr string

func init() {
	flag.StringVar(&itemId, "itemId", configs.Config.ItemId, "商品ID")
	flag.StringVar(&email, "email", configs.Config.Email, "接收登陆二维码邮箱")
	flag.StringVar(&secTimeStr, "secTime", configs.Config.SecTime, "抢购时间")
	flag.Parse()
}

func main() {
	if len(secTimeStr) == 13{
		secTimeStr = secTimeStr[0:10]
	}
	secTime, err := strconv.Atoi(secTimeStr)
	if err != nil {
		logger.Error("抢购时间参数错误")
		os.Exit(0)
	}
	sec := time.Unix(int64(secTime), 0)
	logger.Infof("商品ID:%v,接收二维码邮箱:%v,抢购时间:%v", itemId, email, sec.String())

	//等待秒杀开始前一分钟
	secUnix := sec.Unix()
	for {
		t := secUnix - time.Now().Unix()
		if t <= 60 {
			break
		} else {
			logger.Infof("预约时间未到.距离还有%v秒", t)
			time.Sleep(time.Second * 1)
		}
	}

	//请求登陆二维码并下载
	cookies := getAndSaveQrCode()
	var token = ""
	for _, v := range cookies {
		if v.Name == "wlfstk_smdl" {
			token = v.Value
		}
	}
	logger.Infof("获取到的token:%s", token)
	//发送登陆二维码到邮箱
	mail.SendEmail(email)

	headerMap := checkLogin(token, cookies)
	logger.Infof("确认抢购是否已开始...")
	for {
		t := sec.UnixNano() - time.Now().UnixNano()
		gapTime := time.Duration(t)
		if gapTime <= 0 {
			break
		} else if gapTime > time.Second*2 {
			logger.Infof("抢购时间未到...还有%v秒", gapTime.Seconds())
			time.Sleep(time.Millisecond * 1000)
		} else {
			logger.Infof("抢购时间未到...还有%v秒", gapTime.Seconds())
			time.Sleep(time.Millisecond * 100)
		}
	}

	secKilTask(cookies, headerMap)
}

func checkLogin(token string, cookies []*http.Cookie) map[string]string {
	baseNum := 9000000
	headerMap := make(map[string]string)
	headerMap["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/"
	headerMap["Referer"] = "https://www.jd.com/"
	//循环判断是否扫码
	for {
		time.Sleep(time.Duration(2) * time.Second)
		rand.Seed(time.Now().Unix())
		result := baseNum + rand.Intn(1000000)
		checkQrCodeUrl := fmt.Sprintf("https://qr.m.jd.com/check?callback=jQuery%v&appid=133&token=%v&_=%v", result, token, time.Now().Unix())
		resp, err := http_client.Exec("GET", checkQrCodeUrl, headerMap, cookies)
		if err != nil {
			logger.Errorf("获取二维码扫描状态失败,error:%s", err.Error())
			continue
		}
		respData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Errorf("获取二维码扫描状态失败,error:%s", err.Error())
			continue
		}
		qrResult := getQrCodeResult(respData)
		if qrResult.Code != 200 {
			logger.Errorf("尚未扫码，返回内容:%s", string(respData))
			continue
		} else {
			//扫码完成登陆 保存ticket
			ticket = qrResult.Ticket
			logger.Infof("登陆成功，获取到ticket:%s", qrResult.Ticket)
			break
		}
	}
	return headerMap
}

func secKilTask(cookies []*http.Cookie, headerMap map[string]string) {
	for {
		//获取thor cookie
		cookies = ticketValidation(ticket, cookies)

		//将商品加入购物车
		addItemToCart(headerMap, cookies)

		//addItemToCartCallback(headerMap, cookies)

		//结算
		getOrderInfo(cookies)
		success := submitOrder(cookies)
		if success {
			return
		}
		time.Sleep(time.Millisecond * 300)
	}
}

func getAndSaveQrCode() []*http.Cookie {
	var qrcodeUrl = fmt.Sprintf("https://qr.m.jd.com/show?appid=133&size=147&t=%v", time.Now().Unix())
	resp, _ := http_client.Exec("GET", qrcodeUrl, nil, nil)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	out, _ := os.Create("qrcode.png")
	io.Copy(out, bytes.NewReader(body))
	return resp.Cookies()
}

func getQrCodeResult(respData []byte) *entity.QRCodeResult {
	start := strings.Index(string(respData), "(")
	end := strings.LastIndex(string(respData), ")")
	jsonData := respData[start+1 : end]
	var qrResult = new(entity.QRCodeResult)
	json.Unmarshal(jsonData, &qrResult)
	return qrResult
}

/**
为了获取jd返回的thor cookie
*/
func ticketValidation(ticket string, cookies []*http.Cookie) []*http.Cookie {
	validationUrl := fmt.Sprintf(`https://passport.jd.com/uc/qrCodeTicketValidation?t=%s&ReturnUrl="%s"`, ticket, "https://item.jd.com/100030066232.html")
	headerMap := make(map[string]string)
	headerMap["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36"
	cartResp, _ := http_client.Exec("GET", validationUrl, headerMap, cookies)
	return cartResp.Cookies()
}

/**
加入购物车
*/
func addItemToCart(headerMap map[string]string, cookies []*http.Cookie) {
	cartUrl := fmt.Sprintf("https://cart.jd.com/gate.action?pcount=1&ptype=1&pid=%s", itemId)
	headerMap["Referer"] = "https://item.jd.com/"
	_, err := http_client.Exec("GET", cartUrl, headerMap, cookies)
	if err != nil {
		logger.Errorf("加入购物车失败,error:%s", err.Error())
		return
	}
}

//func addItemToCartCallback(headerMap map[string]string, cookies []*http.Cookie) {
//	cartUrl := fmt.Sprintf("https://cart.jd.com/gateResult?rcd=1&pid=%v&pc=1&eb=1&rid=%v&em=", itemId, time.Now().Unix())
//	headerMap["Referer"] = "https://item.jd.com/"
//	_, err := http_client.Exec("GET", cartUrl, headerMap, cookies)
//	if err != nil {
//		logger.Errorf("加入购物车失败,error:%s", err.Error())
//		return
//	}
//}

/**
结算
*/
func getOrderInfo(cookies []*http.Cookie) {
	orderInfoUrl := "https://trade.jd.com/shopping/order/getOrderInfo.action"
	headerMap := make(map[string]string)
	headerMap["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36"
	headerMap["referer"] = "https://cart.jd.com/"

	orderInfoResp, err := http_client.Exec("GET", orderInfoUrl, headerMap, cookies)
	_, err = ioutil.ReadAll(orderInfoResp.Body)
	if err != nil {
		logger.Errorf("结算失败,error:%s", err.Error())
		return
	}
}

func submitOrder(cookies []*http.Cookie) bool {
	//https://trade.jd.com/shopping/order/submitOrder.action?&presaleStockSign=1
	submitOrderUrl := "https://trade.jd.com/shopping/order/submitOrder.action?&presaleStockSign=1"
	headerMap := make(map[string]string)
	headerMap["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36"
	headerMap["referer"] = "https://trade.jd.com/shopping/order/getOrderInfo.action"

	orderInfoResp, err := http_client.Exec("POST", submitOrderUrl, headerMap, cookies)
	data, err := ioutil.ReadAll(orderInfoResp.Body)
	if err != nil {
		logger.Errorf("下单失败,error:%s", err.Error())
		return false
	}
	returnData := string(data)
	if strings.Contains(returnData, "\"success\":true") {
		logger.Infof("下单成功:%s", returnData)
		return true
	}
	logger.Infof("下单失败,返回结果:%v", returnData)
	return false
}
