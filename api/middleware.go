package api

import (
	"chat-websocket/models"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Middleware struct {
	ctx context.Context
}

func (m *Middleware) AuthMiddleware(ctx *gin.Context) {

	tokenString := ctx.GetHeader("Authorization")
	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": models.ErrTokenInvalid,
		})
		ctx.Abort()
		return
	}
	tokenString = tokenString[len("Bearer "):]
	// TODO: add secret to env
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": models.ErrorInternal,
		})
		log.Println(err)
		ctx.Abort()
		return
	}

	if !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": models.ErrTokenInvalid,
		})
		log.Println(err)
		ctx.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": models.ErrTokenInvalid,
		})
		ctx.Abort()
		return
	}

	ctx.Set("uid", fmt.Sprint(claims["uid"]))
	ctx.Set("context", m.ctx)

	ctx.Next()
}
