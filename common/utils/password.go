package utils

import (
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

var (
	algorithm = jwt.SigningMethodHS256
)

func SignPassword(password string) (signedPassword string) {
	signedPassword, _ = algorithm.Sign(password, []byte(viper.GetString("PASSWORD_KEY")))

	return signedPassword
}

func VerifyPassword(password, signedPassword string) (isPasswordCorrect bool) {
	if err := algorithm.Verify(password, signedPassword, []byte(viper.GetString("PASSWORD_KEY"))); err != nil {
		return false
	}

	return true
}

func SignUser(userDetail string) (signedUser string) {
	signedUser, _ = algorithm.Sign(userDetail, []byte(viper.GetString("USER_KEY")))

	return signedUser
}

func VerifyUser(userDetail, signedUser string) (isUserValid bool) {
	if err := algorithm.Verify(userDetail, signedUser, []byte(viper.GetString("USER_KEY"))); err != nil {
		return false
	}

	return true
}
