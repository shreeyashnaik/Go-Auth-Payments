package web

import (
	"net/http"

	"github.com/Shreeyash-Naik/Go-Auth/src/core/route"
	"github.com/gin-gonic/gin"
)

func Run() {
	router := gin.Default()
	route.MountRoutes(router)

	// Health Checks
	router.POST("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusCreated, gin.H{
			"msg": "OK",
		})
	})

	router.GET("/healthcheck", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "OK",
		})
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	server.ListenAndServe()

}
