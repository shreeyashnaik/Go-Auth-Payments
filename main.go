package main

import (
	"github.com/Shreeyash-Naik/Go-Auth/common/db"
	"github.com/Shreeyash-Naik/Go-Auth/common/utils"
	"github.com/Shreeyash-Naik/Go-Auth/src/core/web"
)

func main() {
	utils.ImportEnv()

	// // AWS Utils
	// utils.ConnectAWS()
	// utils.SendOTPEmail("shreeyashnaikofficial@gmail.com", 123456)

	db.SetDB()
	web.Run()
}
