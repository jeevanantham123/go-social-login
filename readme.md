# social-login with Go

 
 
## Config Project
* Check out this site https://developer.github.com/apps/building-oauth-apps/creating-an-oauth-app/
* After getting credentials from Github save it into the .env file

```.env
CLIENTID = "your github client ID"
CLIENTSECRET = "your github client Secret"
```


## How OAuth2 works with Google
The authorization sequence begins when your application redirects the browser to a Google URL; the URL includes query parameters that indicate the type of access being requested. Google handles the user authentication, session selection, and user consent. The result is an authorization code, which the application can exchange for an access token and a refresh token.

The application should store the refresh token for future use and use the access token to access a Google API. Once the access token expires, the application uses the refresh token to obtain a new one.


## Let's go to the code
We will use the package "golang.org/x/oauth2" that provides support for making OAuth2 authorized and authenticated HTTP requests.

Create a new project(folder) in your workdir in my case I will call it 'go-social-login', and we need to include the package of oauth2.

`go get golang.org/x/oauth2`


So into the project we create a main.go.

```go
package main

import (
	"fmt"
	"net/http"
	"log"
	"github.com/jeevanantham123/go-social-login/handlers"
)

func main() {
	server := &http.Server{
		Addr: fmt.Sprintf(":8000"),
		Handler: handlers.New(),
	}

	log.Printf("Starting HTTP Server. Listening at %q", server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("%v", err)
	} else {
		log.Println("Server closed!")
	}
}

```
We create a simple server using http.Server and run.

Next, we create folder 'handlers' that contains handler of our application, in this folder create 'base.go'.

```go
package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	githubOAuth2 "golang.org/x/oauth2/github"
)

//New func
func New() http.Handler {
	godotenv.Load(".env")
	mux := http.NewServeMux()
	// Root
	config := &oauth2.Config{
		ClientID:     os.Getenv("CLIENTID"),
		ClientSecret: os.Getenv("CLIENTSECRET"),
		RedirectURL:  "http://localhost:8000/oauth/callback",
		Endpoint:     githubOAuth2.Endpoint,
	}
	mux.Handle("/", http.FileServer(http.Dir("templates/")))

	stateConfig := gologin.DebugOnlyCookieConfig
	mux.Handle("/login", github.StateHandler(stateConfig, github.LoginHandler(config, nil)))
	mux.Handle("/oauth/callback", github.StateHandler(stateConfig, github.CallbackHandler(config, issueSession(), nil)))
	mux.Handle("/check", issueSession())
	return mux
}
```

We use **http.ServeMux** to handle our endpoints, next we create the Root endpoint "/" for serving a simple template with a minimmum HTML&CSS in this example we use 'http. http.FileServer', that template is 'index.html' and is in the folder 'templates'.

Also we create two endpoints for Oauth with Google "/login" and "/oauth/callback". Remember when we configured our application in the Github console? The callback url must be the same.


```go
config := &oauth2.Config{
		ClientID:     os.Getenv("CLIENTID"),
		ClientSecret: os.Getenv("CLIENTSECRET"),
		RedirectURL:  "http://localhost:8000/oauth/callback",
		Endpoint:     githubOAuth2.Endpoint,
}
```


### redirect function after getting authorization from user
This redirect func will fires automatically after getting auth from user.We can then use go github interfaces to access the user details.

```go
func issueSession() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		githubUser, err := github.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(githubUser)
	}
	return http.HandlerFunc(fn)
}
```

```
## let's run and test
```bash
go run main.go
```
