package httpsrv

import (
	"third/gin"
	"backend/common"
	"fmt"
	"os"
)
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
