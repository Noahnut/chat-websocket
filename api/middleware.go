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

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("c79ad6f88197be0d2a890f890cab101b837db348ba501723c962176fb54d280d"), nil
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

	print("uid: ", fmt.Sprint(claims["uid"]))

	ctx.Next()
}
