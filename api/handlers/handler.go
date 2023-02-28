package handler

import (
	"context"
	"websocket-service/api/http"
	"websocket-service/config"
	"websocket-service/grpc/client"
	"websocket-service/pkg/errors"
	"websocket-service/pkg/logger"
	"websocket-service/pkg/security"
	"websocket-service/socket"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	cfg      config.Config
	log      logger.LoggerI
	services client.ServiceManagerI
	hub      *socket.Hub
	err      errors.Error
}

func NewHandler(cfg config.Config, log logger.LoggerI, svcs client.ServiceManagerI, hub *socket.Hub) *Handler {
	return &Handler{
		cfg:      cfg,
		log:      log,
		services: svcs,
		hub:      hub,
		err: *errors.NewError(log, "Handler", cfg.HTTPPort),
	}
}

func (h *Handler) handleResponse(c *gin.Context, status http.Status, data interface{}) {
	switch code := status.Code; {
	case code < 300:
		h.log.Info(
			"---Response--->",
			logger.Int("code", status.Code),
			logger.String("status", status.Status),
			logger.Any("description", status.Description),
			// logger.Any("data", data),
		)
	case code < 400:
		h.log.Warn(
			"!!!Response--->",
			logger.Int("code", status.Code),
			logger.String("status", status.Status),
			logger.Any("description", status.Description),
			logger.Any("data", data),
		)
		data = h.err.GetError(data)
	default:
		h.log.Error(
			"!!!Response--->",
			logger.Int("code", status.Code),
			logger.String("status", status.Status),
			logger.Any("description", status.Description),
			logger.Any("data", data),
		)
		data = h.err.GetError(data)
	}

	c.JSON(status.Code, http.Response{
		Status:      status.Status,
		Description: status.Description,
		Data:        data,
	})
}

func (h *Handler) HasAccess(c *gin.Context) {
	reqToken := c.GetHeader("Authorization")
	var err error
	defer h.err.Wrap(&err, "HasAccess", reqToken)

	if len(reqToken) < 25 {
		h.handleResponse(c, http.Forbidden, "token is empty")
		c.Abort()
		return
	}

	token, err := security.ExtractToken(reqToken)
	if err != nil {
		h.handleResponse(c, http.Forbidden, err.Error())
		c.Abort()
		return
	}

	tokenInfo, err := security.ParseClaims(token, h.cfg.SecretKey)
	if err != nil {
		h.handleResponse(c, http.Forbidden, err.Error())
		c.Abort()
		return
	}

	c.Set("ctx", NewContext(c.Request.Context(), "token_info", &tokenInfo))
}

func NewContext(r context.Context, key, value any) context.Context {
	return context.WithValue(r, key, value)
}
