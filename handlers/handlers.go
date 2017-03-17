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
)

const landingPageTemplate = "index.tmpl"
const successTemplate = "success.tmpl"
const internalTemplate = "internal.tmpl"
const errorTemplate = "error.tmpl"
const authTemplate = "auth.tmpl"

const credentialsStore = "./credentials.json"
const redirectUrl = "http://127.0.0.1:9090/auth"
const googleOathUserInfoUrl = "https://www.googleapis.com/oauth2/v3/userinfo"

const sessionKey = "user-id"
const sessionStateKey = "state"

// You have to select your own scope from here ->
// https://developers.google.com/identity/protocols/googlescopes#google_sign-in
var googleAuthScopes = []string{
	"https://www.googleapis.com/auth/userinfo.email",
}

var credentials Credentials
var config *oauth2.Config

// Credentials stores Google API developer credentials.
type Credentials struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

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
	file, err := ioutil.ReadFile(credentialsStore)
	if err != nil {
		log.Printf("Error reading file %s: %v\n", credentialsStore, err)
		os.Exit(1)
	}
	json.Unmarshal(file, &credentials)

	config = &oauth2.Config{
		ClientID:     credentials.ClientId,
		ClientSecret: credentials.ClientSecret,
		RedirectURL:  redirectUrl,
		Scopes:       googleAuthScopes,
		Endpoint:     google.Endpoint,
	}
}

// IndexHandler handels /.
func IndexHandler(context *gin.Context) {
	context.HTML(http.StatusOK, landingPageTemplate, gin.H{})
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
	db := database.MongoDBConnection{}
	if _, mongoErr := db.LoadUser(user.Email); mongoErr == nil {
		seen = true
	} else {
		err = db.SaveUser(&user)
		if err != nil {
			log.Println(err)
			context.HTML(http.StatusBadRequest, errorTemplate, gin.H{"message": "Error while saving user. Please try again."})
			return
		}
	}
	context.HTML(http.StatusOK, successTemplate, gin.H{"email": user.Email, "seen": seen})
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

// InternalPageHandler is a rudimentary handler for logged in users.
func InternalPageHandler(context *gin.Context) {
	session := sessions.Default(context)
	userId := session.Get(sessionKey)
	context.HTML(http.StatusOK, internalTemplate, gin.H{"user": userId})
}
