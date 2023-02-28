package helper

import (
	"context"
	"websocket-service/pkg/security"
)

func GetValueContext(ctx context.Context) *security.TokenInfo {
	if ctx.Value("token_info") == nil {
		return &security.TokenInfo{}
	}
	return ctx.Value("token_info").(*security.TokenInfo)
}
