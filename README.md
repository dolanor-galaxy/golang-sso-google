# golang-sso-google

Go language implementation of Single Sign On with Google.

**Credit** to
[GoProgressQuest](https://github.com/Skarlso/goprogressquest), from
which this current repo was forked.

# Installation

Simply `go get github.com/jamesonwilliams/golang-sso-google`.

# Setup

## Google

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
root: `go get -v ./...`

# Running

To run it, simply build & run and navigate to
http://127.0.0.1:9090/login, nothing else should be required.

```
go build
./golang-sso-google
```

[iam-creds]: https://console.developers.google.com/iam-admin/projects

