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

Uses a DynamoDB table called `Users`. You can create this table with:

    aws dynamodb \
        --profile default \
        --region us-west-2 \
    create-table \
        --table-name Users \
        --key-schema \
            AttributeName=email,KeyType=HASH \
        --attribute-definitions \
            AttributeName=email,AttributeType=S \
        --provisioned-throughput \
            ReadCapacityUnits=1,WriteCapacityUnits=1

You will also need to arrange for AWS credentials that have sufficient
permissions to read and write to that table to be in your environment.

Easiest is to just put them in `~/.aws/config`.

To gather all the libraries this project uses, simply execute from the
root:

    go get -v "github.com/gin-gonic/gin"
    go get -v "github.com/gin-gonic/contrib/sessions"
    go get -v "golang.org/x/oauth2"
    go get -v "golang.org/x/oauth2/google"
    go get -v "gopkg.in/mgo.v2"
    go get -v "gopkg.in/mgo.v2/bson"
    # etc.

Or just:

    go get

## Running

    go build
    ./golang-sso-google &
    google-chrome http://127.0.0.1:9090/login &

## Screenshots

![Login Button][login-button]
![Google Redirect][google-redirect]


[iam-creds]: https://console.developers.google.com/apis/credentials
[login-button]: https://nosemaj.org/dl/login-button.png "Login button"
[google-redirect]: https://nosemaj.org/dl/google-redirect.png "Google Redirect"

