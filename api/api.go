package api

import (
	handler "websocket-service/api/handlers"
	"websocket-service/config"

	"github.com/gin-gonic/gin"
)

func SetUpRouter(h *handler.Handler, cfg config.Config) (r *gin.Engine) {
	r = gin.New()

	r.Use(gin.Logger(), gin.Recovery())

	r.Use(customCORSMiddleware())

	r.GET("/ws", h.HasAccess, h.Ws)
	r.GET("/", h.GetToken)

	return
}

func customCORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.Header("Access-Control-Max-Age", "3600")
		c.Header("Access-Control-Allow-Headers", "Host, Connection, Upgrade, Sec-WebSocket-Key, Sec-WebSocket-Version, Sec-WebSocket-Extensions, Authorization")


		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
