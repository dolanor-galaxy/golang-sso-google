package main

import (
	"github.com/jamesonwilliams/golang-sso-google/auth"
	"github.com/jamesonwilliams/golang-sso-google/handlers"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
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

	router.GET("/", handlers.IndexHandler)
	router.GET("/login", handlers.LoginHandler)
	router.GET("/auth", handlers.AuthHandler)

	authorized := router.Group("/ok")
	authorized.Use(auth.AuthorizeRequest())
	{
		authorized.GET("/internal", handlers.InternalPageHandler)
	}

	router.Run("127.0.0.1:9090")
}
