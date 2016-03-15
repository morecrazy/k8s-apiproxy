package httpsrv

import (
	"third/gin"
	"backend/common"
	"fmt"
	"os"
)

var kubeApiserverPath = ""
var kubeApiserverPort = ""
var registryPath = ""
var registryPort = ""
var skyDNSPath = ""
var skyDNSPort = ""

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
}

func StartServer() {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware())
	//router.Use(common.Log())
	SetupRoutes(router)
	err := router.Run(common.Config.Listen)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
