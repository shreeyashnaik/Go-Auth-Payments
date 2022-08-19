package middleware

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

func CheckLooseAuth(ctx *gin.Context) {

}

type MyCustomClaims struct {
	Foo string `json:"foo"`
	jwt.StandardClaims
}

func CheckStrictAuth(ctx *gin.Context) {
	cookie, err := ctx.Cookie("jwt")
	log.Println(cookie)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg": "User needs to Log In",
		})

		ctx.Abort()
	}

	token, err := jwt.ParseWithClaims(cookie, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("USER_KEY")), nil
	})
	fmt.Println(token.Claims)
	fmt.Println(token.Header)
	fmt.Println(token.Valid)

	if err != nil || !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg": "User Unauthorized",
		})
		ctx.Abort()
	}

	claims := token.Claims.(jwt.MapClaims)
	ctx.Keys = map[string]interface{}{
		"user": claims["iss"],
	}

	ctx.Next()
}
