package httpsrv

import (
	"third/gin"
	"backend/common"
	"strings"
	"encoding/json"
	"net/http"
)

var kubeApiserverPath = ""
var kubeApiserverPort = ""
var registryPath = ""
var registryPort = ""
var skyDNSPath = ""
var skyDNSPort = ""
var loginUrl = ""
var sourceType = ""
var redirectUrl = "https://login_in.codoon.com?next=授权"

func respondWithError(code int, message string, c *gin.Context) {
	resp := map[string]interface{}{"state": 1, "msg": message}

	c.JSON(code, resp)
	c.Abort()
}

func AccountAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//var loginUrl string = "https://login_in.codoon.com/check_session"
		sessionCookie, _ := c.Request.Cookie("login_session_id")
		sessionId := sessionCookie.Value
		clientAddr := c.Request.RemoteAddr
		clients := strings.Split(clientAddr, ":")
		clientIP := clients[0]

		reqJson := map[string]string{
			"session_id": sessionId,
			"client_ip": clientIP,
			"source_type": sourceType,
		}

		code, response, err := common.SendFormRequest("GET", loginUrl, reqJson)
		if err != nil {
			respondWithError(code, "access login.in.codoon.com failed", c)
			return
		}

		bytes := []byte(response)
		var res map[string]interface{}

		if err := json.Unmarshal(bytes, &res); err != nil {
			respondWithError(code, "unmarsh response json failed", c)
			return
		}

		if errcode := res["errcode"].(int); errcode != 0 {
			c.JSON(http.StatusOK, gin.H{"status": map[string]interface{}{
				"state": 5,
				"msg": "please login",
			},"data": map[string]interface{}{
				"rd_url": redirectUrl,
			}})
			return
		}

		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		/**
		if c.Request.Method == "OPTIONS" {
				c.Abort()
				return
		}
		**/
		c.Next()
	}
}

func InitExternalConfig(config common.Configure)  {
	kubeApiserverPath = config.External["kubeApiserverPath"]
	kubeApiserverPort = config.External["kubeApiserverPort"]
	registryPath = config.External["registryPath"]
	registryPort = config.External["registryPort"]
	skyDNSPath = config.External["SkyDNSPath"]
	skyDNSPort = config.External["SkyDNSPort"]
	loginUrl = config.External["loginPath"]
	sourceType = config.External["sourceType"]
}

func StartServer() {
	defer common.MyRecovery()
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware())
	router.Use(AccountAuthMiddleware())

	SetupRoutes(router)
	router.Run(common.Config.Listen)
}
