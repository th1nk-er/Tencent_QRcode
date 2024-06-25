package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strings"
	"tencent_qrcode/core"
)

var cookies []*http.Cookie
var qrsig string

func ReturnQrCode(c *gin.Context) {
	cookies = core.RequestForCookies()
	var code string
	code, qrsig = core.GetQrCode(cookies)
	c.String(200, code)
}

func IsQrCodeExpired(c *gin.Context) {
	var loginSig string
	for _, c := range cookies {
		if c.Name == "pt_login_sig" {
			loginSig = c.Value
		}
	}
	res := core.IsQrCodeExpired(cookies, qrsig, loginSig, core.FromOathUrlGetParams("appid"), core.FromOathUrlGetParams("daid"), core.FromOathUrlGetParams("pt_3rd_aid"))
	status := "二维码未失效"
	if strings.Contains(res, "二维码认证中") {
		status = "二维码认证中"
	} else if strings.Contains(res, "二维码已失效") {
		status = "二维码已失效"
	} else if strings.Contains(res, "登录成功") {
		reg, _ := regexp.Compile("https://[?.a-zA-Z0-9&_=/:%]+")
		str := reg.FindAllString(res, -1)
		status = str[0]
	}
	c.String(200, status)
}
