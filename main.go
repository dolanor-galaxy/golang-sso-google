package main

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/jamesonwilliams/golang-sso-google/auth"
	"github.com/jamesonwilliams/golang-sso-google/handlers"
)

func main() {
	router := gin.Default()
	store := sessions.NewCookieStore([]byte(handlers.RandomToken(64)))
	store.Options(sessions.Options{
		Path:   "/",
		MaxAge: 86400 * 7,
	})
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(sessions.Sessions("goquestsession", store))
	router.Static("/css", "./static/css")
	router.Static("/img", "./static/img")
	router.LoadHTMLGlob("templates/*")

	router.GET("/", handlers.LoginHandler)
	router.GET("/auth", handlers.AuthHandler)

	authorized := router.Group("/ok")
	authorized.Use(auth.AuthorizeRequest())
	{
		authorized.GET("/internal", handlers.ReverseProxy)
	}

	router.Run(":9090")
}
