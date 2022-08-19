package web

import (
	"net/http"

	"github.com/Shreeyash-Naik/Go-Auth/src/core/route"
	"github.com/gin-gonic/gin"
)

func Run() {
	router := gin.Default()
	route.MountRoutes(router)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	server.ListenAndServe()
}
