package auth

import (
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

const errorTemplate = "error.tmpl"
const sessionKey = "user-id"

// AuthorizeRequest is used to authorize a request for a certain
// end-point group.
func AuthorizeRequest() gin.HandlerFunc {
	return func(context *gin.Context) {
		session := sessions.Default(context)
		sessionValue := session.Get(sessionKey)
		if sessionValue == nil {
			context.HTML(http.StatusUnauthorized, errorTemplate, gin.H{"message": "Please login."})
			context.Abort()
		}
		context.Next()
	}
}
