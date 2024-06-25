package core

import (
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const oathUrl = "https://xui.ptlogin2.qq.com/cgi-bin/xlogin?appid=716027609&daid=383&style=33&login_text=%E7%99%BB%E5%BD%95&hide_title_bar=1&hide_border=1&target=self&s_url=https%3A%2F%2Fgraph.qq.com%2Foauth2.0%2Flogin_jump&pt_3rd_aid=101491592&pt_feedback_link=https%3A%2F%2Fsupport.qq.com%2Fproducts%2F77942%3FcustomInfo%3Dmilo.qq.com.appid101491592&theme=2&verify_theme="

var client = http.DefaultClient

func GetQrCode(cookies []*http.Cookie) (string, string) {
	req, _ := http.NewRequest("GET", "https://ssl.ptlogin2.qq.com/ptqrshow", nil)
	values := req.URL.Query()
	values.Add("appid", FromOathUrlGetParams("appid"))
	values.Add("e", "2")
	values.Add("l", "M")
	values.Add("s", "3")
	values.Add("d", "72")
	values.Add("v", "4")
	values.Add("t", "0.9288003470066593")
	values.Add("daid", FromOathUrlGetParams("daid"))
	values.Add("pt_3rd_aid", FromOathUrlGetParams("pt_3rd_aid"))
	values.Add("u1", "https://graph.qq.com/oauth2.0/login_jump")
	req.URL.RawQuery = values.Encode()
	for _, c := range cookies {
		req.AddCookie(c)
	}
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)

	ck := resp.Cookies()
	var qrsig string
	for _, c := range ck {
		if c.Name == "qrsig" {
			qrsig = c.Value
		}
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(body), qrsig
}
func FromOathUrlGetParams(key string) string {
	u, _ := url.Parse(oathUrl)
	query := u.Query()
	return query.Get(key)
}
func RequestForCookies() []*http.Cookie {
	req, _ := http.NewRequest("GET", oathUrl, nil)
	resp, _ := client.Do(req)
	cookies := resp.Cookies()
	for _, c := range cookies {
		log.Println(c.Name, c.Value)
	}
	return cookies
}
func IsQrCodeExpired(cookies []*http.Cookie, qrsig string, login_sig string, aid string, daid string, pt_3rd_aid string) string {
	req, _ := http.NewRequest("GET", "https://ssl.ptlogin2.qq.com/ptqrlogin", nil)
	values := req.URL.Query()
	values.Add("u1", "https://graph.qq.com/oauth2.0/login_jump")
	values.Add("ptqrtoken", strconv.Itoa(Hash33(qrsig)))
	values.Add("ptredirect", "0")
	values.Add("h", "1")
	values.Add("t", "1")
	values.Add("g", "1")
	values.Add("from_ui", "1")
	values.Add("ptlang", "2052")
	values.Add("action", "0-0-"+strconv.Itoa(int(time.Now().UnixMilli())))
	values.Add("js_ver", "24042411")
	values.Add("js_type", "1")
	values.Add("login_sig", login_sig)
	values.Add("pt_uistyle", "40")
	values.Add("aid", aid)
	values.Add("daid", daid)
	values.Add("pt_3rd_aid", pt_3rd_aid)
	values.Add("o1vId", "09fe2c176ac5fdfc385ce0e743ac86c6")
	values.Add("pt_js_version", "v1.48.3")
	req.URL.RawQuery = values.Encode()
	for _, c := range cookies {
		req.AddCookie(c)
	}
	req.AddCookie(&http.Cookie{Name: "qrsig", Value: qrsig})
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	log.Println(resp.Status, string(body))
	for _, c := range resp.Cookies() {
		log.Println(c.Name, c.Value)
	}
	return string(body)

}
func Hash33(t string) int {
	e := 0
	for n := 0; n < len(t); n++ {
		e += (e << 5) + int(t[n])
	}
	return 2147483647 & e
}
