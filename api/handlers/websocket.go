package handler

import (
	"context"
	"websocket-service/pkg/logger"

	"github.com/gin-gonic/gin"
)

func (h Handler) Ws(c *gin.Context) {
	h.log.Info("Info", logger.Any("room", c.GetHeader("room")), logger.Any("room", c.GetHeader("userId")))

	ctx := NewContext(c.Value("ctx").(context.Context), "room", c.GetHeader("room"))

	h.hub.ServeWs(c.Writer, c.Request, ctx)
}
