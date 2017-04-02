package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/jamesonwilliams/golang-sso-google/auth"
	"github.com/jamesonwilliams/golang-sso-google/database"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"net/http/httputil"
)

const landingPageTemplate = "index.tmpl"
const successTemplate = "success.tmpl"
const internalTemplate = "internal.tmpl"
const errorTemplate = "error.tmpl"
const authTemplate = "auth.tmpl"

const redirectUrl = "http://ec2-54-190-25-232.us-west-2.compute.amazonaws.com:9090/auth"
const googleOathUserInfoUrl = "https://www.googleapis.com/oauth2/v3/userinfo"

const sessionKey = "user-id"
const sessionStateKey = "state"

var db = database.DynamoDatabase{
	Region:    "us-west-2",
	TableName: "Users",
}

// You have to select your own scope from here ->
// https://developers.google.com/identity/protocols/googlescopes#google_sign-in
var googleAuthScopes = []string{
	"https://www.googleapis.com/auth/userinfo.email",
}

var config *oauth2.Config

// RandomToken generates a random @length length token.
func RandomToken(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)

	return base64.StdEncoding.EncodeToString(bytes)
}

func getLoginURL(state string) string {
	return config.AuthCodeURL(state)
}

func init() {
	config = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  redirectUrl,
		Scopes:       googleAuthScopes,
		Endpoint:     google.Endpoint,
	}
}

// AuthHandler handles authentication of a user and initiates a session.
func AuthHandler(context *gin.Context) {
	// Handle the exchange code to initiate a transport.
	session := sessions.Default(context)
	retrievedState := session.Get(sessionStateKey)
	queryState := context.Request.URL.Query().Get(sessionStateKey)
	if retrievedState != queryState {
		log.Printf("Invalid session state: retrieved: %s; Param: %s", retrievedState, queryState)
		context.HTML(http.StatusUnauthorized, errorTemplate, gin.H{"message": "Invalid session state."})
		return
	}
	code := context.Request.URL.Query().Get("code")
	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Println(err)
		context.HTML(http.StatusBadRequest, errorTemplate, gin.H{"message": "Login failed. Please try again."})
		return
	}

	client := config.Client(oauth2.NoContext, tok)
	userinfo, err := client.Get(googleOathUserInfoUrl)
	if err != nil {
		log.Println(err)
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}
	defer userinfo.Body.Close()
	data, _ := ioutil.ReadAll(userinfo.Body)
	user := auth.User{}
	if err = json.Unmarshal(data, &user); err != nil {
		log.Println(err)
		context.HTML(http.StatusBadRequest, errorTemplate, gin.H{"message": "Error marshalling response. Please try agian."})
		return
	}
	session.Set(sessionKey, user.Email)
	err = session.Save()
	if err != nil {
		log.Println(err)
		context.HTML(http.StatusBadRequest, errorTemplate, gin.H{"message": "Error while saving session. Please try again."})
		return
	}
	seen := false
	if _, dbErr := db.RetrieveUser(user.Email); dbErr == nil {
		seen = true
	} else {
		err = db.SaveUser(&user)
		if err != nil {
			log.Println(err)
			context.HTML(http.StatusBadRequest, errorTemplate, gin.H{"message": "Error while saving user. Please try again."})
			return
		}
	}
	context.HTML(http.StatusOK, successTemplate, gin.H{"name": user.GivenName, "seen": seen, "picture": user.Picture})
}

// LoginHandler handles the login procedure.
func LoginHandler(context *gin.Context) {
	state := RandomToken(32)
	session := sessions.Default(context)
	session.Set(sessionStateKey, state)
	log.Printf("Stored session: %v\n", state)
	session.Save()
	link := getLoginURL(state)
	context.HTML(http.StatusOK, authTemplate, gin.H{"link": link})
}

func ReverseProxy(c *gin.Context) {
	director := func(req *http.Request) {
		r := c.Request
		req = r
		req.URL.Scheme = "http"
        req.URL.Host = "ec2-54-190-25-232.us-west-2.compute.amazonaws.com:9112"
		req.URL.Path = "/"
	}
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(c.Writer, c.Request)
}
