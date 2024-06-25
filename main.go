package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"tencent_qrcode/core"
	"tencent_qrcode/handlers"
	"time"
)

func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))
	r.GET("/qrcode", handlers.ReturnQrCode)
	r.GET("/status", handlers.IsQrCodeExpired)
	//test()
	_ = r.Run()
}

func test() {
	cookies := core.RequestForCookies()
	_, qrsig := core.GetQrCode(cookies)
	var loginSig string
	for _, c := range cookies {
		if c.Name == "pt_login_sig" {
			loginSig = c.Value
		}
	}
	log.Println("qrsig", qrsig)
	log.Println("login_sig", loginSig)
	for i := 0; i < 100; i++ {
		core.IsQrCodeExpired(cookies, qrsig, loginSig, core.FromOathUrlGetParams("appid"), core.FromOathUrlGetParams("daid"), core.FromOathUrlGetParams("pt_3rd_aid"))
		time.Sleep(time.Second)
	}

}
