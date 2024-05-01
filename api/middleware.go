package api

import (
	"context"
	"crypto/md5"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	ctx context.Context
}

func (m *Middleware) AuthMiddleware(ctx *gin.Context) {
	test_token := md5.New().Sum([]byte("test"))

	ctx.Set("uid", hex.EncodeToString(test_token))
	ctx.Set("context", m.ctx)

	ctx.Next()
}
