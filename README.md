# golang-sso-google

Go language implementation of Single Sign On with Google.

**Credit** to
[Skarlso](https://github.com/Skarlso), from whom I forked this repo. 

## Installation

    go get github.com/jamesonwilliams/golang-sso-google

## Google Credentials

In order for the Google Authentication to work, you'll need developer
credentials which the this application gathers from a file in the root
directory called `credentials.json`. The structure of this file should
be like this:

```json
{
  "clientId":"hash.apps.googleusercontent.com",
  "clientSecret":"somesecrethash"
}
```

To obtain these credentials, please navigate to this site and follow the
procedure to setup a new project: [Google Developer Console][iam-creds].

## Dependencies

To gather all the libraries this project uses, simply execute from the
root:

    go get -v "github.com/gin-gonic/gin"
    go get -v "github.com/gin-gonic/contrib/sessions"
    go get -v "golang.org/x/oauth2"
    go get -v "golang.org/x/oauth2/google"
    go get -v "gopkg.in/mgo.v2"
    go get -v "gopkg.in/mgo.v2/bson"
    # etc.

Additionally, install MongoDB:

    sudo aptitude install mongodb-server

## Running

    go build
    ./golang-sso-google &
    google-chrome http://127.0.0.1:9090/login &

## Screenshots

![Login Button][login-button]
![Google Redirect][google-redirect]


[iam-creds]: https://console.developers.google.com/iam-admin/projects
[login-button]: https://nosemaj.org/dl/login-button.png "Login button"
[google-redirect]: https://nosemaj.org/dl/google-redirect.png "Google Redirect"

