package controller

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/Shreeyash-Naik/Go-Auth/common/cache"
	"github.com/Shreeyash-Naik/Go-Auth/common/db"
	"github.com/Shreeyash-Naik/Go-Auth/common/schemas"
	"github.com/Shreeyash-Naik/Go-Auth/common/utils"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

func Signup(ctx *gin.Context) {
	payload := new(schemas.SignupPayload)
	if err := ctx.ShouldBindJSON(payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid Params",
		})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(payload.Password), 14)
	if err := db.CreateUser(payload.Email, payload.Name, string(hashedPassword)); err != nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"msg": "User Already exists with email",
		})
		return
	}

	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"msg": "User successfully Created",
	})
}

func Login(ctx *gin.Context) {
	payload := new(schemas.LoginPayload)
	if err := ctx.ShouldBindJSON(payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid Params",
		})
		return
	}

	// Get hashed password value from DB and verify it with user entered password
	hashedOriginalPassword := db.RedisDB.LIndex(db.RedisDB.Context(), payload.Email, 0).Val()
	if err := bcrypt.CompareHashAndPassword([]byte(hashedOriginalPassword), []byte(payload.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg": "Incorrect Password",
		})
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    payload.Email,
		ExpiresAt: time.Now().Add(time.Minute * 60).Unix(),
	})

	jwtToken, err := claims.SignedString([]byte(viper.GetString("USER_KEY")))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Could not log in",
		})
		return
	}

	ctx.SetCookie("jwt", jwtToken, 3600, "/", "http://127.0.0.1", false, true)

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"msg": "User successfully Logged In",
		"jwt": jwtToken,
	})
}

func Home(ctx *gin.Context) {

}

func Orders(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("Hello user %s", ctx.Keys["user"]),
	})
}

func CreateOrder(ctx *gin.Context) {

}

func LoginOTP(ctx *gin.Context) {
	p := new(schemas.LoginOTPPayload)
	if err := ctx.ShouldBindJSON(p); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg": "Invalid Params",
		})
		return
	}

	// For OTP, generate 6 digit random number
	otp := 100000 + rand.Intn(999999-100000+1)

	// Store OTP in cache for further verification
	cache.SetValueEx(fmt.Sprint(p.Email, "_otp"), otp, time.Minute*5)

	// Send OTP to requested email via AWS Pinpoint
	go utils.SendOTPEmail(p.Email, otp)

	ctx.JSON(http.StatusOK, gin.H{
		"msg":   "OTP sent to email",
		"email": p.Email,
	})
}

func VerifyOTP(ctx *gin.Context) {
	p := new(schemas.VerifyOTPPayload)
	if err := ctx.ShouldBindJSON(p); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg": "Invalid Params",
		})
		return
	}

	otp := cache.CacheClient.Get(cache.CacheClient.Context(), fmt.Sprintf("%s_otp", p.Email))
	if otp.Val() != strconv.Itoa(p.OTP) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg": "Invalid OTP",
		})
		return
	}

	ctx.JSON(http.StatusUnauthorized, gin.H{
		"msg":   "Successful OTP Verification",
		"email": p.Email,
	})
}
