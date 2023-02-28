package handler

import (
	"time"
	"websocket-service/pkg/security"

	"github.com/gin-gonic/gin"
)

func (h Handler) GetToken(c *gin.Context) {
	m := map[string]interface{}{
		"user_id": c.Query("user_id"),
		"email":   c.Query("email"),
	}
	token, _ := security.GenerateJWT(m, 100*time.Minute, h.cfg.SecretKey)
	c.JSON(200, gin.H{
		"token": token,
	})
}
