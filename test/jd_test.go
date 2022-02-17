package test

import (
	"crypto/tls"
	"fmt"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"jd-sec/http_client"
	"jd-sec/logger"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test1(t *testing.T) {

	cartUrl := fmt.Sprintf("https://cart.jd.com/gate.action?pcount=1&ptype=1&pid=%s", "100012955747")
	headerMap := make(map[string]string)
	headerMap["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36"
	headerMap["cookie"] = "thor=FFD6505A6F6BDEA8CE0726F1F2EA218E6B06C460783009620477DC3888145CCDA26B1F018D21FBEC1DBB686109F83DD09EA597451F501CFD6E9354A3C54C72C15CAD7E020441A350A073877FB52DA66064046ADE09FA4A659091CA59EA1ED149D5FF00CAE416AD21DEF2F14E69687B2483F66FC52FE697D72A390A974707424EA07DE415F03D85E780A7EEA2E7BF889722BB2102FE807FB6677D724B9E1A7CA7;"
	headerMap["referer"] = "https://cart.jd.com/gateResult?rcd=1&pid=100012955747&pc=1&eb=1&rid=1642665576318&em="

	cartResp, err := http_client.Exec("GET", cartUrl, headerMap, nil)
	cartRespData, err := ioutil.ReadAll(cartResp.Body)
	if err != nil {
		logger.Errorf("加入购物车失败,error:%s", err.Error())
		return
	}
	println(string(cartRespData))
}

func TestOrderInfo(t *testing.T) {

	cartUrl := fmt.Sprintf("https://trade.jd.com/shopping/order/getOrderInfo.action")
	headerMap := make(map[string]string)
	headerMap["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36"
	headerMap["cookie"] = "thor=FFD6505A6F6BDEA8CE0726F1F2EA218E6B06C460783009620477DC3888145CCDBD6E521495E31E7CD97A09D431ACE0084C66DAFFE6B925E1900BE9ACF55F8942841A790D2AE9B49B659798AA7A6D3DD704C0CE8AE7747DB99F9A291CFD0806EDBCA32795BF0612F082B67A18108217A621919257FAEA687FA789834F08618E3C0D4094AD74066C0157702F60A98C558C27F19DD1C5C000C590E86509564EDD58;"
	headerMap["referer"] = "https://cart.jd.com/"

	cartResp, err := http_client.Exec("GET", cartUrl, headerMap, nil)
	cartRespData, err := ioutil.ReadAll(cartResp.Body)
	if err != nil {
		logger.Errorf("结算页面失败,error:%s", err.Error())
		return
	}
	println(string(cartRespData))
}

func TestTicketValidation(t *testing.T) {

	validationUrl := fmt.Sprintf(`https://passport.jd.com/uc/qrCodeTicketValidation?t=%s&ReturnUrl="%s"`, "AAEAMHRyPKDY2-Rk0OeLLt-UOAZRa-leQQUihkuf8Qbt3XSZEIxNAPdkvqmcXbUPMhodsg", "https://item.jd.com/100030066232.html")
	headerMap := make(map[string]string)
	headerMap["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36"
	headerMap["cookie"] = "shshshfpa=07cae11f-7b8c-ad78-807a-c9c83fc5612b-1627022135; __jdu=16270221349561790981632; _c_id=sldr370z2yuewmd8ogs1627460122911zn8m; DeviceSeq=4f855c849df84d809b2146d96ed325b3; pinId=pSTA13Al3c7H1hhVenHW8rV9-x-f3wj7; shshshfpb=btGffzAgocE%2BZ8LcoRMromw%3D%3D; unpl=JF8EAQZnNSttXxxdBRpSSEdHSF9cWw1bSUQLamJWVAkIGQcHHgBPEUB7XlVdXhRKFx9uYRRUWFNJUg4bBisSFXtdVV9dDE4UAGlmNWQYGxsGUQcCHhUTTVVRWFwPTBAGaGMGXVxZV1cFGwEaFBNLXVdfWzhIFzNuVwVTXl5KVwYbAB4aFExdUlZcAEIXBG1kNWRdUE9kBRoDGhMYQ1hQVl0KexcCbmQMXV1fTFBrHgMaExNLVVNZVGZJFgJtZAxcVF9JZAUSAhkaEUpdU11tCXsWMy0J09vIgPj00ZOMz5S_V4TozIS5-s6s07Ok7V0VS1MGHQMYERBJWFxaWghNHwJnbgVTX1t7VTUY; __jdv=122270672|www.zhihu.com|t_1001542270_1002881665_4000219668_3003289863|tuiguang|6f810abdf2394ab3b944b1eacb343e2a|1641537436281; pin=jd_7c028a8e03721; unick=%E9%85%B8%E8%8F%9C%E9%A6%85%E9%A5%BC; _tp=DBr4HkbA2TL%2FPCCY%2FmHASbIkQO3QdM4CuTz2rDDEGoE%3D; _pst=jd_7c028a8e03721; areaId=1; user-key=4afb52d6-6a99-4876-819e-25d20362a5f9; ipLoc-djd=1-2901-4135-0.1464564712; ipLocation=%u5317%u4eac; PCSYCityID=CN_110000_110100_0; mt_xid=V2_52007VwMVUVtZUVkfSxxUAWADFFpcUFtaGEsabFViA0ICVQtQRhgbHw4ZYgUUUEELAlwdVUxYUGEDFFIJD1NfTHkaXQZkHxNSQVtTSx9MElwDbAASYl1oUmoWQRBZA2MGEldbX1RYGEkaWwBnMxBWXFo%3D; shshshfp=13726cb3f2eaec41277a83b2e6aa91bb; ip_cityCode=2802; _distM=236737574026; __jdc=122270672; token=9d059234b23ab2bc2a96e20f2a46c69b,3,912639; __tk=05a6f005610d955ba2a3513d636f98e8,3,912639; __jda=122270672.16270221349561790981632.1627022134.1642731367.1642750827.26; TrackID=1rRLJlgf1hiloG-ovJdcd1A6PfKqUHlQn9c3p27Kf1xqdeoom3Vc4UOUMWsnN81NVpmwMgo5S18Fui7T5Wg6hKHvGu92xHoknSRUwzwmxrfHI-CQ9lPQGTzW8Wfeq13ZZ; shshshsID=d0c9d1a36ed21e83141fc4b11e15d0b8_8_1642752324690; cn=2; alc=/NxsHMZM5kEGO+M/oluiTA==; _t=nO0nBz/8jYLVFugVt1lUkZxSdiX4ovYnhhYuZ0ND5LQ=; __jdb=122270672.16.16270221349561790981632|26.1642750827; wlfstk_smdl=67cg57fiy690btsqbvwi10c5owp721k8; 3AB9D23F7A4B3C9B=TO5WPGM2SICB6OWVDT4H4L3STDR43BDUGEGEYFOS2WUNE3JP34CBIU4QVIZ3IM67GEWPZH6HC4WPF3FDHX65SDO4WI"
	cartResp, err := http_client.Exec("GET", validationUrl, headerMap, nil)
	cartRespData, err := ioutil.ReadAll(cartResp.Body)
	if err != nil {
		logger.Errorf("结算页面失败,error:%s", err.Error())
		return
	}
	for _, c := range cartResp.Cookies() {
		println(c.Name + ":" + c.Value)
	}
	println(string(cartRespData))
}

func TestSendMail(t *testing.T) {
	//dwzlhixmpnhibdde
	d := gomail.NewDialer("smtp.qq.com", 25, "669484592@qq.com", "dwzlhixmpnhibdde")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	m := gomail.NewMessage()
	m.SetHeader("From", "669484592@qq.com")
	m.SetHeader("To", "lyq7845@126.com")
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")
	pa, _ := os.Getwd()
	index := strings.LastIndex(pa, "/")
	qrCodePath := pa[0:index] + "/qrcode.png"
	m.Attach(qrCodePath)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func TestGetPWD(t *testing.T) {
	//D:\Go\workspace\port
	pa, _ := os.Getwd()
	index := strings.LastIndex(pa, "/")
	println(pa[0:index]+"/qrcode.png")
}

func TestGetCurrentDirectory(t *testing.T) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
}

